package user

import (
	"fmt"
	"market/pkg/dtos"
	"market/pkg/storage"
)

type IRepository interface {
	EditUsersSegments(toCreate []string, toDelete []string, userID uint64) error
	GetUsersSegments(userID uint64) ([]dtos.SegmentDtoResponse, error)
}

type Repository struct {
	storage storage.IStorage
}

func NewRepository(storage storage.IStorage) IRepository {
	return &Repository{storage}
}

func (r *Repository) EditUsersSegments(toCreate []string, toDelete []string, userID uint64) error {
	toCreateIds, toCreateMissingNames, err := r.storage.GetIDsAndMissingNames(toCreate)
	if err != nil {
		return err
	}
	toDeleteIds, toDeleteMissingNames, err := r.storage.GetIDsAndMissingNames(toDelete)
	if err != nil {
		return err
	}

	var missingNames []string

	if len(toCreateMissingNames) > 0 {
		for _, name := range toCreateMissingNames {
			missingNames = append(missingNames, name)
		}
	}

	if len(toDeleteMissingNames) > 0 {
		for _, name := range toDeleteMissingNames {
			missingNames = append(missingNames, name)
		}
	}

	if len(missingNames) > 0 {
		return fmt.Errorf("missing names: %v", missingNames)
	}

	if err := r.storage.AddSegmentsToUser(toCreateIds, toDeleteIds, userID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUsersSegments(userID uint64) ([]dtos.SegmentDtoResponse, error) {
	segments, err := r.storage.GetSegmentsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var segmentsDto []dtos.SegmentDtoResponse
	for _, segment := range segments {
		segmentsDto = append(segmentsDto, dtos.ToSegmentDto(&segment))
	}

	return segmentsDto, nil
}
