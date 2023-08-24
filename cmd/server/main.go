package main

import (
	"context"
	"fmt"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"log"
	"market/pkg/repository"
	"market/pkg/router"
	"market/pkg/storage"
	"net/http"
)

const (
	configPath = "config"
	envPath    = "config/.env"
)

func LoadConfiguration() (Config, error) {
	ctx := context.Background()
	var cfg Config
	err := confita.NewLoader(
		file.NewBackend(fmt.Sprintf("%s/default.yaml", configPath)),
	).Load(ctx, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

func InitializeStorage(cfg *Config) (storage.IStorage, error) {
	db, err := storage.NewStorage(cfg.DatabaseType, &cfg.DatabaseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	if err := db.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func main() {
	cfg, err := LoadConfiguration()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	db, err := InitializeStorage(&cfg)
	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}

	err = db.MakeMigrations()
	if err != nil {
		log.Fatalf("failed to make migrations: %v", err)
	}

	repos := repository.NewRepositories(db)
	r := router.SetupRouter(repos)

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Close(); err != nil {
		log.Fatalf("failed to close connection to database: %v", err)
	}
}
