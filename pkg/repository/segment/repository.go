package segment

import (
	"fmt"
	"market/pkg/errors"
	"market/pkg/storage"
)

type IRepository interface {
	CreateSegment(name string) error
	DeleteSegment(name string) error
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
