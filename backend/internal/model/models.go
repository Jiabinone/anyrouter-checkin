package model

import (
	"github.com/dromara/carbon/v2"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type User struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	Username  string          `gorm:"uniqueIndex;size:50" json:"username"`
	Password  string          `gorm:"size:255" json:"-"`
	CreatedAt carbon.DateTime `json:"created_at" swaggertype:"string" format:"date-time"`
	UpdatedAt carbon.DateTime `json:"updated_at" swaggertype:"string" format:"date-time"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-" swaggerignore:"true"`
}

type Account struct {
	ID          uint             `gorm:"primarykey" json:"id"`
	Session     string           `gorm:"type:text" json:"-"`
	UserID      int              `json:"user_id"`
	Username    string           `gorm:"size:100" json:"username"`
	Role        int              `json:"role"`
	Status      int              `gorm:"default:1" json:"status"`
	Balance     decimal.Decimal  `gorm:"type:decimal(20,2);default:0" json:"balance"`
	LastCheckin *carbon.DateTime `json:"last_checkin" swaggertype:"string" format:"date-time"`
	LastResult  string           `gorm:"size:255" json:"last_result"`
	CreatedAt   carbon.DateTime  `json:"created_at" swaggertype:"string" format:"date-time"`
	UpdatedAt   carbon.DateTime  `json:"updated_at" swaggertype:"string" format:"date-time"`
}

type CronTask struct {
	ID         uint             `gorm:"primarykey" json:"id"`
	Name       string           `gorm:"size:100" json:"name"`
	CronExpr   string           `gorm:"size:50" json:"cron_expr"`
	TaskType   string           `gorm:"size:50" json:"task_type"`
	AccountIDs string           `gorm:"type:text" json:"account_ids"`
	Status     int              `gorm:"default:1" json:"status"`
	LastRun    *carbon.DateTime `json:"last_run" swaggertype:"string" format:"date-time"`
	NextRun    *carbon.DateTime `json:"next_run" swaggertype:"string" format:"date-time"`
	CreatedAt  carbon.DateTime  `json:"created_at" swaggertype:"string" format:"date-time"`
}

type Config struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	Key       string          `gorm:"uniqueIndex;size:100" json:"key"`
	Value     string          `gorm:"type:text" json:"value"`
	Category  string          `gorm:"size:50;index" json:"category"`
	UpdatedAt carbon.DateTime `json:"updated_at" swaggertype:"string" format:"date-time"`
}

type CheckinLog struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	AccountID uint            `gorm:"index" json:"account_id"`
	Success   bool            `json:"success"`
	Message   string          `gorm:"type:text" json:"message"`
	CreatedAt carbon.DateTime `json:"created_at" swaggertype:"string" format:"date-time"`
}
