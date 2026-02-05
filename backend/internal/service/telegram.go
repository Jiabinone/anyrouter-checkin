package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	neturl "net/url"
	"strings"
	"text/template"

	"anyrouter-checkin/internal/repository"

	"gorm.io/gorm"
)

var ErrNoSuccessfulCheckinLog = errors.New("暂无成功签到记录")

func SendTelegramMessage(message string) error {
	if GetConfig("telegram.enabled") != "true" {
		return nil
	}

	botToken := GetConfig("telegram.bot_token")
	chatID := GetConfig("telegram.chat_id")
	if botToken == "" || chatID == "" {
		return fmt.Errorf("telegram 配置不完整")
	}

	apiBase := strings.TrimSpace(GetConfig("telegram.api_base"))
	apiBase = strings.TrimRight(apiBase, "/")
	if apiBase == "" {
		apiBase = "https://api.telegram.org"
	}
	endpoint := fmt.Sprintf("%s/bot%s/sendMessage", apiBase, botToken)
	payload := map[string]string{
		"chat_id":    chatID,
		"text":       message,
		"parse_mode": "HTML",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	proxyURL := strings.TrimSpace(GetConfig("telegram.proxy_url"))
	if proxyURL != "" {
		proxy, err := neturl.Parse(proxyURL)
		if err != nil {
			return fmt.Errorf("代理地址无效")
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}

	resp, err := client.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram API 返回 %d", resp.StatusCode)
	}
	return nil
}

func renderCheckinTemplate(accountName string, success bool, result string) (string, error) {
	tplStr := GetConfig("telegram.template")
	if tplStr == "" {
		tplStr = `<b>AnyRouter 签到系统</b>
用户名：<code>{{.Username}}</code>
状态：{{if .Success}}<b>成功 ✅</b>{{else}}<b>失败 ❌</b>{{end}}
结果：
<pre>{{.Result}}</pre>`
	}

	tpl, err := template.New("msg").Parse(tplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, map[string]interface{}{
		"Username": html.EscapeString(accountName),
		"Success":  success,
		"Result":   html.EscapeString(result),
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func SendCheckinNotification(accountName string, success bool, result string) {
	message, err := renderCheckinTemplate(accountName, success, result)
	if err != nil {
		return
	}
	SendTelegramMessage(message)
}

func SendTestCheckinNotification() error {
	log, err := repository.GetLatestSuccessfulCheckinLog()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoSuccessfulCheckinLog
		}
		return err
	}

	accountName := fmt.Sprintf("账号ID:%d", log.AccountID)
	account, err := repository.GetAccountByID(log.AccountID)
	if err == nil {
		if displayName := strings.TrimSpace(account.Username); displayName != "" {
			accountName = displayName
		}
	}

	message, err := renderCheckinTemplate(accountName, log.Success, log.Message)
	if err != nil {
		return err
	}

	return SendTelegramMessage(message)
}
