package models

import (
	"time"
)

type EventModel interface {
	BeforeSave() error
	Prepare()
	Validate(action string) error
}

type Event struct {
	ID        uint32    `gorm:"not null;auto_increment" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Hooks

func (event *Event) BeforeSave() error {
	return nil
}

// Prepare/Validate

func (event *Event) Prepare() {
	event.ID = 0
	event.CreatedAt = time.Now()
	event.UpdateAt = time.Now()
}

func (event *Event) Validate(action string) error {

	return nil
}

// Query to database
