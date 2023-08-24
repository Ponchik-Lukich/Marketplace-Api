package postgresql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"market/pkg/errors"
	"market/pkg/models"
	"time"
)

type Storage struct {
	cfg Config
	db  *gorm.DB
}

func NewStorage(cfg Config) *Storage {
	return &Storage{cfg: cfg.withDefaults()}
}

func (s *Storage) Connect() error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.cfg.Host, s.cfg.Port, s.cfg.User, s.cfg.Password, s.cfg.Database)

	database, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}

	s.db = database
	return nil
}

func (s *Storage) Close() error {
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func (s *Storage) Init() *gorm.DB {
	return s.db
}

func (s *Storage) MakeMigrations() error {
	if err := s.db.AutoMigrate(&models.User{}, &models.Segment{}, &models.UserSegment{}, &models.Log{}); err != nil {
		return fmt.Errorf("failed migrations: %w", err)
	}
	return nil
}

func (s *Storage) CreateSegment(name string) error {
	db := s.Init()
	segment := models.Segment{Name: name}
	return db.Create(&segment).Error
}

func (s *Storage) DeleteSegment(name string) error {
	db := s.Init()
	segment := models.Segment{Name: name}
	return db.Where("name = ?", name).Delete(&segment).Error
}

func (s *Storage) GetSegmentByName(name string) (models.Segment, error) {
	var segment models.Segment

	db := s.Init()
	tx := db.Where("name = ?", name).First(&segment)
	if tx.Error != nil {
		if tx.RowsAffected == 0 {
			return models.Segment{}, fmt.Errorf(errors.SegmentNotFoundErr)
		}
		return models.Segment{}, tx.Error
	}
	return segment, nil
}

func (s *Storage) GetSegmentsByUserID(userID uint64) ([]models.Segment, error) {
	var segments []models.Segment

	db := s.Init()
	err := db.Model(&models.User{ID: userID}).Association("Segments").Find(&segments)
	return segments, err
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

func (s *Storage) AddSegmentsToUser(toCreate []uint64, userID uint64) ([]models.Log, error) {
	db := s.Init()

	var logs []models.Log

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	for _, segmentID := range toCreate {
		var existingUserSegment models.UserSegment

		if err := tx.Unscoped().Where("user_id = ? AND segment_id = ?", userID, segmentID).First(&existingUserSegment).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				tx.Rollback()
				return nil, err
			}
		} else {
			if existingUserSegment.DeletedAt.Valid {
				if err := tx.Unscoped().Model(&existingUserSegment).Update("deleted_at", nil).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
				continue
			} else {
				return nil, fmt.Errorf(errors.UserAlreadyHasSegmentErr + " " + fmt.Sprintf("%d", segmentID))
			}
		}

		userSegment := models.UserSegment{
			UserID:    userID,
			SegmentID: segmentID,
		}

		if err := tx.Create(&userSegment).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		timeStamp := time.Now()

		log := models.Log{
			UserID:    userID,
			EventType: "добавление",
			Segment:   "",
			Time:      &timeStamp,
		}

		logs = append(logs, log)
	}

	return logs, tx.Commit().Error
}

func (s *Storage) DeleteSegmentsFromUser(toDelete []uint64, userID uint64) ([]models.Log, error) {
	db := s.Init()

	var logs []models.Log

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	for _, segmentID := range toDelete {
		var existingUserSegment models.UserSegment

		if err := tx.Unscoped().Where("user_id = ? AND segment_id = ?", userID, segmentID).First(&existingUserSegment).Error; err != nil {
			return nil, err
		} else {
			if existingUserSegment.DeletedAt.Valid {
				return nil, fmt.Errorf(errors.UserDoesNotHaveSegmentErr + " " + fmt.Sprintf("%d", segmentID))
			} else {
				if err := tx.Delete(&existingUserSegment).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}

		timeStamp := time.Now()

		log := models.Log{
			UserID:    userID,
			EventType: "добавление",
			Segment:   "",
			Time:      &timeStamp,
		}

		logs = append(logs, log)
	}

	return logs, tx.Commit().Error
}

func (s *Storage) AddLogs(logs []models.Log) error {
	db := s.Init()
	return db.Create(&logs).Error
}

func (s *Storage) CreateUser(userID uint64) error {
	db := s.Init()
	user := models.User{ID: userID}
	return db.Create(&user).Error
}
