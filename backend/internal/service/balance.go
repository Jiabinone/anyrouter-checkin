package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time" // 仅用于 time.Duration 类型

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type userSelfResponse struct {
	Data struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Role     int    `json:"role"`
		Status   int    `json:"status"`
		Quota    int64  `json:"quota"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type AccountSelfInfo struct {
	UserID   int
	Username string
	Role     int
	Status   int
	Balance  decimal.Decimal
}

func fetchAccountSelf(sessionCookie string, userID int) (AccountSelfInfo, error) {
	info, err := fetchAccountSelfInternal(sessionCookie, userID)
	if err != nil {
		zap.L().Warn("获取账号信息失败", zap.Int("user_id", userID), zap.Error(err))
	}
	return info, err
}

func fetchAccountSelfInternal(sessionCookie string, userID int) (AccountSelfInfo, error) {
	info, err := fetchAccountSelfAttempt(sessionCookie, userID)
	if err == nil || userID <= 0 || !errors.Is(err, ErrInvalidSession) {
		return info, err
	}
	return fetchAccountSelfAttempt(sessionCookie, 0)
}

func fetchAccountSelfAttempt(sessionCookie string, userID int) (AccountSelfInfo, error) {
	baseURL := "https://anyrouter.top"
	client := &http.Client{Timeout: 30 * time.Second}
	sessionValue := extractSessionValue(sessionCookie)
	headers := map[string]string{
		"accept":          "application/json, text/plain, */*",
		"accept-language": "zh-CN,zh;q=0.9",
		"pragma":          "no-cache",
		"user-agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
	}
	if userID > 0 {
		headers["new-api-user"] = strconv.Itoa(userID)
	}

	if sessionValue == "" {
		return AccountSelfInfo{}, fmt.Errorf("session 为空")
	}
	acwScV2, err := fetchAcwScV2(baseURL, headers)
	if err != nil {
		return AccountSelfInfo{}, fmt.Errorf("获取 acw_sc__v2 失败: %v", err)
	}

	req, err := http.NewRequest("GET", baseURL+"/api/user/self", nil)
	if err != nil {
		return AccountSelfInfo{}, fmt.Errorf("创建请求失败: %v", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	cookieHeader := buildSelfCookieHeader(sessionValue, acwScV2)
	req.Header.Set("cookie", cookieHeader)

	resp, err := client.Do(req)
	if err != nil {
		return AccountSelfInfo{}, fmt.Errorf("请求失败: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return AccountSelfInfo{}, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return AccountSelfInfo{}, ErrInvalidSession
	}
	if resp.StatusCode != http.StatusOK {
		return AccountSelfInfo{}, fmt.Errorf("请求失败: %s", resp.Status)
	}

	var payload userSelfResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return AccountSelfInfo{}, fmt.Errorf("解析响应失败: %v", err)
	}
	if !payload.Success {
		if isUnauthorizedMessage(payload.Message) {
			return AccountSelfInfo{}, ErrInvalidSession
		}
		msg := payload.Message
		if msg == "" {
			msg = "请求失败"
		}
		return AccountSelfInfo{}, errors.New(msg)
	}

	// quota 返回值比实际美元额度放大 100 倍，因此需除以 500000（5000 * 100）
	balance := decimal.NewFromInt(payload.Data.Quota).DivRound(decimal.NewFromInt(500000), 2)
	return AccountSelfInfo{
		UserID:   payload.Data.ID,
		Username: payload.Data.Username,
		Role:     payload.Data.Role,
		Status:   payload.Data.Status,
		Balance:  balance,
	}, nil
}

func fetchAcwScV2(baseURL string, headers map[string]string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", baseURL+"/", nil)
	if err != nil {
		return "", err
	}
	if v := headers["accept-language"]; v != "" {
		req.Header.Set("accept-language", v)
	}
	if v := headers["user-agent"]; v != "" {
		req.Header.Set("user-agent", v)
	}
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`(?i)arg1='([a-f0-9]+)'`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", fmt.Errorf("未获取到 arg1")
	}
	return generateAcwScV2(matches[1])
}

func buildSelfCookieHeader(sessionValue, acwScV2 string) string {
	return strings.Join([]string{
		"session=" + sessionValue,
		"acw_sc__v2=" + acwScV2,
	}, "; ")
}

func isUnauthorizedMessage(message string) bool {
	if message == "" {
		return false
	}
	lower := strings.ToLower(message)
	return strings.Contains(message, "未授权") || strings.Contains(lower, "unauthorized")
}

func FetchAcwScV2() (string, error) {
	return fetchAcwScV2("https://anyrouter.top", map[string]string{
		"accept-language": "zh-CN,zh;q=0.9",
		"user-agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
	})
}
