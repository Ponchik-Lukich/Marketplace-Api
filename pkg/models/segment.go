package models

import "gorm.io/gorm"

type Segment struct {
	gorm.Model

	ID    uint64 `gorm:"primary_key;auto_increment"`
	Name  string
	Users []User `gorm:"many2many:user_segments;"`
}
