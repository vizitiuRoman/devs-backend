package models

import (
	"html"
	"strings"
	"time"

	. "github.com/devsmd/pkg/db"
)

type TokenModel interface {
	Prepare()
	CreateOrUpdate(userID uint32) (*Token, error)
	DeleteById(userID uint32) (int64, error)
}

type Token struct {
	ID        uint32    `gorm:"not nul;auto_increment" json:"id"`
	UserID    uint32    `gorm:"unique" json:"user_id"`
	Token     string    `gorm:"size:255" json:"token"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Prepare

func (token *Token) Prepare() {
	token.ID = 0
	token.Token = html.EscapeString(strings.TrimSpace(token.Token))
	token.CreatedAt = time.Now()
	token.UpdateAt = time.Now()
}

// Query to database

func (token *Token) CreateOrUpdate(userID uint32) (*Token, error) {
	err := DB.Debug().Model(&Token{}).Where("user_id = ?", userID).Update(&token).Take(&token).Error
	if err == nil {
		return token, nil
	}

	err = DB.Debug().Model(&Token{}).Create(&token).Error
	if err == nil {
		return &Token{}, err
	}
	return token, err
}

func (token *Token) DeleteById(userID uint32) (int64, error) {
	err := DB.Debug().Model(&Token{}).Where("user_id = ?", userID).Delete(&token)
	if err != nil {
		return 0, err.Error
	}
	return err.RowsAffected, nil
}
