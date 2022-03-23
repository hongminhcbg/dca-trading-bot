package store

import (
	"context"
	"dca-bot/models"
	"gorm.io/gorm"
)

type OrderTrackingStore interface {
	GetOrderByIndexNum(ctx context.Context, indexNum int64) (*models.OrderTracking, error)
	Save(ctx context.Context, r *models.OrderTracking) error
}

func NewOrderTrackingStore(db *gorm.DB) OrderTrackingStore {
	return &orderTrackingStoreImpl{db: db}
}

type orderTrackingStoreImpl struct {
	db *gorm.DB
}

func (s *orderTrackingStoreImpl) Save(ctx context.Context, r *models.OrderTracking) error {
	return s.db.WithContext(ctx).Save(r).Error
}

func (s *orderTrackingStoreImpl) GetOrderByIndexNum(ctx context.Context, indexNum int64) (*models.OrderTracking, error) {
	var result models.OrderTracking
	err := s.db.WithContext(ctx).Where("index_num = ?", indexNum).First(&result).Error
	if err != gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}
