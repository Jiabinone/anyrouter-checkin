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
	logs, err := repository.ListCheckinLogs(limit)
	if err != nil {
		return CheckinLogSummary{}, err
	}

	now := carbon.Now()
	start := now.StartOfDay().StdTime()
	end := now.EndOfDay().StdTime()

	count, err := repository.CountSuccessfulAccounts(start, end)
	if err != nil {
		return CheckinLogSummary{}, err
	}

	return CheckinLogSummary{
		Logs:                     logs,
		TodayCheckinAccountCount: count,
	}, nil
}
