package ttl

import (
	"log"
	"market/pkg/repository/segment"
	"time"
)

func DeleteExpiredSegments(repo segment.IRepository, moment time.Time) {
	err := repo.DeleteExpiredSegments(&moment)
	if err != nil {
		log.Printf("Error deleting expired segments: %v\n", err)
	}
}

func WorkerTtl(repo segment.IRepository) {
	for {
		now := time.Now().Add(time.Hour * 3)
		next := now.Add(time.Minute / 2).Truncate(time.Minute / 2)
		time.Sleep(next.Sub(now))

		DeleteExpiredSegments(repo, now)
	}
}
