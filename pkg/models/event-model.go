package models

import (
	"time"
)

type EventModel interface {
	BeforeSave()
	BeforeCreate()
	AfterSave()
	AfterCreate()
	Prepare()
	Validate() error
}

type Event struct {
	ID        uint32    `gorm:"not nul;auto_increment" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Hooks https://gorm.io/docs/hooks.html

func (event *Event) BeforeSave() {}

func (event *Event) BeforeCreate() {}

func (event *Event) AfterSave() {}

func (event *Event) AfterCreate() {}

// Prepare/Validate

func (event *Event) Prepare() {
	event.ID = 0
	event.CreatedAt = time.Now()
	event.UpdateAt = time.Now()
}

func (event *Event) Validate() error {

	return nil
}

// Query to database
