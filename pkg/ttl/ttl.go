package ttl

import (
	"log"
	"market/pkg/repository/segment"
	"time"
)

func DeleteExpiredSegments(repo segment.IRepository, moment time.Time) {
	err := repo.DeleteExpiredSegments(&moment)
	if err != nil {
		log.Fatalf("Error deleting expired segments: %v", err)
	}
}

func TtlWorker(repo segment.IRepository) {
	for {
		now := time.Now().Add(time.Hour * 3)
		next := now.Add(time.Minute).Truncate(time.Minute)
		time.Sleep(next.Sub(now))

		log.Println("Deleting expired segments")
		DeleteExpiredSegments(repo, now)
	}
}
