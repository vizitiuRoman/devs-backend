package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	. "github.com/devs-backend/user-service/pkg/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserModel interface {
	BeforeSave() error
	Prepare()
	Validate(action string) error
	Create() (*User, error)
	Update() (*User, error)
	DeleteByID(AccessUUID, RefreshUUID string) error
	FindByID() (*User, error)
	FindByEmail() (*User, error)
	FindAll() (*[]User, error)
}

type User struct {
	ID        uint32    `gorm:"not nul;auto_increment" json:"id"`
	Email     string    `gorm:"not null;unique;" json:"email"`
	Password  string    `gorm:"not null;" json:"password"`
	Name      string    `gorm:"not null;" json:"name"`
	LastName  string    `gorm:"not null;" json:"lastName"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Hooks

func (user *User) BeforeSave() error {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

// Prepare/Validate/VerifyPassword/HashPassword

func (user *User) Prepare() {
	user.ID = 0
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.Password = html.EscapeString(strings.TrimSpace(user.Password))
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.LastName = html.EscapeString(strings.TrimSpace(user.LastName))
	user.CreatedAt = time.Now()
	user.UpdateAt = time.Now()
}

func (user *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case LOGIN:
		if user.Email == "" {
			return errors.New("Required Email")
		}
		if user.Password == "" {
			return errors.New("Required Password")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	case UPDATE:
		if user.Email == "" {
			return errors.New("Required Email")
		}
		if user.Name == "" {
			return errors.New("Required Name")
		}
		if user.LastName == "" {
			return errors.New("Required Last Name")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	default:
		if user.Email == "" {
			return errors.New("Required Email")
		}
		if user.Name == "" {
			return errors.New("Required Name")
		}
		if user.LastName == "" {
			return errors.New("Required Last Name")
		}
		if user.Password == "" {
			return errors.New("Required Password")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// Query to database

func (user *User) Create() (*User, error) {
	err := DB.Model(&User{}).Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) Update() (*User, error) {
	err := user.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Model(&User{}).Where("id = ?", user.ID).Update(&User{
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		LastName: user.LastName,
		UpdateAt: time.Now(),
	}).Take(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) DeleteByID(AccessUUID, RefreshUUID string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&User{}).Where("id = ?", user.ID).Take(&user).Delete(&user).Error
		if err != nil {
			return err
		}
		token := TokenDetails{
			AccessUUID:  AccessUUID,
			RefreshUUID: RefreshUUID,
		}
		err = token.DeleteByUUID()
		if err != nil {
			return err
		}
		return nil
	})
}

func (user *User) FindByID() (*User, error) {
	err := DB.Model(&User{}).Where("id = ?", user.ID).Take(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) FindByEmail() (*User, error) {
	err := DB.Model(&User{}).Where("email = ?", user.Email).Take(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) FindAll() (*[]User, error) {
	var users []User
	err := DB.Model(&[]User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, nil
}
