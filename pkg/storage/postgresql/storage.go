package postgresql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"market/pkg/constant"
	"market/pkg/dtos"
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

func (s *Storage) CreateSegment(name string) (uint64, error) {
	db := s.Init()
	var segment models.Segment

	if err := db.Unscoped().Where("name = ?", name).First(&segment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			segment = models.Segment{Name: name}
			if err := db.Create(&segment).Error; err != nil {
				return 0, err
			}
			return segment.ID, nil
		} else {
			return 0, err
		}
	}

	if segment.DeletedAt.Valid {
		println("segment is deleted")
		if err := db.Unscoped().Model(&segment).Update("deleted_at", nil).Error; err != nil {
			return 0, err
		}
	}

	return segment.ID, nil
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
			return models.Segment{}, fmt.Errorf(errors.SegmentNotFoundErr400)
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
		return nil, errors.UserAlreadyHasSegment{Err: fmt.Sprintf("Existing Segments: %v", existingNames)}
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
		return nil, errors.UserDoesNotHaveSegment{Err: fmt.Sprintf("Non existing Segments: %v", nonExistingNames)}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.Transaction{Err: err.Error()}
	}

	return logs, nil
}

func (s *Storage) AddLogs(logs []models.Log) error {
	db := s.Init()
	return db.Create(&logs).Error
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

func (s *Storage) DeleteExpiredSegments(moment *time.Time) ([]models.Log, error) {
	db := s.Init()

	var segmentsToDelete []struct {
		UserID      uint64
		SegmentID   uint64
		SegmentName string
	}
	if err := db.Table("user_segments").
		Select("user_segments.user_id, user_segments.segment_id, segments.name as segment_name").
		Joins("JOIN segments ON user_segments.segment_id = segments.id").
		Where("user_segments.expiration_date < ?", moment).
		Where("user_segments.deleted_at IS NULL").
		Scan(&segmentsToDelete).Error; err != nil {
		return nil, err
	}

	log.Println("Delete number of segments:", len(segmentsToDelete))

	var logs []models.Log
	timeStamp := time.Now().Add(time.Hour * 3)

	for _, segment := range segmentsToDelete {
		log := models.Log{
			UserID:    segment.UserID,
			EventType: "удаление",
			Segment:   segment.SegmentName,
			Time:      &timeStamp,
		}
		logs = append(logs, log)
	}

	if err := db.Where("expiration_date < ?", moment).Delete(&models.UserSegment{}).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *Storage) CountUsersNumber() (uint64, error) {
	db := s.Init()

	var count int64
	if err := db.Table("users").Count(&count).Error; err != nil {
		return 0, err
	}

	return uint64(count), nil
}

func (s *Storage) AddSegmentsToUsersByPercent(totalUsers uint64, segmentId uint64, percent int) error {
	numUsersToAssign := int(float64(totalUsers) * float64(percent) / 100)

	var users []models.User
	db := s.Init()
	db.Model(&models.User{}).Order(gorm.Expr("RAND()")).Limit(numUsersToAssign).Find(&users)

	for _, user := range users {
		userSegment := models.UserSegment{
			UserID:    user.ID,
			SegmentID: segmentId,
		}
		db.Create(&userSegment)
	}
	return nil
}
