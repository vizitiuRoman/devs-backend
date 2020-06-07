package models

import (
	"time"
)

type StoryModel interface {
	BeforeSave()
	BeforeCreate()
	AfterSave()
	AfterCreate()
	Prepare()
	Validate() error
}

type Story struct {
	ID        uint32    `gorm:"not nul;auto_increment" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Hooks https://gorm.io/docs/hooks.html

func (story *Story) BeforeSave() {}

func (story *Story) BeforeCreate() {}

func (story *Story) AfterSave() {}

func (story *Story) AfterCreate() {}

// Prepare/Validate

func (story *Story) Prepare() {
	story.ID = 0
	story.CreatedAt = time.Now()
	story.UpdateAt = time.Now()
}

func (story *Story) Validate() error {

	return nil
}

// Query to database
