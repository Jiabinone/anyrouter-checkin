package repository

import (
	"time"

	"anyrouter-checkin/internal/model"
)

func CreateCheckinLog(log *model.CheckinLog) error {
	return DB.Create(log).Error
}

func GetLatestSuccessfulCheckinLog() (*model.CheckinLog, error) {
	var log model.CheckinLog
	if err := DB.Where("success = ?", true).Order("id desc").First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func ListCheckinLogs(limit int) ([]model.CheckinLog, error) {
	var logs []model.CheckinLog
	query := DB.Order("id desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func CountSuccessfulAccounts(start, end time.Time) (int64, error) {
	var count int64
	if err := DB.Model(&model.CheckinLog{}).
		Distinct("account_id").
		Where("success = ?", true).
		Where("created_at >= ? AND created_at <= ?", start, end).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
