package service

import (
	"errors"
	"fmt"

	"anyrouter-checkin/internal/model"
	"anyrouter-checkin/internal/repository"

	"gorm.io/gorm"
)

var ErrInvalidSession = errors.New("session 无效")
var ErrAccountDisabled = errors.New("账号已禁用")

func ListAccounts() ([]model.Account, error) {
	return repository.ListAccounts()
}

func CreateAccount(session string) (model.Account, error) {
	info, err := ParseSession(session)
	if err != nil {
		return model.Account{}, fmt.Errorf("%w: %v", ErrInvalidSession, err)
	}

	account := model.Account{
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

func UpdateAccount(id uint, session string) (model.Account, error) {
	account, err := repository.GetAccountByID(id)
	if err != nil {
		return model.Account{}, err
	}

	if session == "" {
		return *account, nil
	}
	info, err := ParseSession(session)
	if err != nil {
		return model.Account{}, fmt.Errorf("%w: %v", ErrInvalidSession, err)
	}
	selfInfo, err := fetchAccountSelf(session, info.UserID)
	if err != nil {
		return model.Account{}, fmt.Errorf("获取账号信息失败: %w", err)
	}
	account.Session = session
	account.UserID = selfInfo.UserID
	account.Username = selfInfo.Username
	account.Role = selfInfo.Role
	account.Balance = selfInfo.Balance
	if err := repository.SaveAccount(account); err != nil {
		return model.Account{}, err
	}

	return *account, nil
}

func RefreshAccount(id uint) (model.Account, error) {
	account, err := repository.GetAccountByID(id)
	if err != nil {
		return model.Account{}, err
	}
	if account.Status != 1 {
		return model.Account{}, ErrAccountDisabled
	}

	sessionInfo, err := ParseSession(account.Session)
	if err != nil {
		return model.Account{}, fmt.Errorf("%w: %v", ErrInvalidSession, err)
	}

	info, err := fetchAccountSelf(account.Session, sessionInfo.UserID)
	if err != nil {
		return model.Account{}, fmt.Errorf("获取账号信息失败: %w", err)
	}

	account.UserID = info.UserID
	account.Username = info.Username
	account.Role = info.Role
	account.Balance = info.Balance

	if err := repository.SaveAccount(account); err != nil {
		return model.Account{}, err
	}

	return *account, nil
}

func UpdateAccountStatus(id uint, status int) (model.Account, error) {
	account, err := repository.GetAccountByID(id)
	if err != nil {
		return model.Account{}, err
	}

	account.Status = status
	if err := repository.SaveAccount(account); err != nil {
		return model.Account{}, err
	}

	return *account, nil
}

func DeleteAccount(id uint) error {
	if err := removeAccountFromCronTasks(id); err != nil {
		return err
	}
	return repository.DeleteAccount(id)
}

func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
