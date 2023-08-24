package storage

import (
	"gorm.io/gorm"
)

type IStorage interface {
	Connect() error
	Close() error
	Init() *gorm.DB
	MakeMigrations() error
}
