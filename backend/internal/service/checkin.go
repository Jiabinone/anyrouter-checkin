package service

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time" // 仅用于 time.Duration 类型

	"anyrouter-checkin/internal/model"
	"anyrouter-checkin/internal/repository"

	"github.com/dromara/carbon/v2"
)

type SessionInfo struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	Status   int    `json:"status"`
	Group    string `json:"group"`
}

func ParseSession(sessionCookie string) (*SessionInfo, error) {
	decoded, err := base64.URLEncoding.DecodeString(sessionCookie)
	if err != nil {
		return nil, err
	}

	str := string(decoded)
	parts := strings.Split(str, "|")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid session format")
	}

	gobData, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	info := &SessionInfo{}
	kvPairs := make(map[string]interface{})
	var lastKey string
	data := gobData

	for i := 0; i < len(data); {
		if i+6 < len(data) && string(data[i:i+6]) == "string" {
			i += 6
			if i+4 < len(data) && data[i] == 0x0c {
				valLen := int(data[i+3])
				if valLen > 0 && i+4+valLen <= len(data) {
					val := string(data[i+4 : i+4+valLen])
					if lastKey == "" {
						lastKey = val
					} else {
						kvPairs[lastKey] = val
						lastKey = val
					}
					i += 4 + valLen
					continue
				}
			}
		}
		if i+3 < len(data) && string(data[i:i+3]) == "int" {
			i += 3
			if i+4 < len(data) && data[i] == 0x04 {
				encType := data[i+1]
				var val int
				switch encType {
				case 0x04:
					val = int(data[i+4])<<8 | int(data[i+5])
					i += 6
				case 0x02:
					val = int(data[i+3])
					i += 4
				}
				if lastKey != "" {
					kvPairs[lastKey] = val
					lastKey = ""
				}
				continue
			}
		}
		i++
	}

	if v, ok := kvPairs["id"].(int); ok {
		info.UserID = v
	}
	if v, ok := kvPairs["username"].(string); ok {
		info.Username = v
	}
	if v, ok := kvPairs["role"].(int); ok {
		info.Role = v
	}
	if v, ok := kvPairs["status"].(int); ok {
		info.Status = v
	}
	if v, ok := kvPairs["group"].(string); ok {
		info.Group = v
	}

	return info, nil
}

func generateAcwScV2(arg1 string) (string, error) {
	m := []int{0xf, 0x23, 0x1d, 0x18, 0x21, 0x10, 0x1, 0x26, 0xa, 0x9,
		0x13, 0x1f, 0x28, 0x1b, 0x16, 0x17, 0x19, 0xd, 0x6, 0xb,
		0x27, 0x12, 0x14, 0x8, 0xe, 0x15, 0x20, 0x1a, 0x2, 0x1e,
		0x7, 0x4, 0x11, 0x5, 0x3, 0x1c, 0x22, 0x25, 0xc, 0x24}
	p := "3000176000856006061501533003690027800375"

	q := make([]byte, len(m))
	for x := 0; x < len(arg1); x++ {
		for z := 0; z < len(m); z++ {
			if m[z] == x+1 {
				q[z] = arg1[x]
			}
		}
	}
	u := string(q)

	v := ""
	minLen := len(u)
	if len(p) < minLen {
		minLen = len(p)
	}
	for x := 0; x < minLen; x += 2 {
		a, err := strconv.ParseInt(u[x:x+2], 16, 64)
		if err != nil {
			return "", err
		}
		b, err := strconv.ParseInt(p[x:x+2], 16, 64)
		if err != nil {
			return "", err
		}
		v += fmt.Sprintf("%02x", a^b)
	}

	return v, nil
}

func Checkin(sessionCookie string) (string, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", fmt.Errorf("初始化 Cookie 失败: %v", err)
	}
	client := &http.Client{Jar: jar, Timeout: 30 * time.Second}

	baseURL := "https://anyrouter.top"
	headers := map[string]string{
		"accept":          "application/json, text/plain, */*",
		"accept-language": "zh-CN,zh;q=0.9",
		"cache-control":   "no-store",
		"origin":          baseURL,
		"referer":         baseURL + "/console/token",
		"user-agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
	}

	req, err := http.NewRequest("GET", baseURL+"/", nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	re := regexp.MustCompile(`arg1='([A-F0-9]+)'`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		acwScV2, err := generateAcwScV2(matches[1])
		if err != nil {
			return "", fmt.Errorf("生成 Cookie 失败: %v", err)
		}
		u, err := url.Parse(baseURL)
		if err != nil {
			return "", fmt.Errorf("解析地址失败: %v", err)
		}
		jar.SetCookies(u, []*http.Cookie{{Name: "acw_sc__v2", Value: acwScV2}})
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("解析地址失败: %v", err)
	}
	jar.SetCookies(u, []*http.Cookie{{Name: "session", Value: sessionCookie}})

	req, err = http.NewRequest("POST", baseURL+"/api/user/sign_in", nil)
	if err != nil {
		return "", fmt.Errorf("创建签到请求失败: %v", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err = client.Do(req)
	if err != nil {
		return "", fmt.Errorf("签到请求失败: %v", err)
	}
	body, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("读取签到响应失败: %v", err)
	}

	return string(body), nil
}

func CheckinAccount(accountID uint) (bool, string) {
	account, err := repository.GetAccountByID(accountID)
	if err != nil {
		return false, "账号不存在"
	}

	result, err := Checkin(account.Session)
	if err != nil {
		return false, err.Error()
	}

	now := carbon.DateTime{Carbon: carbon.Now()}
	account.LastCheckin = &now
	account.LastResult = result
	if err := repository.SaveAccount(account); err != nil {
		return false, "保存签到结果失败: " + err.Error()
	}

	success := strings.Contains(result, "success") || strings.Contains(result, "已签到")
	if err := repository.CreateCheckinLog(&model.CheckinLog{
		AccountID: accountID,
		Success:   success,
		Message:   result,
	}); err != nil {
		return false, "记录签到日志失败: " + err.Error()
	}

	SendCheckinNotification(account.Name, success, result)

	return success, result
}
