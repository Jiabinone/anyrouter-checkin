package service

import "anyrouter-checkin/internal/model"

func ListCronTasks() ([]model.CronTask, error) {
	return listCronTasks()
}

func CreateCronTask(task model.CronTask) (model.CronTask, error) {
	created, err := createCronTask(&task)
	if err != nil {
		return model.CronTask{}, err
	}
	if created.Status == 1 {
		if err := RegisterTask(created); err != nil {
			return model.CronTask{}, err
		}
	}
	return created, nil
}

func UpdateCronTask(id uint, req model.CronTask) (model.CronTask, error) {
	updated, err := updateCronTask(id, req)
	if err != nil {
		return model.CronTask{}, err
	}
	if updated.Status == 1 {
		if err := RegisterTask(updated); err != nil {
			return model.CronTask{}, err
		}
	} else {
		UnregisterTask(updated.ID)
	}
	return updated, nil
}

func DeleteCronTask(id uint) error {
	UnregisterTask(id)
	return deleteCronTask(id)
}
