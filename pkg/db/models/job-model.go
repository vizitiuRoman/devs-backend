package models

import (
	"time"
)

type JobModel interface {
	BeforeSave() error
	Prepare()
	Validate(action string) error
}

type Job struct {
	ID        uint32    `gorm:"not null;auto_increment" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Hooks

func (job *Job) BeforeSave() error {
	return nil
}

// Prepare/Validate

func (job *Job) Prepare() {
	job.ID = 0
	job.CreatedAt = time.Now()
	job.UpdateAt = time.Now()
}

func (job *Job) Validate(action string) error {

	return nil
}

// Query to database
