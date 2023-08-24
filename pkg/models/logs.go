package models

import (
	"gorm.io/gorm"
	"time"
)

type Log struct {
	gorm.Model
	ID        uint64     `gorm:"primaryKey;autoIncrement"`
	UserID    uint64     `gorm:"not null"`
	Segment   string     `gorm:"not null"`
	EventType string     `gorm:"not null"`
	Time      *time.Time `gorm:"not null"`
}
