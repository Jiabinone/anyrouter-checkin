package repository

import "anyrouter-checkin/internal/model"

func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CountUsers() (int64, error) {
	var count int64
	if err := DB.Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func CreateUser(user *model.User) error {
	return DB.Create(user).Error
}

func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func SaveUser(user *model.User) error {
	return DB.Save(user).Error
}
