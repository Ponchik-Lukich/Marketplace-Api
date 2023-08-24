package postgresql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"market/pkg/errors"
	"market/pkg/models"
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
	if err := s.db.AutoMigrate(&models.User{}, &models.Segment{}, &models.UserSegment{}); err != nil {
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

func (s *Storage) AddSegmentsToUser(toCreate []string, toDelete []string, userID uint64) []error {
	// add segments to users_segments
	//var segments []models.Segment
	//tx := s.db.Where("name IN ?", toCreate).Find(&segments)
	//if tx.Error != nil {
	//	return []error{tx.Error}
	//}
	//if tx.RowsAffected != int64(len(toCreate)) {
	//	return []error{fmt.Errorf(errors.SegmentsNotFoundErr)}
	//}
	//
	//var user models.User
	//tx = s.db.First(&user, userID)
	//if tx.Error != nil {
	//	return []error{tx.Error}
	//}
	//if tx.RowsAffected == 0 {
	//	return []error{fmt.Errorf(errors.UserNotFoundErr)}
	//}
	return nil
}
