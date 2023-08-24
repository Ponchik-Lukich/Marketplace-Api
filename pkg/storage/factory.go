package storage

import (
	"errors"
	"market/pkg/constant"
	"market/pkg/storage/postgresql"
)

func NewStorage(dbType string, cfg Config) (IStorage, error) {
	switch dbType {
	case constant.Postgres:
		if pCfg, ok := cfg.(*postgresql.Config); ok {
			return postgresql.NewStorage(*pCfg), nil
		} else {
			return nil, errors.New("invalid config for postgresql")
		}
	default:
		return nil, errors.New("unsupported database type")
	}
}
