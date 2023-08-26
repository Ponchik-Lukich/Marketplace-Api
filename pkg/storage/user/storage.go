package user

import (
	"market/pkg/dtos"
	"market/pkg/errors"
	"market/pkg/models"
	"time"
)

type IStorage interface {
	CreateUser(userID uint64) error
	GetIDsAndMissingNames(names []string) ([]uint64, []string, error)
	GetUserLogs(start *time.Time, end *time.Time, userID uint64) ([]models.Log, error)
	GetUserSegmentsByID(userID uint64) ([]models.Segment, error)
	AddSegmentsToUser(toCreate []uint64, toCreateDto []dtos.CreateSegmentDto, userID uint64) ([]models.Log, errors.CustomError)
	DeleteSegmentsFromUser(toDelete []uint64, toDeleteDto []string, userID uint64) ([]models.Log, errors.CustomError)
}
