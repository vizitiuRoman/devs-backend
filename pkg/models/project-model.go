package models

import (
	"time"
)

type ProjectModel interface {
	BeforeSave()
	BeforeCreate()
	AfterSave()
	AfterCreate()
	Prepare()
	Validate() error
}

type Project struct {
	ID        uint32    `gorm:"not nul;auto_increment" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Hooks https://gorm.io/docs/hooks.html

func (project *Project) BeforeSave() {}

func (project *Project) BeforeCreate() {}

func (project *Project) AfterSave() {}

func (project *Project) AfterCreate() {}

// Prepare/Validate

func (project *Project) Prepare() {
	project.ID = 0
	project.CreatedAt = time.Now()
	project.UpdateAt = time.Now()
}

func (project *Project) Validate() error {

	return nil
}

// Query to database
