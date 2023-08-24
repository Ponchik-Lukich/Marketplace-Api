package segment

import "market/pkg/storage"

type IRepository interface {
}

type Repository struct {
	storage storage.IStorage
}

func NewRepository(storage storage.IStorage) IRepository {
	return &Repository{storage}
}
