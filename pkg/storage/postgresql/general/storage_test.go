package general

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"market/pkg/models"
	"market/pkg/storage/testsetup"
	"os"
	"testing"
	"time"
)

var genStorage *Storage

func TestMain(m *testing.M) {
	testsetup.Setup()
	genStorage = NewStorage(testsetup.DB)
	os.Exit(m.Run())
}

func TestInit(t *testing.T) {
	db := genStorage.Init()
	assert.IsType(t, &gorm.DB{}, db, "DB initialization failed")
}

func TestCountUsersNumber(t *testing.T) {
	users := []models.User{{ID: 1}, {ID: 2}}
	testsetup.DB.Create(&users)

	count, err := genStorage.CountUsersNumber()

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, uint64(2), count, "Count mismatch")
}

func TestAddLogs(t *testing.T) {
	logTime := time.Now().Add(time.Hour * 3)
	logs := []models.Log{
		{UserID: 1, Segment: "Segment1", EventType: "добавление", Time: &logTime},
		{UserID: 2, Segment: "Segment2", EventType: "удаление", Time: &logTime},
	}

	err := genStorage.AddLogs(logs)

	assert.Nil(t, err, "Error should be nil")

	var count int64
	testsetup.DB.Table("logs").Count(&count)
	assert.Equal(t, int64(2), count, "Count of logs should match number of inserted logs")
}
