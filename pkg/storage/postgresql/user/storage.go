package user

import (
	"fmt"
	"gorm.io/gorm"
	"market/pkg/constant"
	"market/pkg/dtos"
	"market/pkg/errors"
	"market/pkg/models"
	"time"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Init() *gorm.DB {
	return s.db
}

func (s *Storage) GetIDsAndMissingNames(names []string) ([]uint64, []string, error) {
	var existingIds []uint64
	var missingNames []string

	db := s.Init()
	var segments []models.Segment
	if err := db.Where("name IN ?", names).Find(&segments).Error; err != nil {
		return nil, nil, err
	}

	existingNamesMap := make(map[string]bool)
	for _, segment := range segments {
		existingIds = append(existingIds, segment.ID)
		existingNamesMap[segment.Name] = true
	}

	for _, name := range names {
		if _, ok := existingNamesMap[name]; !ok {
			missingNames = append(missingNames, name)
		}
	}

	return existingIds, missingNames, nil
}

func (s *Storage) AddSegmentsToUser(toCreate []uint64, toCreateDto []dtos.CreateSegmentDto, userID uint64) ([]models.Log, errors.CustomError) {
	var logs []models.Log
	var existingNames []string
	flag := false

	db := s.Init()
	tx := db.Begin()
	if tx.Error != nil {
		return nil, errors.Transaction{Err: tx.Error.Error()}
	}

	for i, segmentID := range toCreate {
		var existingUserSegment models.UserSegment
		stop := false

		if err := tx.Unscoped().Where("user_id = ? AND segment_id = ?", userID, segmentID).First(&existingUserSegment).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				tx.Rollback()
				return nil, errors.AddSegments{Err: err.Error()}
			}
		} else {
			if existingUserSegment.DeletedAt.Valid {
				if err := tx.Unscoped().Model(&existingUserSegment).Update("deleted_at", nil).Error; err != nil {
					tx.Rollback()
					return nil, errors.AddSegments{Err: err.Error()}
				}
				stop = true
			} else {
				flag = true
				existingNames = append(existingNames, toCreateDto[i].Name)
			}
		}

		if !stop {
			userSegment := models.UserSegment{
				UserID:    userID,
				SegmentID: segmentID,
			}

			if toCreateDto[i].DeleteTime != "" {
				expireTime, err := time.Parse(constant.FullLayout, toCreateDto[i].DeleteTime)
				if err != nil {
					tx.Rollback()
					return nil, errors.DateParsing{Err: err.Error()}
				}
				userSegment.ExpirationDate = expireTime
			}

			if err := tx.Create(&userSegment).Error; err != nil {
				tx.Rollback()
				return nil, errors.CreateSegment{Err: err.Error()}
			}
		}

		timeStamp := time.Now().Add(time.Hour * 3)

		addLog := models.Log{
			UserID:    userID,
			EventType: "добавление",
			Segment:   "",
			Time:      &timeStamp,
		}

		logs = append(logs, addLog)
	}
	if flag {
		tx.Rollback()
		return nil, errors.UserAlreadyHasSegment{Err: fmt.Sprintf("Ids: %v", existingNames)}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, errors.Transaction{Err: tx.Error.Error()}
	}

	return logs, nil
}

func (s *Storage) DeleteSegmentsFromUser(toDelete []uint64, toDeleteDto []string, userID uint64) ([]models.Log, errors.CustomError) {
	var logs []models.Log
	var nonExistingNames []string
	flag := false

	db := s.Init()
	tx := db.Begin()
	if tx.Error != nil {
		return nil, errors.Transaction{Err: tx.Error.Error()}
	}

	for i, segmentID := range toDelete {
		var existingUserSegment models.UserSegment

		if err := tx.Where("user_id = ? AND segment_id = ?", userID, segmentID).First(&existingUserSegment).Error; err != nil {
			flag = true
			nonExistingNames = append(nonExistingNames, toDeleteDto[i])
		} else {
			if err := tx.Delete(&existingUserSegment).Error; err != nil {
				tx.Rollback()
				return nil, errors.DeleteSegments{Err: err.Error()}
			}
		}

		timeStamp := time.Now().Add(time.Hour * 3)

		addLog := models.Log{
			UserID:    userID,
			EventType: "удаление",
			Segment:   "",
			Time:      &timeStamp,
		}

		logs = append(logs, addLog)
	}

	if flag {
		tx.Rollback()
		return nil, errors.UserDoesNotHaveSegment{Err: fmt.Sprintf("Ids: %v", nonExistingNames)}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.Transaction{Err: err.Error()}
	}

	return logs, nil
}

func (s *Storage) CreateUser(userID uint64) error {
	db := s.Init()
	user := models.User{}
	if err := db.First(&user, userID).Error; err == gorm.ErrRecordNotFound {
		newUser := models.User{ID: userID}
		return db.Create(&newUser).Error
	}
	return nil
}

func (s *Storage) GetUserLogs(start *time.Time, end *time.Time, userID uint64) ([]models.Log, error) {
	var logs []models.Log

	db := s.Init()
	if err := db.Where("user_id = ? AND time > ? AND time < ?", userID, start, end).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (s *Storage) GetUserSegmentsByID(userID uint64) ([]models.Segment, error) {
	var segments []models.Segment

	db := s.Init()
	err := db.Model(&models.User{ID: userID}).Association("Segments").Find(&segments)
	return segments, err
}
