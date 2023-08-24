package models

type UserSegment struct {
	UserID    uint64 `gorm:"primary_key;auto_increment"`
	SegmentID uint64 `gorm:"primary_key;auto_increment"`
}
