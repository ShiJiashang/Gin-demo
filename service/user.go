package service

import (
	"errors"
	"time"

	"gin_gorm_demo/auth"
	"gin_gorm_demo/dto"
	"gin_gorm_demo/model"
	"gin_gorm_demo/repository"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrNameExists      = errors.New("name already exists")
	ErrInvalidAge      = errors.New("invalid age")
	ErrInvalidPassword = errors.New("invalid password")
)

func ListUsers() ([]dto.UserResponse, error) {
	users, err := repository.ListUsers()
	if err != nil {
		return nil, err
	}

	result := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		result = append(result, toUserResponse(user))
	}

	return result, nil
}

func GetUserByID(id uint) (dto.UserResponse, error) {
	user, err := repository.GetUserByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return dto.UserResponse{}, ErrUserNotFound
		}

		return dto.UserResponse{}, err
	}

	return toUserResponse(user), nil
}

func CreateUser(req dto.CreateUserRequest) (dto.UserResponse, error) {
	if req.Age <= 0 {
		return dto.UserResponse{}, ErrInvalidAge
	}

	exists, err := repository.IsNameUsedByOtherUser(req.Name, 0)
	if err != nil {
		return dto.UserResponse{}, err
	}
	if exists {
		return dto.UserResponse{}, ErrNameExists
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return dto.UserResponse{}, err
	}

	user := model.User{
		Name:     req.Name,
		Age:      req.Age,
		Password: hashedPassword,
		Role:     "user",
	}

	user, err = repository.CreateUser(user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return toUserResponse(user), nil
}

func UpdateUser(id uint, req dto.UpdateUserRequest) (dto.UserResponse, error) {
	if req.Age <= 0 {
		return dto.UserResponse{}, ErrInvalidAge
	}

	user, err := repository.GetUserByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return dto.UserResponse{}, ErrUserNotFound
		}

		return dto.UserResponse{}, err
	}

	exists, err := repository.IsNameUsedByOtherUser(req.Name, id)
	if err != nil {
		return dto.UserResponse{}, err
	}
	if exists {
		return dto.UserResponse{}, ErrNameExists
	}

	user.Name = req.Name
	user.Age = req.Age

	user, err = repository.UpdateUser(user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return toUserResponse(user), nil
}

func DeleteUser(id uint) error {
	err := repository.DeleteUser(id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return ErrUserNotFound
		}

		return err
	}

	return nil
}

func toUserResponse(user model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Age:  user.Age,
		Role: user.Role,
	}
}

// 登陆：根据用户名和密码验证用户是否存在，密码是否正确，返回token
func Login(req dto.LoginRequest, secret string, duration time.Duration) (dto.LoginResponse, error) {

	user, err := repository.GetUserByName(req.Name)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return dto.LoginResponse{}, ErrUserNotFound
		}

		return dto.LoginResponse{}, err
	}
	if !auth.CheckPassword(req.Password, user.Password) {
		return dto.LoginResponse{}, ErrInvalidPassword
	}

	token, err := auth.GenerateToken(user.ID, secret, duration)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	return dto.LoginResponse{
		Token: token,
	}, nil
}
