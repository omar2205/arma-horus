package models

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key;"`
	Email   string    // `gorm:"unique"` ? user can sign in multiple devices
	Token   string    `gorm:"unique"`
	Expires time.Time
}

func NewUser(email, token string) *User {
	return &User{
		ID:      uuid.NewV4(),
		Email:   email,
		Token:   token,
		Expires: time.Now().Add(time.Minute * 30), // expires in 30mins
	}
}

var sqlite_file = sqlite.Open("./db.sqlite")
var gorm_config = &gorm.Config{}

func (u *User) Save() error {
	db, err := gorm.Open(sqlite_file, gorm_config)
	if err != nil {
		return err
	}

	// Auto migrate the users table if it does not exist
	err = db.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	err = db.Create(u).Error
	if err != nil {
		return err
	}

	return nil
}

func ValidateUserToken(token string) (*User, error) {
	db, err := gorm.Open(sqlite_file, gorm_config)
	if err != nil {
		return nil, err
	}

	var user User

	err = db.Where("token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}

	if time.Now().After(user.Expires) {
		err = db.Where("token = ?", token).Delete(&user).Error
		return nil, errors.New("Token expired")
	}

	return &user, nil
}
