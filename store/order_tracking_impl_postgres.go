package store

import (
	"context"
	"dca-bot/models"
	"gorm.io/gorm"
)

type orderTrackingPostgresImpl struct {
	*gorm.DB
}

func NewOrderTrackingPostgresImpl(db *gorm.DB) OrderTrackingStore {
	return &orderTrackingPostgresImpl{db}
}

func (x *orderTrackingPostgresImpl) Close() error {
	return x.Close()
}

func (x *orderTrackingPostgresImpl) GetOrderByIndexNum(ctx context.Context, indexNum int64) (*models.OrderTracking, error) {
	result := new(models.OrderTracking)

	err := x.WithContext(ctx).Where("index_num=?", indexNum).Order("id DESC").First(result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return result, nil
}

func (x *orderTrackingPostgresImpl) Save(ctx context.Context, r *models.OrderTracking) error {
	return x.WithContext(ctx).Save(r).Error
}
