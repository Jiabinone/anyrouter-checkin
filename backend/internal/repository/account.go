package repository

import "anyrouter-checkin/internal/model"

func ListAccounts() ([]model.Account, error) {
	var accounts []model.Account
	if err := DB.Order("id desc").Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func GetAccountByID(id uint) (*model.Account, error) {
	var account model.Account
	if err := DB.First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func CreateAccount(account *model.Account) error {
	return DB.Create(account).Error
}

func SaveAccount(account *model.Account) error {
	return DB.Save(account).Error
}

func DeleteAccount(id uint) error {
	return DB.Delete(&model.Account{}, id).Error
}
