package user

import (
	"encoding/csv"
	"fmt"
	"market/pkg/dtos"
	"market/pkg/errors"
	"market/pkg/models"
	"market/pkg/storage"
	"os"
	"strconv"
	"time"
)

type IRepository interface {
	EditUsersSegments(toCreate []dtos.CreateSegmentDto, toDelete []string, userID uint64) error
	GetUsersSegments(userID uint64) ([]dtos.SegmentDtoResponse, error)
	CreateUserLogs(date string, userID uint64) (string, error)
}

type Repository struct {
	storage storage.IStorage
}

func NewRepository(storage storage.IStorage) IRepository {
	return &Repository{storage}
}

func (r *Repository) EditUsersSegments(toCreateDto []dtos.CreateSegmentDto, toDelete []string, userID uint64) error {
	toCreate := make([]string, len(toCreateDto))

	for i, segment := range toCreateDto {
		toCreate[i] = segment.Name
	}

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

	createLogs, err := r.storage.AddSegmentsToUser(toCreateIds, toCreateDto, userID)
	if err != nil {
		return fmt.Errorf("%s: %v", errors.CreateSegmentsErr, err)
	}

	deleteLogs, err := r.storage.DeleteSegmentsFromUser(toDeleteIds, userID)
	if err != nil {
		return fmt.Errorf("%s: %v", errors.DeleteSegmentsErr, err)
	}

	if len(createLogs) == 0 && len(deleteLogs) == 0 {
		return nil
	}

	logs := make([]models.Log, 0, len(createLogs)+len(deleteLogs))

	if len(createLogs) > 0 {
		for i, _ := range createLogs {
			createLogs[i].Segment = toCreate[i]
		}
		logs = append(logs, createLogs...)
	}

	if len(deleteLogs) > 0 {
		for i, _ := range deleteLogs {
			deleteLogs[i].Segment = toDelete[i]
		}
		logs = append(logs, deleteLogs...)
	}

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

func (r *Repository) CreateUserLogs(date string, userID uint64) (string, error) {

	t, err := time.Parse("2006-01", date)
	if err != nil {
		return "", fmt.Errorf(errors.InvalidDateErr)
	}

	startTime := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	endTime := startTime.AddDate(0, 1, 0).Add(-time.Nanosecond)

	logs, err := r.storage.GetUserLogs(&startTime, &endTime, userID)
	if err != nil {
		return "", err
	}

	file, err := os.Create("/tmp/" + strconv.FormatUint(userID, 10) + ".csv")
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, log := range logs {
		record := []string{
			strconv.FormatUint(log.UserID, 10),
			log.Segment,
			log.EventType,
			log.CreatedAt.String(),
		}
		err := writer.Write(record)
		if err != nil {
			return "", err
		}
	}

	return file.Name(), nil
}
