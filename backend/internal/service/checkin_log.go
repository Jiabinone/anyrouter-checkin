package service

import (
	"anyrouter-checkin/internal/model"
	"anyrouter-checkin/internal/repository"

	"github.com/dromara/carbon/v2"
)

type CheckinLogSummary struct {
	Logs                     []model.CheckinLog `json:"logs"`
	TodayCheckinAccountCount int64              `json:"today_checkin_account_count"`
}

func GetCheckinLogSummary(limit int) (CheckinLogSummary, error) {
	var logs []model.CheckinLog
	query := repository.DB.Order("id desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&logs).Error; err != nil {
		return CheckinLogSummary{}, err
	}

	now := carbon.Now()
	start := now.StartOfDay().StdTime()
	end := now.EndOfDay().StdTime()

	var count int64
	if err := repository.DB.Model(&model.CheckinLog{}).
		Distinct("account_id").
		Where("success = ?", true).
		Where("created_at >= ? AND created_at <= ?", start, end).
		Count(&count).Error; err != nil {
		return CheckinLogSummary{}, err
	}

	return CheckinLogSummary{
		Logs:                     logs,
		TodayCheckinAccountCount: count,
	}, nil
}
