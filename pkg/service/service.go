package service

import (
	"market/pkg/service/segment"
	"market/pkg/service/user"
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
	return user.NewRepository(r.storage.GetUserStorage(), r.storage.GetGeneralStorage())
}

func (r *Repositories) GetSegmentsRepo() segment.IRepository {
	return segment.NewRepository(r.storage.GetSegmentStorage(), r.storage.GetGeneralStorage())
}

func NewRepositories(storage storage.IStorage) IRepositories {
	return &Repositories{storage}
}
