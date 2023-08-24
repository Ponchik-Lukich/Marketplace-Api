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
	GetSegmentByName(name string) (models.Segment, error)
	GetIDsAndMissingNames(names []string) ([]uint64, []string, error)
	AddSegmentsToUser(toCreate []uint64, toDelete []uint64, userID uint64) error
}
