package models

import (
	"time"
)

type ProjectModel interface {
	BeforeSave() error
	Prepare()
	Validate(action string) error
}

type Project struct {
	ID        uint32    `gorm:"not null;auto_increment" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Hooks

func (project *Project) BeforeSave() error {
	return nil
}

// Prepare/Validate

func (project *Project) Prepare() {
	project.ID = 0
	project.CreatedAt = time.Now()
	project.UpdateAt = time.Now()
}

func (project *Project) Validate(action string) error {

	return nil
}

// Query to database
