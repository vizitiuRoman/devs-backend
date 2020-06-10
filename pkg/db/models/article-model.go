package models

import (
	"time"
)

type ArticleModel interface {
	BeforeSave() error
	Prepare()
	Validate(action string) error
}

type Article struct {
	ID        uint32    `gorm:"not null;auto_increment" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Hooks

func (article *Article) BeforeSave() error {
	return nil
}

// Prepare/Validate

func (article *Article) Prepare() {
	article.ID = 0
	article.CreatedAt = time.Now()
	article.UpdateAt = time.Now()
}

func (article *Article) Validate(action string) error {

	return nil
}

// Query to database
