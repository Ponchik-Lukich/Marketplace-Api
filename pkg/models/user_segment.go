package models

import "gorm.io/gorm"

type UserSegment struct {
	gorm.Model

	UserID    uint64 `gorm:"primary_key;auto_increment"`
	User      User
	SegmentID uint64 `gorm:"primary_key;auto_increment"`
	Segment   Segment
}
