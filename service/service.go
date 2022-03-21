package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"sync"
)

type DCAService struct {
	mu sync.Mutex
	biCli *binance.Client
}
func NewDCAService(biCli *binance.Client) *DCAService {
	return &DCAService{
		mu:    sync.Mutex{},
		biCli: biCli,
	}
}
func (s *DCAService) Wallet() {
}

func (s *DCAService) MakeAnOrder() {
	resp , err := s.biCli.NewCreateOrderService().Symbol("BNBUSDT").
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeMarket).
		Quantity("0.1").
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	b, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Println("create order success", string(b))
	s.DoubleCheckOrder(resp.OrderID)
}

func (s *DCAService) DoubleCheckOrder(orderId int64) {
	order, err := s.biCli.NewGetOrderService().Symbol("BNBUSDT").
		OrderID(orderId).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	b, _ := json.MarshalIndent(order, "", "\t")
	fmt.Println("double check order success", string(b))
}

func (s *DCAService) GetAccountInfo() {
	res, err := s.biCli.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	b, _ := json.MarshalIndent(res, "", "\t")
	fmt.Println("get account success", string(b))
}