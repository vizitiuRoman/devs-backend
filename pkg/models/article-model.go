package models

import (
	"time"
)

type ArticleModel interface {
	BeforeSave()
	BeforeCreate()
	AfterSave()
	AfterCreate()
}

type Article struct {
	ID        uint32    `gorm:"not nul;auto_increment" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Hooks https://gorm.io/docs/hooks.html

func (article *Article) BeforeSave() {}

func (article *Article) BeforeCreate() {}

func (article *Article) AfterSave() {}

func (article *Article) AfterCreate() {}

// Prepare/Validate

func (article *Article) Prepare() {
	article.ID = 0
	article.CreatedAt = time.Now()
	article.UpdateAt = time.Now()
}

func (article *Article) Validate() error {

	return nil
}

// Query to database