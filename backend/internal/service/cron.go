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

	tasks, err := repository.ListEnabledCronTasks()
	if err != nil {
		return
	}
	for _, task := range tasks {
		if err := RegisterTask(task); err != nil {
			continue
		}
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
	task, err := repository.GetCronTaskByID(taskID)
	if err != nil {
		return
	}

	var accountIDs []uint
	if err := json.Unmarshal([]byte(task.AccountIDs), &accountIDs); err != nil {
		return
	}

	for _, accID := range accountIDs {
		account, err := repository.GetAccountByID(accID)
		if err != nil {
			continue
		}
		success, result := CheckinAccount(accID)
		SendCheckinNotification(account.Name, success, result)
	}

	now := carbon.DateTime{Carbon: carbon.Now()}
	task.LastRun = &now
	if err := repository.SaveCronTask(task); err != nil {
		return
	}
	updateNextRun(taskID)
}

func updateNextRun(taskID uint) {
	task, err := repository.GetCronTaskByID(taskID)
	if err != nil {
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
	if err := repository.SaveCronTask(task); err != nil {
		return
	}
}

func listCronTasks() ([]model.CronTask, error) {
	return repository.ListCronTasks()
}

func createCronTask(task *model.CronTask) (model.CronTask, error) {
	if err := repository.CreateCronTask(task); err != nil {
		return model.CronTask{}, err
	}
	return *task, nil
}

func updateCronTask(id uint, req model.CronTask) (model.CronTask, error) {
	task, err := repository.GetCronTaskByID(id)
	if err != nil {
		return model.CronTask{}, err
	}
	task.Name = req.Name
	task.CronExpr = req.CronExpr
	task.AccountIDs = req.AccountIDs
	task.Status = req.Status
	if err := repository.SaveCronTask(task); err != nil {
		return model.CronTask{}, err
	}
	return *task, nil
}

func deleteCronTask(id uint) error {
	return repository.DeleteCronTask(id)
}
