package repository

import (
	"errors"

	"anyrouter-checkin/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	if err := DB.AutoMigrate(
		&model.User{},
		&model.Account{},
		&model.CronTask{},
		&model.Config{},
		&model.CheckinLog{},
	); err != nil {
		return err
	}
	return nil
}

func InitDefaultConfigs() {
	defaultTelegramTemplate := `<b>AnyRouter 签到系统</b>
用户名：<code>{{.Username}}</code>
状态：{{if .Success}}<b>成功 ✅</b>{{else}}<b>失败 ❌</b>{{end}}
结果：
<pre>{{.Result}}</pre>`
	defaults := []model.Config{
		{Key: "telegram.api_base", Value: "https://api.telegram.org", Category: "telegram"},
		{Key: "telegram.bot_token", Value: "", Category: "telegram"},
		{Key: "telegram.chat_id", Value: "", Category: "telegram"},
		{Key: "telegram.enabled", Value: "false", Category: "telegram"},
		{Key: "telegram.proxy_url", Value: "", Category: "telegram"},
		{Key: "telegram.template", Value: defaultTelegramTemplate, Category: "telegram"},
	}
	for _, c := range defaults {
		var cfg model.Config
		if err := DB.Where("`key` = ?", c.Key).First(&cfg).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := DB.Create(&c).Error; err != nil {
					continue
				}
			}
			continue
		}

		if c.Key == "telegram.template" {
			continue
		}
	}
}
