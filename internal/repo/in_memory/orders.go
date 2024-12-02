package in_memory

import (
	"applicationDesignTest/internal/model"
	"applicationDesignTest/internal/repo"
	"applicationDesignTest/pkg/collections"
	"applicationDesignTest/pkg/in_memory_storage"
	"context"
	"errors"
)

type InMemoryOrdersRepo struct {
	storage in_memory_storage.InMemoryStorage
}

func NewInMemoryOrdersRepo(storage in_memory_storage.InMemoryStorage) repo.OrdersRepo {
	return &InMemoryOrdersRepo{
		storage: storage,
	}
}

func (o *InMemoryOrdersRepo) CreatOrder(ctx context.Context, order model.Order) error {
	err := o.storage.CreatOrder(ctx, order)
	if err != nil {
		if errors.Is(err, collections.ErrEntityAlreadyExists) {
			return repo.ErrOrderAlreadyCreated
		}
		return err
	}
	return nil
}

func (o *InMemoryOrdersRepo) OrderExists(ctx context.Context, order model.Order) bool {
	return o.storage.OrderExists(ctx, order)
}

func (o *InMemoryOrdersRepo) CancelOrder(ctx context.Context, order model.Order) error {
	return o.storage.CancelOrder(ctx, order)
}
