package general

import "market/pkg/models"

type IStorage interface {
	AddLogs(logs []models.Log) error
	CountUsersNumber() (uint64, error)
}
