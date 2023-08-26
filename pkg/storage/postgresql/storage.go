package postgresql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"market/pkg/models"
	"market/pkg/storage/general"
	stGeneral "market/pkg/storage/postgresql/general"
	stSegment "market/pkg/storage/postgresql/segment"
	stUser "market/pkg/storage/postgresql/user"
	"market/pkg/storage/segment"
	"market/pkg/storage/user"
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

func (s *Storage) GetSegmentStorage() segment.IStorage {
	return stSegment.NewStorage(s.db)
}

func (s *Storage) GetUserStorage() user.IStorage {
	return stUser.NewStorage(s.db)
}

func (s *Storage) GetGeneralStorage() general.IStorage {
	return stGeneral.NewStorage(s.db)
}
