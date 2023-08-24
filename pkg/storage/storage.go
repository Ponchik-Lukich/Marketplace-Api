package storage

import (
	"gorm.io/gorm"
	"market/pkg/models"
)

type IStorage interface {
	Connect() error
	Close() error
	Init() *gorm.DB
	MakeMigrations() error

	CreateSegment(name string) error
	DeleteSegment(name string) error
	GetSegmentsByUserID(userID uint64) ([]models.Segment, error)
	AddSegmentsToUser(toCreate []string, toDelete []string, userID uint64) []error
	GetSegmentByName(name string) (models.Segment, error)
}
