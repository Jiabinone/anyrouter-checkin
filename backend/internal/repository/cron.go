package repository

import "anyrouter-checkin/internal/model"

func ListCronTasks() ([]model.CronTask, error) {
	var tasks []model.CronTask
	if err := DB.Order("id desc").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func ListEnabledCronTasks() ([]model.CronTask, error) {
	var tasks []model.CronTask
	if err := DB.Where("status = ?", 1).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetCronTaskByID(id uint) (*model.CronTask, error) {
	var task model.CronTask
	if err := DB.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func CreateCronTask(task *model.CronTask) error {
	return DB.Create(task).Error
}

func SaveCronTask(task *model.CronTask) error {
	return DB.Save(task).Error
}

func DeleteCronTask(id uint) error {
	return DB.Delete(&model.CronTask{}, id).Error
}
