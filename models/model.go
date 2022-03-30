package models

import "time"

type OrderTracking struct {
	Id          int64      `json:"id" gorm:"column:id"`
	IndexNum    int64      `json:"index_num" gorm:"column:index_num"`
	SelectedNum int64      `json:"selected_num" gorm:"column:selected_num"`
	Status      string     `json:"status" gorm:"column:status"`
	Error       string     `json:"error" gorm:"column:error"`
	RawResponse string     `json:"raw_response" gorm:"column:raw_response"`
	ExecutedQty string     `json:"executed_qty" gorm:"column:executed_qty"`
	UsdtQty     string     `json:"usdt_qty" gorm:"column:usdt_qty"`
	CreatedAt   *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (OrderTracking) TableName() string {
	return "order_tracking"
}
