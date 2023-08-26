package segment

import (
	"market/pkg/errors"
	"market/pkg/storage/general"
	"market/pkg/storage/segment"
	"time"
)

type IRepository interface {
	CreateSegment(name string, percent int) errors.CustomError
	DeleteSegment(name string) errors.CustomError
	DeleteExpiredSegments(moment *time.Time) errors.CustomError
}

type Repository struct {
	segStorage segment.IStorage
	genStorage general.IStorage
}

func NewRepository(segStorage segment.IStorage, genStorage general.IStorage) IRepository {
	return &Repository{segStorage, genStorage}
}

func (r *Repository) CreateSegment(name string, percent int) errors.CustomError {
	_, err := r.segStorage.GetSegmentByName(name)
	if err == nil {
		return errors.SegmentAlreadyExist{}
	} else {
		if err.Error() != errors.SegmentNotFoundErr400 {
			return errors.GetSegmentByName{Err: err.Error()}
		}
	}

	segmentId, err := r.segStorage.CreateSegment(name)
	if err != nil {
		return errors.CreateSegment{Err: err.Error()}
	}

	totalUsers, err := r.genStorage.CountUsersNumber()
	if err != nil {
		return errors.CountUsersNumber{Err: err.Error()}
	}

	logs, err := r.segStorage.AddSegmentsToUsersByPercent(totalUsers, segmentId, percent)
	if err != nil {
		return errors.AddPercent{Err: err.Error()}
	}

	if len(logs) > 0 {
		err = r.genStorage.AddLogs(logs)
		if err != nil {
			return errors.AddLogs{Err: err.Error()}
		}
	}

	return nil
}

func (r *Repository) DeleteSegment(name string) errors.CustomError {
	_, err := r.segStorage.GetSegmentByName(name)
	if err != nil {
		return errors.GetSegmentByName{Err: err.Error()}
	}
	err = r.segStorage.DeleteSegment(name)
	if err != nil {
		return errors.DeleteSegments{Err: err.Error()}
	}
	return nil
}

func (r *Repository) DeleteExpiredSegments(moment *time.Time) errors.CustomError {
	logs, err := r.segStorage.DeleteExpiredSegments(moment)
	if err != nil {
		return errors.DeleteSegments{Err: err.Error()}
	}

	if len(logs) > 0 {
		err = r.genStorage.AddLogs(logs)
		if err != nil {
			return errors.AddLogs{Err: err.Error()}
		}
	}

	return nil
}
