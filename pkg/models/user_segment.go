package models

import "gorm.io/gorm"

type UserSegment struct {
	gorm.Model

	UserID    uint64
	SegmentID uint64
}
