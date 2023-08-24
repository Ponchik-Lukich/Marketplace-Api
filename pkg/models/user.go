package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	ID       uint64        `gorm:"primary_key;auto_increment"`
	Segments []UserSegment `gorm:"foreignKey:UserID"`
}
