package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"text/template"

	"anyrouter-checkin/internal/model"
	"anyrouter-checkin/internal/repository"

	"gorm.io/gorm"
)

var ErrNoSuccessfulCheckinLog = errors.New("暂无成功签到记录")

func GetConfig(key string) string {
	var cfg model.Config
	repository.DB.Where("`key` = ?", key).First(&cfg)
	return cfg.Value
}

func SetConfig(key, value, category string) error {
	var cfg model.Config
	result := repository.DB.Where("`key` = ?", key).First(&cfg)
	if result.Error != nil {
		cfg = model.Config{Key: key, Value: value, Category: category}
		return repository.DB.Create(&cfg).Error
	}
	cfg.Value = value
	return repository.DB.Save(&cfg).Error
}

func SendTelegramMessage(message string) error {
	if GetConfig("telegram.enabled") != "true" {
		return nil
	}

	botToken := GetConfig("telegram.bot_token")
	chatID := GetConfig("telegram.chat_id")
	if botToken == "" || chatID == "" {
		return fmt.Errorf("Telegram 配置不完整")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	payload := map[string]string{
		"chat_id":    chatID,
		"text":       message,
		"parse_mode": "HTML",
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Telegram API 返回 %d", resp.StatusCode)
	}
	return nil
}

func renderCheckinTemplate(accountName string, success bool, result string) (string, error) {
	tplStr := GetConfig("telegram.template")
	if tplStr == "" {
		tplStr = `<b>AnyRouter 签到系统</b>
你好，<code>{{.Name}}</code>
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
		"Name":    html.EscapeString(accountName),
		"Success": success,
		"Result":  html.EscapeString(result),
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
	var log model.CheckinLog
	if err := repository.DB.Where("success = ?", true).Order("id desc").First(&log).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoSuccessfulCheckinLog
		}
		return err
	}

	accountName := fmt.Sprintf("账号ID:%d", log.AccountID)
	var account model.Account
	if err := repository.DB.First(&account, log.AccountID).Error; err == nil && account.Name != "" {
		accountName = account.Name
	}

	message, err := renderCheckinTemplate(accountName, log.Success, log.Message)
	if err != nil {
		return err
	}

	return SendTelegramMessage(message)
}
