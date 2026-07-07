package repository

import (
	"errors"

	"gin_gorm_demo/database"
	"gin_gorm_demo/model"

	"gorm.io/gorm"
)

var ErrRecordNotFound = errors.New("record not found")

func ListUsers() ([]model.User, error) {
	var users []model.User

	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByID(id uint) (model.User, error) {
	var user model.User

	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, ErrRecordNotFound
		}

		return model.User{}, err
	}

	return user, nil
}

func CreateUser(user model.User) (model.User, error) {
	if err := database.DB.Create(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}

func UpdateUser(user model.User) (model.User, error) {
	if err := database.DB.Save(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}

// 按照用户名查找user
func GetUserByName(name string) (model.User, error) {
	var user model.User

	if err := database.DB.Where("name=?", name).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, ErrRecordNotFound
		}
		return model.User{}, err
	}
	return user, nil
}
func DeleteUser(id uint) error {
	user, err := GetUserByID(id)
	if err != nil {
		return err
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

func IsNameUsedByOtherUser(name string, currentUserID uint) (bool, error) {
	var count int64

	err := database.DB.Model(&model.User{}).
		Where("name = ? AND id <> ?", name, currentUserID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
