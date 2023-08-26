package general

import (
	"gorm.io/gorm"
	"market/pkg/models"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Init() *gorm.DB {
	return s.db
}

func (s *Storage) CountUsersNumber() (uint64, error) {
	db := s.Init()

	var count int64
	if err := db.Table("users").Count(&count).Error; err != nil {
		return 0, err
	}

	return uint64(count), nil
}

func (s *Storage) AddLogs(logs []models.Log) error {
	db := s.Init()
	return db.Create(&logs).Error
}
