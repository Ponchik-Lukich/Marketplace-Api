package segment

import (
	"fmt"
	"market/pkg/errors"
	"market/pkg/storage"
	"time"
)

type IRepository interface {
	CreateSegment(name string, percent int) error
	DeleteSegment(name string) error
	DeleteExpiredSegments(moment *time.Time) error
}

type Repository struct {
	storage storage.IStorage
}

func NewRepository(storage storage.IStorage) IRepository {
	return &Repository{storage}
}

func (r *Repository) CreateSegment(name string, percent int) error {
	_, err := r.storage.GetSegmentByName(name)
	if err == nil {
		return fmt.Errorf(errors.SegmentAlreadyExist)
	} else {
		if err.Error() != errors.SegmentNotFoundErr {
			return fmt.Errorf("%s: %w", errors.GetSegmentByNameErr, err)
		}
	}

	segmentId, err := r.storage.CreateSegment(name)
	if err != nil {
		return err
	}

	totalUsers, err := r.storage.CountUsersNumber()
	if err != nil {
		return fmt.Errorf("%s: %w", errors.CountUsersNumberErr, err)
	}

	err = r.storage.AddSegmentsToUsersByPercent(totalUsers, segmentId, percent)
	if err != nil {
		return fmt.Errorf("%s: %w", errors.AddingPercentErr, err)
	}

	return nil
}

func (r *Repository) DeleteSegment(name string) error {
	_, err := r.storage.GetSegmentByName(name)
	if err != nil {
		return err
	}
	return r.storage.DeleteSegment(name)
}

func (r *Repository) DeleteExpiredSegments(moment *time.Time) error {
	logs, err := r.storage.DeleteExpiredSegments(moment)
	if err != nil {
		return err
	}

	if len(logs) > 0 {
		err = r.storage.AddLogs(logs)
		if err != nil {
			return err
		}
	}

	return nil
}
