package models

import (
	"time"
)

type JobModel interface {
	BeforeSave()
	BeforeCreate()
	AfterSave()
	AfterCreate()
}

type Job struct {
	ID        uint32    `gorm:"not nul;auto_increment" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Hooks https://gorm.io/docs/hooks.html

func (job *Job) BeforeSave() {}

func (job *Job) BeforeCreate() {}

func (job *Job) AfterSave() {}

func (job *Job) AfterCreate() {}

// Prepare/Validate

func (job *Job) Prepare() {
	job.ID = 0
	job.CreatedAt = time.Now()
	job.UpdateAt = time.Now()
}

func (job *Job) Validate() error {

	return nil
}

// Query to database
