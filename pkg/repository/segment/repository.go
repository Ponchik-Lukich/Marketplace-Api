package segment

import (
	"market/pkg/errors"
	"market/pkg/storage"
	"time"
)

type IRepository interface {
	CreateSegment(name string, percent int) errors.CustomError
	DeleteSegment(name string) errors.CustomError
	DeleteExpiredSegments(moment *time.Time) errors.CustomError
}

type Repository struct {
	storage storage.IStorage
}

func NewRepository(storage storage.IStorage) IRepository {
	return &Repository{storage}
}

func (r *Repository) CreateSegment(name string, percent int) errors.CustomError {
	_, err := r.storage.GetSegmentByName(name)
	if err == nil {
		return errors.SegmentAlreadyExist{}
	} else {
		if err.Error() != errors.SegmentNotFoundErr400 {
			return errors.GetSegmentByName{Err: err.Error()}
		}
	}

	segmentId, err := r.storage.CreateSegment(name)
	if err != nil {
		return errors.CreateSegment{Err: err.Error()}
	}

	totalUsers, err := r.storage.CountUsersNumber()
	if err != nil {
		return errors.CountUsersNumber{Err: err.Error()}
	}

	err = r.storage.AddSegmentsToUsersByPercent(totalUsers, segmentId, percent)
	if err != nil {
		return errors.AddPercent{Err: err.Error()}
	}

	return nil
}

func (r *Repository) DeleteSegment(name string) errors.CustomError {
	_, err := r.storage.GetSegmentByName(name)
	if err != nil {
		return errors.GetSegmentByName{Err: err.Error()}
	}
	err = r.storage.DeleteSegment(name)
	if err != nil {
		return errors.DeleteSegments{Err: err.Error()}
	}
	return nil
}

func (r *Repository) DeleteExpiredSegments(moment *time.Time) errors.CustomError {
	logs, err := r.storage.DeleteExpiredSegments(moment)
	if err != nil {
		return errors.DeleteSegments{Err: err.Error()}
	}

	if len(logs) > 0 {
		err = r.storage.AddLogs(logs)
		if err != nil {
			return errors.AddLogs{Err: err.Error()}
		}
	}

	return nil
}
