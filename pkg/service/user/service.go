package user

import (
	"encoding/csv"
	"fmt"
	"market/pkg/constant"
	"market/pkg/dtos"
	"market/pkg/errors"
	"market/pkg/models"
	"market/pkg/storage/general"
	"market/pkg/storage/user"
	"os"
	"strconv"
	"time"
)

type IRepository interface {
	EditUsersSegments(toCreate []dtos.CreateSegmentDto, toDelete []string, userID uint64) errors.CustomError
	GetUsersSegments(userID uint64) ([]dtos.SegmentDtoResponse, errors.CustomError)
	CreateUserLogs(date string, userID uint64) (string, errors.CustomError)
}

type Repository struct {
	userStorage user.IStorage
	genStorage  general.IStorage
}

func NewRepository(userStorage user.IStorage, genStorage general.IStorage) IRepository {
	return &Repository{userStorage, genStorage}
}

func (r *Repository) EditUsersSegments(toCreateDto []dtos.CreateSegmentDto, toDelete []string, userID uint64) errors.CustomError {
	toCreate := make([]string, len(toCreateDto))
	for i, segment := range toCreateDto {
		toCreate[i] = segment.Name
	}

	toCreateIds, toCreateMissingNames, err := r.userStorage.GetIDsAndMissingNames(toCreate)
	if err != nil {
		return errors.GetMissingNames{Err: fmt.Sprintf("%v", toCreateMissingNames)}
	}
	toDeleteIds, toDeleteMissingNames, err := r.userStorage.GetIDsAndMissingNames(toDelete)
	if err != nil {
		return errors.GetMissingNames{Err: fmt.Sprintf("%v", toDeleteMissingNames)}
	}

	missingNames := append(toCreateMissingNames, toDeleteMissingNames...)
	if len(missingNames) > 0 {
		return errors.MissingNames{Err: fmt.Sprintf("Missing names: %v", missingNames)}
	}

	if err := r.userStorage.CreateUser(userID); err != nil {
		return errors.UpdateUser{Err: err.Error()}
	}

	createLogs, customErr := r.userStorage.AddSegmentsToUser(toCreateIds, toCreateDto, userID)
	if customErr != nil {
		return customErr
	}

	deleteLogs, customErr := r.userStorage.DeleteSegmentsFromUser(toDeleteIds, toDelete, userID)
	if customErr != nil {
		return customErr
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

	err = r.genStorage.AddLogs(logs)
	if err != nil {
		return errors.AddLogs{Err: err.Error()}
	}

	return nil
}

func (r *Repository) GetUsersSegments(userID uint64) ([]dtos.SegmentDtoResponse, errors.CustomError) {
	segments, err := r.userStorage.GetUserSegmentsByID(userID)
	if err != nil {
		return nil, errors.GetSegmentByName{}
	}

	err = r.userStorage.CreateUser(userID)
	if err != nil {
		return nil, errors.UpdateUser{Err: err.Error()}
	}

	var segmentsDto []dtos.SegmentDtoResponse
	for _, segment := range segments {
		segmentsDto = append(segmentsDto, dtos.ToSegmentDto(&segment))
	}

	return segmentsDto, nil
}

func (r *Repository) CreateUserLogs(date string, userID uint64) (string, errors.CustomError) {

	t, err := time.Parse(constant.Layout, date)
	if err != nil {
		return "", errors.DateParsing{Err: err.Error()}
	}

	startTime := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	endTime := startTime.AddDate(0, 1, 0).Add(-time.Nanosecond)

	logs, err := r.userStorage.GetUserLogs(&startTime, &endTime, userID)
	if err != nil {
		return "", errors.GetLogs{Err: err.Error()}
	}

	file, err := os.Create(constant.DockerPath + strconv.FormatUint(userID, 10) + "_" + date + ".csv")
	if err != nil {
		return "", errors.CreateFile{Err: err.Error()}
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
			return "", errors.WriteFile{Err: err.Error()}
		}
	}

	return file.Name(), nil
}
