package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Users struct {
	//id,created_at,updated at insert themselves automaticaly when gorm.Model is used
	gorm.Model
	Name     string `json:"name" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Contact  string `json:"contact" gorm:"unique"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Avatar   string `json:"avatar"`
}

func (user *Users) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *Users) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
