package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password string `gorm:"type:varchar(100) ;not null" json:"-"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
}

func (u *User) CheckPassword(rawPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(rawPassword),
	)
	return err == nil
}
