package store

import (
	"context"
	"dca-bot/models"
)

type OrderTrackingStore interface {
	GetOrderByIndexNum(ctx context.Context, indexNum int64) (*models.OrderTracking, error)
	Save(ctx context.Context, r *models.OrderTracking) error
}