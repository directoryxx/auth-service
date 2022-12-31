package usecase

import (
	"auth/internal/domain"
	"auth/internal/helper"
	"auth/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// UserUseCase represent the user's usecase contract
type UserUseCase interface {
	Login(ctx context.Context, login *domain.LoginValidation) (user *domain.User, uuid string, err error)
	Register(ctx context.Context, register *domain.RegisterValidation) (user *domain.User, uuid string, err error)
	Profile(ctx context.Context, uuid string) (user *domain.User, err error)
	CheckUsername(ctx context.Context, username string) (user *domain.User, err error)
}

type UserUseCaseImpl struct {
	UserRepo repository.UserRepository
}

// NewMysqlAuthorRepository will create an implementation of author.Repository
func NewUserUseCase(UserRepo repository.UserRepository) UserUseCase {
	return &UserUseCaseImpl{
		UserRepo: UserRepo,
	}
}

func (uc *UserUseCaseImpl) Login(ctx context.Context, login *domain.LoginValidation) (user *domain.User, uuidGen string, err error) {
	usernameCheck, _ := uc.UserRepo.GetOneByUsername(ctx, login.Username)

	if usernameCheck == nil {
		return nil, "", errors.New("username / password salah")
	}

	passwordCheck, _ := helper.ComparePasswordAndHash(login.Password, usernameCheck.Password)

	if !passwordCheck {
		return nil, "", errors.New("username / password salah")
	}

	uuidGenerate := uuid.NewString()

	rememberSession := uc.UserRepo.RememberUUID(ctx, usernameCheck, uuidGenerate)

	if rememberSession != nil {
		return nil, "", rememberSession
	}

	// uc.UserRepo.Publish(ctx, "test")

	return usernameCheck, uuidGenerate, nil
}

func (uc *UserUseCaseImpl) Register(context context.Context, register *domain.RegisterValidation) (user *domain.User, uuidGen string, err error) {
	_, err = uc.CheckUsername(context, register.Username)

	if err == nil {
		return nil, "", errors.New("username telah terdaftar")
	}

	hashpassword, err := helper.CreateHash(register.Password, helper.DefaultParams)

	if err != nil {
		return nil, "", err
	}

	userInput := &domain.User{
		Email:    register.Email,
		Name:     register.Name,
		Username: register.Username,
		Password: hashpassword,
	}

	user, err = uc.UserRepo.Insert(context, userInput)

	uuidGenerate := uuid.NewString()

	if err != nil {
		return nil, "", err
	}

	mail := &domain.Message{
		To:      user.Email,
		From:    "admin@email.com",
		Subject: user.Username + ", Your account is registered",
		Data:    "Hi, " + userInput.Name + ". Your account is registered. Please Login",
		Uuid:    uuidGenerate,
	}

	b, _ := json.Marshal(mail)

	uc.UserRepo.Publish(context, string(b), "mail")

	rememberSession := uc.UserRepo.RememberUUID(context, user, uuidGenerate)

	if rememberSession != nil {
		return nil, "", err
	}

	fmt.Println(user)

	return user, uuidGenerate, nil
}

func (uc *UserUseCaseImpl) Profile(ctx context.Context, uuid string) (user *domain.User, err error) {
	userMod := &domain.User{}
	res, err := uc.UserRepo.GetUUID(ctx, uuid)

	if err != nil {
		return nil, err
	}

	var jsonData = []byte(res)

	var _ = json.Unmarshal(jsonData, &userMod)

	return userMod, nil
}

func (uc *UserUseCaseImpl) CheckUsername(ctx context.Context, username string) (user *domain.User, err error) {
	user, err = uc.UserRepo.GetOneByUsername(ctx, username)

	if err != nil {
		return nil, err
	}

	return user, nil
}
