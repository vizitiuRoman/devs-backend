package models

import (
	"time"
)

type EmployerModel interface {
	BeforeSave()
	BeforeCreate()
	AfterSave()
	AfterCreate()
}

type Employer struct {
	ID        uint32    `gorm:"not nul;auto_increment" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Hooks https://gorm.io/docs/hooks.html

func (employer *Employer) BeforeSave() {}

func (employer *Employer) BeforeCreate() {}

func (employer *Employer) AfterSave() {}

func (employer *Employer) AfterCreate() {}

// Prepare/Validate

func (employer *Employer) Prepare() {
	employer.ID = 0
	employer.CreatedAt = time.Now()
	employer.UpdateAt = time.Now()
}

func (employer *Employer) Validate() error {

	return nil
}

// Query to database