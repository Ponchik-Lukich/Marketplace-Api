package user

import (
	"fmt"
	"market/pkg/dtos"
	"market/pkg/errors"
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

	err = r.storage.CreateUser(userID)
	if err != nil {
		return fmt.Errorf(errors.UpdatingUserErr)
	}

	createLogs, err := r.storage.AddSegmentsToUser(toCreateIds, userID)
	if err != nil {
		return fmt.Errorf("%s: %v", errors.CreateSegmentsErr, err)
	}

	deleteLogs, err := r.storage.DeleteSegmentsFromUser(toDeleteIds, userID)
	if err != nil {
		return fmt.Errorf("%s: %v", errors.DeleteSegmentsErr, err)
	}

	logs := append(createLogs, deleteLogs...)

	err = r.storage.AddLogs(logs)
	if err != nil {
		return fmt.Errorf("%s: %v", errors.AddingLogsErr, err)
	}

	return nil
}

func (r *Repository) GetUsersSegments(userID uint64) ([]dtos.SegmentDtoResponse, error) {
	segments, err := r.storage.GetSegmentsByUserID(userID)
	if err != nil {
		return nil, err
	}

	err = r.storage.CreateUser(userID)
	if err != nil {
		return nil, fmt.Errorf(errors.UpdatingUserErr)
	}

	var segmentsDto []dtos.SegmentDtoResponse
	for _, segment := range segments {
		segmentsDto = append(segmentsDto, dtos.ToSegmentDto(&segment))
	}

	return segmentsDto, nil
}
