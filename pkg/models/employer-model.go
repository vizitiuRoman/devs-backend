package models

import (
	"time"
)

type EmployerModel interface {
	BeforeSave() error
	Prepare()
	Validate(action string) error
}

type Employer struct {
	ID        uint32    `gorm:"not null;auto_increment" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Hooks

func (employer *Employer) BeforeSave() error {
	return nil
}

// Prepare/Validate

func (employer *Employer) Prepare() {
	employer.ID = 0
	employer.CreatedAt = time.Now()
	employer.UpdateAt = time.Now()
}

func (employer *Employer) Validate(action string) error {

	return nil
}

// Query to database
