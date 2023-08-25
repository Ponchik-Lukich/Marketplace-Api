package segment

import (
	"fmt"
	"market/pkg/errors"
	"market/pkg/storage"
	"time"
)

type IRepository interface {
	CreateSegment(name string) error
	DeleteSegment(name string) error
	DeleteExpiredSegments(moment *time.Time) error
}

type Repository struct {
	storage storage.IStorage
}

func NewRepository(storage storage.IStorage) IRepository {
	return &Repository{storage}
}

func (r *Repository) CreateSegment(name string) error {
	_, err := r.storage.GetSegmentByName(name)
	if err == nil {
		return fmt.Errorf(errors.SegmentAlreadyExist)
	} else {
		if err.Error() != errors.SegmentNotFoundErr {
			return err
		}
	}
	return r.storage.CreateSegment(name)
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
