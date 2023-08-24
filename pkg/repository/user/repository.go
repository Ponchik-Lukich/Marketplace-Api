package user

import (
	"fmt"
	"market/pkg/storage"
)

type IRepository interface {
	EditUsersSlugs(toCreate []string, toDelete []string, userID uint64) error
}

type Repository struct {
	storage storage.IStorage
}

func NewRepository(storage storage.IStorage) IRepository {
	return &Repository{storage}
}

func (r *Repository) EditUsersSlugs(toCreate []string, toDelete []string, userID uint64) error {
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
