package models

import (
	"errors"
	"time"
)

type UserModel interface {
	BeforeSave()
	BeforeCreate()
	AfterSave()
	AfterCreate()
	Prepare()
	Validate() error
}

type User struct {
	ID        uint32    `gorm:"not nul;auto_increment" json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Hooks https://gorm.io/docs/hooks.html

func (user *User) BeforeSave() {}

func (user *User) BeforeCreate() {}

func (user *User) AfterSave() {}

func (user *User) AfterCreate() {}

// Prepare/Validate

func (user *User) Prepare() {
	user.ID = 0
	user.CreatedAt = time.Now()
	user.UpdateAt = time.Now()
}

func (user *User) Validate() error {
	if user.Email == "" {
		return errors.New("REQUIRE EMAIL")
	}
	return nil
}

// Query to database
