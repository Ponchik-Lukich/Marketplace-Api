package testsetup

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

const (
	configPath = "config"
)

var DB *gorm.DB

func Setup() {
	//ctx := context.Background()
	//var cfg Config
	//err := confita.NewLoader(
	//	file.NewBackend(fmt.Sprintf("%s/default.yaml", configPath)),
	//).Load(ctx, &cfg)
	//if err != nil {
	//	log.Fatalf("parse config err: %v", err)
	//}

	var err error
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "user", "password", "avito")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open db connection: %v", err)
	}
	Cleanup()
}

func Cleanup() {
	DB.Exec("TRUNCATE users CASCADE")
	DB.Exec("TRUNCATE segments CASCADE")
	DB.Exec("TRUNCATE user_segments CASCADE")
	DB.Exec("TRUNCATE logs CASCADE")
	DB.Exec("ALTER SEQUENCE segments_id_seq RESTART WITH 1")
	DB.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
	DB.Exec("ALTER SEQUENCE user_segments_id_seq RESTART WITH 1")
	DB.Exec("ALTER SEQUENCE logs_id_seq RESTART WITH 1")
}
