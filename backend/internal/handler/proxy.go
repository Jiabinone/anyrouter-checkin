package handler

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"anyrouter-checkin/internal/service"
	"anyrouter-checkin/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	anyRouterBaseURL = "https://anyrouter.top"
	proxyTimeout     = 30 * time.Second
)

// AnyRouterProxy 反向代理
// @Summary 代理转发请求到 AnyRouter
// @Tags 代理
// @Accept json
// @Produce json
// @Param path path string true "目标路径"
// @Param Authorization header string true "Bearer {session}"
// @Success 200 {object} map[string]any
// @Router /anyrouter/{path} [get]
func AnyRouterProxy(c *gin.Context) {
	session := extractProxySession(c.GetHeader("Authorization"))
	if session == "" {
		response.Unauthorized(c)
		return
	}

	sessionInfo, err := service.ParseSession(session)
	if err != nil {
		zap.L().Warn("解析 session 失败", zap.Error(err), zap.String("session_prefix", session[:min(50, len(session))]))
		response.Error(c, 400, "invalid session: "+err.Error())
		return
	}

	acwScV2, err := service.FetchAcwScV2()
	if err != nil {
		zap.L().Warn("获取 acw_sc__v2 失败", zap.Error(err))
	}

	targetPath := c.Param("path")
	if targetPath == "" {
		targetPath = "/"
	}
	targetURL := anyRouterBaseURL + targetPath
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	proxyReq, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		zap.L().Error("创建代理请求失败", zap.Error(err))
		response.Error(c, 500, "failed to create request")
		return
	}

	setProxyHeaders(proxyReq, session, acwScV2, sessionInfo.UserID)

	client := &http.Client{Timeout: proxyTimeout}
	resp, err := client.Do(proxyReq)
	if err != nil {
		zap.L().Error("代理请求失败", zap.String("url", targetURL), zap.Error(err))
		response.Error(c, 502, "proxy request failed")
		return
	}
	defer resp.Body.Close()

	copyResponseHeaders(c, resp)
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

func extractProxySession(auth string) string {
	if auth == "" {
		return ""
	}
	return strings.TrimPrefix(auth, "Bearer ")
}

func setProxyHeaders(req *http.Request, session, acwScV2 string, userID int) {
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("cache-control", "no-store")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Not(A:Brand";v="8", "Chromium";v="144", "Google Chrome";v="144"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("referer", anyRouterBaseURL+"/console/personal")

	if userID > 0 {
		req.Header.Set("new-api-user", strconv.Itoa(userID))
	}

	cookie := "session=" + session
	if acwScV2 != "" {
		cookie += "; acw_sc__v2=" + acwScV2
	}
	req.Header.Set("cookie", cookie)
}

func copyResponseHeaders(c *gin.Context, resp *http.Response) {
	for key, values := range resp.Header {
		lower := strings.ToLower(key)
		if lower == "transfer-encoding" || lower == "connection" || lower == "content-length" {
			continue
		}
		for _, v := range values {
			c.Header(key, v)
		}
	}
}
