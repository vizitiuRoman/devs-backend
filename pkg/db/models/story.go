package models

import (
	"time"
)

type StoryModel interface {
	BeforeSave() error
	Prepare()
	Validate(action string) error
}

type Story struct {
	ID        uint32    `gorm:"not null;auto_increment" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Hooks

func (story *Story) BeforeSave() error {
	return nil
}

// Prepare/Validate

func (story *Story) Prepare() {
	story.ID = 0
	story.CreatedAt = time.Now()
	story.UpdateAt = time.Now()
}

func (story *Story) Validate(action string) error {

	return nil
}

// Query to database
