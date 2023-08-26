package storage

import (
	"gorm.io/gorm"
	"market/pkg/dtos"
	"market/pkg/errors"
	"market/pkg/models"
	"time"
)

type IStorage interface {
	Connect() error
	Close() error
	Init() *gorm.DB
	MakeMigrations() error

	CreateSegment(name string) (uint64, error)
	DeleteSegment(name string) error
	GetSegmentsByUserID(userID uint64) ([]models.Segment, error)
	GetSegmentByName(name string) (models.Segment, error)
	DeleteExpiredSegments(moment *time.Time) ([]models.Log, error)

	CreateUser(userID uint64) error
	GetIDsAndMissingNames(names []string) ([]uint64, []string, error)
	GetUserLogs(start *time.Time, end *time.Time, userID uint64) ([]models.Log, error)
	CountUsersNumber() (uint64, error)
	AddSegmentsToUsersByPercent(uint64, uint64, int) error
	AddSegmentsToUser(toCreate []uint64, toCreateDto []dtos.CreateSegmentDto, userID uint64) ([]models.Log, errors.CustomError)
	DeleteSegmentsFromUser(toDelete []uint64, toDeleteDto []string, userID uint64) ([]models.Log, errors.CustomError)

	AddLogs(logs []models.Log) error
}
