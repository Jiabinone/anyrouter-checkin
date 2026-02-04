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

	return DB.AutoMigrate(
		&model.User{},
		&model.Account{},
		&model.CronTask{},
		&model.Config{},
		&model.CheckinLog{},
	)
}

func InitDefaultConfigs() {
	defaultTelegramTemplate := `<b>AnyRouter 签到系统</b>
你好，<code>{{.Name}}</code>
状态：{{if .Success}}<b>成功 ✅</b>{{else}}<b>失败 ❌</b>{{end}}
结果：
<pre>{{.Result}}</pre>`
	legacyTelegramTemplates := map[string]struct{}{
		"签到结果: {{.Result}}": {},
		"【{{.Name}}】签到{{if .Success}}成功{{else}}失败{{end}}: {{.Result}}":                                                                 {},
		"签到通知\n账号：{{.Name}}\n状态：{{if .Success}}成功{{else}}失败{{end}}\n结果：{{.Result}}":                                                    {},
		"<b>签到通知</b>\n账号：<code>{{.Name}}</code>\n状态：{{if .Success}}<b>成功</b>{{else}}<b>失败</b>{{end}}\n结果：\n<pre>{{.Result}}</pre>":     {},
		"<b>签到提醒</b>\n你好，<code>{{.Name}}</code>\n状态：{{if .Success}}<b>成功 ✅</b>{{else}}<b>失败 ❌</b>{{end}}\n结果：\n<pre>{{.Result}}</pre>": {},
	}

	defaults := []model.Config{
		{Key: "telegram.bot_token", Value: "", Category: "telegram"},
		{Key: "telegram.chat_id", Value: "", Category: "telegram"},
		{Key: "telegram.enabled", Value: "false", Category: "telegram"},
		{Key: "telegram.template", Value: defaultTelegramTemplate, Category: "telegram"},
	}
	for _, c := range defaults {
		var cfg model.Config
		if err := DB.Where("`key` = ?", c.Key).First(&cfg).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				DB.Create(&c)
			}
			continue
		}

		if c.Key == "telegram.template" {
			if _, ok := legacyTelegramTemplates[cfg.Value]; ok {
				cfg.Value = c.Value
				DB.Save(&cfg)
			}
		}
	}
}
