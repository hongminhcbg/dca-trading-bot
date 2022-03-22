package service

import (
	"context"
	"dca-bot/conf"
	"dca-bot/noti"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"log"
	"strconv"
	"sync"
	"time"
)

const (
	MAIN_SYMBOL = "BNBUSDT"
	MAIN_ASSET  = "BNB"
)

type DCAService struct {
	mu            sync.Mutex
	biCli         *binance.Client
	currentAmount float64

	priceWillTookProfit float64
	block               int64
	intervalCheckTp     int64
	intervalCheckBuy    int64

	// tookProfit if this field is true my dream will become true
	tookProfit bool

	noti noti.TelegramNoti
}

func NewDCAService(biCli *binance.Client, config *conf.Config) *DCAService {
	return &DCAService{
		mu:    sync.Mutex{},
		biCli: biCli,
		block: config.Bock,
		intervalCheckTp: config.IntervalCheckTp,
		intervalCheckBuy: config.IntervalCheckBuy,
		priceWillTookProfit: config.PriceWillTookProfit,
	}
}

func (s *DCAService) Wallet() {
}

func (s *DCAService) MakeAnOrder() {
	resp, err := s.biCli.NewCreateOrderService().Symbol("BNBUSDT").
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

func (s *DCAService) DoubleCheckOrder(orderId int64) error {
	order, err := s.biCli.NewGetOrderService().Symbol(MAIN_SYMBOL).
		OrderID(orderId).Do(context.Background())
	if err != nil {
		fmt.Println("DoubleCheckOrder error", err)
		return err
	}

	if order.Status == binance.OrderStatusTypeFilled {
		return nil
	}

	return fmt.Errorf("order status is %s", order.Status)
}

func (s *DCAService) GetAccountInfo() (*binance.Account, error) {
	res, err := s.biCli.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	b, _ := json.MarshalIndent(res, "", "\t")
	fmt.Println("get account success", string(b))
	return res, nil
}

func (s *DCAService) orderExec(sideType binance.SideType, quantity string) error {
	resp, err := s.biCli.NewCreateOrderService().Symbol(MAIN_SYMBOL).
		Side(sideType).
		Type(binance.OrderTypeMarket).
		Quantity(quantity).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	if resp.Status == binance.OrderStatusTypeFilled {
		return nil
	}

	time.Sleep(time.Second)
	return s.DoubleCheckOrder(resp.OrderID)
}

func (s *DCAService) TPAllAhihi() {
	for {
		account, err := s.GetAccountInfo()
		if err != nil {
			_ = s.noti.Send("TP but get error when get amount, please to manual")
			continue
		}

		for _, asset := range account.Balances {
			if asset.Asset != MAIN_ASSET {
				continue
			}

			qty, _ := strconv.ParseFloat(asset.Free, 64)
			realQty := qty - 1000
			err = s.orderExec(binance.SideTypeSell, fmt.Sprintf("%f", realQty))
			if err == nil {
				log.Printf("sell success %f BNB\n", realQty)
				//_ = s.noti.Send("TP success, congrats")
				return
			}
		}
		time.Sleep(2 * time.Second)
	}
}
func (s *DCAService) shouldTookProfit() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	prices, err := s.biCli.NewListPricesService().Symbol(MAIN_SYMBOL).Do(context.Background())
	if err != nil {
		log.Println(err, "got price error")
		return false
	}

	for _, p := range prices {
		if p.Symbol != MAIN_SYMBOL {
			continue
		}

		nowPrice, err := strconv.ParseFloat(p.Price, 64)
		if err != nil || nowPrice < s.priceWillTookProfit {
			log.Printf("now price is %f, too low\n", nowPrice)
			continue
		}

		log.Println("Congrats, we will TP at price: ", nowPrice)
		go s.TPAllAhihi()
		s.tookProfit = true
		return true
	}

	return false
}

func (s *DCAService) StartConsumerCheckTp() {
	ticker := time.NewTicker(time.Duration(s.intervalCheckTp) * time.Second)
	defer ticker.Stop()
	for _ = range ticker.C {
		if tookProfit := s.shouldTookProfit(); tookProfit {
			return
		}
	}
}

func (s *DCAService) checkBuy() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.tookProfit {
		log.Println("do nothing because took profit, please restart service manual")
		return true
	}

	log.Println("Buy success")
	return false
}

func (s *DCAService) StartConsumerCheckBuy()  {
	ticker := time.NewTicker(time.Duration(s.intervalCheckBuy) * time.Second)
	defer ticker.Stop()
	for _ = range ticker.C {
		if shouldStopConsumer := s.checkBuy(); shouldStopConsumer {
			return
		}
	}
}
