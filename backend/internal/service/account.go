package service

import (
	"errors"
	"fmt"

	"anyrouter-checkin/internal/model"
	"anyrouter-checkin/internal/repository"

	"gorm.io/gorm"
)

var ErrInvalidSession = errors.New("Session 无效")

func ListAccounts() ([]model.Account, error) {
	return repository.ListAccounts()
}

func CreateAccount(name, session string) (model.Account, error) {
	info, err := ParseSession(session)
	if err != nil {
		return model.Account{}, fmt.Errorf("%w: %v", ErrInvalidSession, err)
	}

	account := model.Account{
		Name:     name,
		Session:  session,
		UserID:   info.UserID,
		Username: info.Username,
		Role:     info.Role,
		Status:   1,
	}

	if err := repository.CreateAccount(&account); err != nil {
		return model.Account{}, err
	}

	return account, nil
}

func UpdateAccount(id uint, name, session string) (model.Account, error) {
	account, err := repository.GetAccountByID(id)
	if err != nil {
		return model.Account{}, err
	}

	if session != "" && session != "unchanged" {
		info, err := ParseSession(session)
		if err != nil {
			return model.Account{}, fmt.Errorf("%w: %v", ErrInvalidSession, err)
		}
		account.Session = session
		account.UserID = info.UserID
		account.Username = info.Username
		account.Role = info.Role
	}

	account.Name = name
	if err := repository.SaveAccount(account); err != nil {
		return model.Account{}, err
	}

	return *account, nil
}

func DeleteAccount(id uint) error {
	return repository.DeleteAccount(id)
}

func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
