package segment

import (
	"fmt"
	"gorm.io/gorm"
	"log"
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
		if err := db.Unscoped().Model(&segment).Update("deleted_at", nil).Error; err != nil {
			return 0, err
		}
	}

	return segment.ID, nil
}

func (s *Storage) DeleteSegment(name string) ([]models.Log, error) {
	db := s.Init()
	var segment models.Segment
	var userSegments []models.UserSegment
	var logs []models.Log

	if err := db.Where("name = ?", name).First(&segment).Error; err != nil {
		return nil, err
	}

	if err := db.Where("segment_id = ?", segment.ID).Find(&userSegments).Error; err != nil {
		return nil, err
	}

	for _, userSegment := range userSegments {
		timeStamp := time.Now().Add(time.Hour * 3)
		addLog := models.Log{
			UserID:    userSegment.UserID,
			EventType: "удаление",
			Segment:   segment.Name,
			Time:      &timeStamp,
		}
		logs = append(logs, addLog)
	}

	if err := db.Model(&segment).Association("Users").Clear(); err != nil {
		return nil, err
	}

	if err := db.Delete(&segment).Error; err != nil {
		return nil, err
	}

	return logs, nil
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
		Where("user_segments.expiration_date <= ?", moment).
		Where("user_segments.deleted_at IS NULL").
		Scan(&segmentsToDelete).Error; err != nil {
		return nil, err
	}

	log.Println("Delete number of segments:", len(segmentsToDelete))

	var logs []models.Log

	for _, segment := range segmentsToDelete {
		timeStamp := time.Now().Add(time.Hour * 3)
		addLog := models.Log{
			UserID:    segment.UserID,
			EventType: "удаление",
			Segment:   segment.SegmentName,
			Time:      &timeStamp,
		}
		logs = append(logs, addLog)
	}

	if err := db.Where("expiration_date <= ?", moment).Delete(&models.UserSegment{}).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *Storage) AddSegmentsToUsersByPercent(totalUsers uint64, segmentId uint64, percent int) ([]models.Log, error) {
	numUsersToAssign := int(float64(totalUsers) * float64(percent) / 100)
	var users []models.User
	var logs []models.Log
	var segment models.Segment

	db := s.Init()

	if err := db.Model(&models.Segment{}).Where("id = ?", segmentId).First(&segment).Error; err != nil {
		return nil, err
	}

	db.Model(&models.User{}).Order(gorm.Expr("RAND()")).Limit(numUsersToAssign).Find(&users)

	for _, user := range users {
		userSegment := models.UserSegment{
			UserID:    user.ID,
			SegmentID: segmentId,
		}
		err := db.Create(&userSegment)
		if err.Error != nil {
			return nil, err.Error
		}

		timeStamp := time.Now().Add(time.Hour * 3)
		addLog := models.Log{
			UserID:    user.ID,
			EventType: "добавление",
			Segment:   segment.Name,
			Time:      &timeStamp,
		}
		logs = append(logs, addLog)
	}
	return logs, nil
}
