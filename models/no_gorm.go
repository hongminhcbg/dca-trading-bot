package models

type OrderExecuted struct {
	ExecutedQty string `json:"executed_qty"`
	UsdtQty     string `json:"usdt_qty"`
}
