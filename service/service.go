package service

import (
	"context"
	"crypto/rand"
	"dca-bot/models"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"sync"
	"time"

	"dca-bot/conf"
	"dca-bot/noti"
	"dca-bot/store"

	"github.com/adshao/go-binance/v2"
)

const (
	MAIN_SYMBOL = "BNBUSDT"
	MAIN_ASSET  = "BNB"
)

var Fibonacci = []int64{1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181, 80000}

type DCAService struct {
	mu            sync.Mutex
	biCli         *binance.Client
	currentAmount float64

	priceWillTookProfit float64
	block               int64
	intervalCheckTp     int64
	intervalCheckBuy    int64
	amountUsdtEachBlock string
	// tookProfit if this field is true my dream will become true
	tookProfit bool

	noti       noti.TelegramNoti
	orderStore store.OrderTrackingStore
}

func NewDCAService(biCli *binance.Client, notiSer noti.TelegramNoti, orderStore store.OrderTrackingStore, config *conf.Config) *DCAService {
	return &DCAService{
		mu:                  sync.Mutex{},
		biCli:               biCli,
		block:               config.Bock,
		intervalCheckTp:     config.IntervalCheckTp,
		intervalCheckBuy:    config.IntervalCheckBuy,
		priceWillTookProfit: config.PriceWillTookProfit,
		amountUsdtEachBlock: config.AmountUsdtEachBlock,
		orderStore:          orderStore,
		noti:                notiSer,
	}
}

func (s *DCAService) MakeAnOrder(quantity string) {
	resp, err := s.biCli.NewCreateOrderService().Symbol(MAIN_SYMBOL).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeMarket).
		QuoteOrderQty("11.0").
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	b, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Println("create order success", string(b))
	s.DoubleCheckOrder(resp.OrderID)
}

func (s *DCAService) DoubleCheckOrder(orderId int64) (string, error) {
	order, err := s.biCli.NewGetOrderService().Symbol(MAIN_SYMBOL).
		OrderID(orderId).Do(context.Background())
	if err != nil {
		fmt.Println("DoubleCheckOrder error", err)
		return "", err
	}

	if order.Status == binance.OrderStatusTypeFilled {
		b, _ := json.MarshalIndent(order, "", "\t")
		return string(b), nil
	}

	return "", fmt.Errorf("order status is %s", order.Status)
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

func (s *DCAService) orderExec(sideType binance.SideType, QuoteOrderQty string, Quantity string) (string, error) {
	fibonacciLevel := 1
	for {
		if fibonacciLevel > 12 {
			_ = s.noti.Send("retry many time but not success, please check manual")
			return "", fmt.Errorf("retry many time but not success")
		}

		cli := s.biCli.NewCreateOrderService().Symbol(MAIN_SYMBOL).
			Side(sideType).
			Type(binance.OrderTypeMarket)
		if len(QuoteOrderQty) > 0 {
			cli = cli.QuoteOrderQty(QuoteOrderQty)
		} else {
			cli = cli.Quantity(Quantity)
		}

		resp, err := cli.Do(context.Background())
		if err != nil {
			time.Sleep(time.Duration(Fibonacci[fibonacciLevel]) * time.Second)
			fibonacciLevel += 1
			_ = s.noti.Send("orderExec error: " + err.Error())
			log.Println(err, "orderExec internal server error")
			continue
		}

		if resp.Status == binance.OrderStatusTypeFilled {
			b, _ := json.MarshalIndent(resp, "", "\t")
			return string(b), nil
		}

		raw, err := s.DoubleCheckOrder(resp.OrderID)
		if err == nil {
			return raw, nil
		}

		time.Sleep(time.Duration(Fibonacci[fibonacciLevel]) * time.Second)
		fibonacciLevel += 1
		s.noti.Send("orderExec error: " + err.Error())
		log.Println(err, "orderExec got an error")
	}

}

func (s *DCAService) TPAllAhihi() {
	fibonacciLevel := 1
	for {
		time.Sleep(time.Duration(Fibonacci[fibonacciLevel]) * time.Second)
		fibonacciLevel += 1

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
			_, err = s.orderExec(binance.SideTypeSell, "", fmt.Sprintf("%0.2f", qty))
			if err == nil {
				log.Printf("sell success %f BNB\n", qty)
				//_ = s.noti.Send("TP success, congrats")
				return
			}
		}
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

// random return [from, max)
func random(from, max int64) int64 {
	maxBig := big.NewInt(max - from)
	for {
		r, err := rand.Int(rand.Reader, maxBig)
		if err == nil {
			return from + r.Int64()
		}
	}
}

func (s *DCAService) handleBuyOrder(r *models.OrderTracking, currentNumInOneBlock int64) {
	if r.Status == "SUCCESS" {
		log.Println("[DEBUG] order is success, do nothing")
		return
	}

	if currentNumInOneBlock != r.SelectedNum && r.SelectedNum != 0 {
		log.Println("[DEBUG] not buy now, do nothing", currentNumInOneBlock, r.SelectedNum)
		return
	}

	defer func() {
		_ = s.orderStore.Save(context.Background(), r)
	}()

	raw, err := s.orderExec(binance.SideTypeBuy, s.amountUsdtEachBlock, "")
	if err != nil {
		r.Error = fmt.Sprintf("%s\nmake_order_error:%s", r.Error, err.Error())
		r.Status = "ERROR"
		return
	}
	r.Status = "SUCCESS"
	s.noti.Send(raw)
	return
}

func (s *DCAService) checkBuy() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.tookProfit {
		log.Println("do nothing because took profit, please restart service manual")
		return true
	}

	ts := time.Now().Unix()
	indexNum := ts / s.block
	mod := ts % s.block
	currentIndexInBlock := mod / s.intervalCheckBuy
	maxIndexInOneBlock := s.block / s.intervalCheckBuy

	r, err := s.orderStore.GetOrderByIndexNum(context.Background(), indexNum)
	if err != nil {
		s.noti.Send("[ERROR] internal server" + err.Error())
		log.Println(err, "internal server error")
		return false
	}

	if r == nil {
		selectedNum := random(0, maxIndexInOneBlock)
		r = &models.OrderTracking{
			IndexNum:    indexNum,
			SelectedNum: selectedNum,
			Status:      "NONE",
			Error:       "",
		}

		err = s.orderStore.Save(context.Background(), r)
		if err != nil {
			s.noti.Send("[ERROR] checkBuy internal server" + err.Error())
			log.Println(err, "checkBuy internal server error")
			return false
		}
	}

	s.handleBuyOrder(r, currentIndexInBlock)
	return false
}

func (s *DCAService) StartConsumerCheckBuy() {
	ticker := time.NewTicker(time.Duration(s.intervalCheckBuy) * time.Second)
	defer ticker.Stop()
	for _ = range ticker.C {
		if shouldStopConsumer := s.checkBuy(); shouldStopConsumer {
			return
		}
	}
}
