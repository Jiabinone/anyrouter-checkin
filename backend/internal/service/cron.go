package service

import (
	"encoding/json"
	"sync"

	"anyrouter-checkin/internal/model"
	"anyrouter-checkin/internal/repository"

	"github.com/dromara/carbon/v2"
	"github.com/robfig/cron/v3"
)

var (
	scheduler *cron.Cron
	taskIDs   = make(map[uint]cron.EntryID)
	mu        sync.Mutex
)

func InitCron() {
	scheduler = cron.New()
	scheduler.Start()

	var tasks []model.CronTask
	repository.DB.Where("status = ?", 1).Find(&tasks)
	for _, task := range tasks {
		RegisterTask(task)
	}
}

func RegisterTask(task model.CronTask) error {
	mu.Lock()
	defer mu.Unlock()

	if oldID, exists := taskIDs[task.ID]; exists {
		scheduler.Remove(oldID)
	}

	entryID, err := scheduler.AddFunc(task.CronExpr, func() {
		ExecuteTask(task.ID)
	})
	if err != nil {
		return err
	}

	taskIDs[task.ID] = entryID
	updateNextRun(task.ID)
	return nil
}

func UnregisterTask(taskID uint) {
	mu.Lock()
	defer mu.Unlock()

	if entryID, exists := taskIDs[taskID]; exists {
		scheduler.Remove(entryID)
		delete(taskIDs, taskID)
	}
}

func ExecuteTask(taskID uint) {
	var task model.CronTask
	if err := repository.DB.First(&task, taskID).Error; err != nil {
		return
	}

	var accountIDs []uint
	json.Unmarshal([]byte(task.AccountIDs), &accountIDs)

	for _, accID := range accountIDs {
		var account model.Account
		if repository.DB.First(&account, accID).Error == nil {
			success, result := CheckinAccount(accID)
			SendCheckinNotification(account.Name, success, result)
		}
	}

	now := carbon.DateTime{Carbon: carbon.Now()}
	task.LastRun = &now
	repository.DB.Save(&task)
	updateNextRun(taskID)
}

func updateNextRun(taskID uint) {
	var task model.CronTask
	if repository.DB.First(&task, taskID).Error != nil {
		return
	}

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(task.CronExpr)
	if err != nil {
		return
	}

	nextTime := schedule.Next(carbon.Now().StdTime())
	next := carbon.DateTime{Carbon: carbon.CreateFromStdTime(nextTime)}
	task.NextRun = &next
	repository.DB.Save(&task)
}
