package models

import (
	"gorm.io/gorm"
	"time"
)

type UserSegment struct {
	gorm.Model

	UserID         uint64
	SegmentID      uint64
	ExpirationDate time.Time `gorm:"default:null"`
}
