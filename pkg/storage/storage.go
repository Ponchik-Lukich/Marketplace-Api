package storage

import (
	"gorm.io/gorm"
	"market/pkg/storage/general"
	"market/pkg/storage/segment"
	"market/pkg/storage/user"
)

type IStorage interface {
	Connect() error
	Close() error
	Init() *gorm.DB
	MakeMigrations() error

	GetSegmentStorage() segment.IStorage
	GetUserStorage() user.IStorage
	GetGeneralStorage() general.IStorage
}
