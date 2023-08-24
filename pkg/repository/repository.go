package repository

import (
	"market/pkg/repository/segment"
	"market/pkg/repository/user"
	"market/pkg/storage"
)

type IRepositories interface {
	GetUsersRepo() user.IRepository
	GetSegmentsRepo() segment.IRepository
}

type Repositories struct {
	storage storage.IStorage
}

func (r *Repositories) GetUsersRepo() user.IRepository {
	return user.NewRepository(r.storage)
}

func (r *Repositories) GetSegmentsRepo() segment.IRepository {
	return segment.NewRepository(r.storage)
}

func NewRepositories(storage storage.IStorage) IRepositories {
	return &Repositories{storage}
}
