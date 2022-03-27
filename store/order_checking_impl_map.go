package store

import (
	"context"
	"dca-bot/models"
	"sync"
)

type orderCheckingStoreMapImpl struct {
	mu   sync.Mutex
	data map[int64]*models.OrderTracking
}

func NewOrderCheckingStoreMapImpl() OrderTrackingStore {
	m := make(map[int64]*models.OrderTracking)
	return &orderCheckingStoreMapImpl{data: m}
}

func (x *orderCheckingStoreMapImpl) GetOrderByIndexNum(ctx context.Context, indexNum int64) (*models.OrderTracking, error) {
	order, ok := x.data[indexNum]
	if !ok || order == nil {
		return nil, nil
	}

	return order, nil
}

func (x *orderCheckingStoreMapImpl) Save(ctx context.Context, r *models.OrderTracking) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	x.data = map[int64]*models.OrderTracking{
		r.IndexNum: r,
	}

	return nil
}
