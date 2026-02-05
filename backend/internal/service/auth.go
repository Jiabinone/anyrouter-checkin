package service

import (
	"errors"

	"anyrouter-checkin/internal/config"
	"anyrouter-checkin/internal/middleware"
	"anyrouter-checkin/internal/model"
	"anyrouter-checkin/internal/repository"

	"github.com/dromara/carbon/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(username, password string) (string, error) {
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("密码错误")
	}

	return generateToken(user.ID, user.Username)
}

func generateToken(userID uint, username string) (string, error) {
	now := carbon.Now().StdTime()
	claims := middleware.Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(config.C.JWT.Expire)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.C.JWT.Secret))
}

func InitAdminUser() error {
	count, err := repository.CountUsers()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(config.C.Admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return repository.CreateUser(&model.User{
		Username: config.C.Admin.Username,
		Password: string(hashed),
	})
}

func ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user.Password = string(hashed)
	return repository.SaveUser(user)
}
