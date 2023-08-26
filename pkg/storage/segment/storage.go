package segment

import (
	"market/pkg/models"
	"time"
)

type IStorage interface {
	CreateSegment(name string) (uint64, error)
	DeleteSegment(name string) error
	GetSegmentByName(name string) (models.Segment, error)
	DeleteExpiredSegments(moment *time.Time) ([]models.Log, error)
	AddSegmentsToUsersByPercent(uint64, uint64, int) ([]models.Log, error)
}
