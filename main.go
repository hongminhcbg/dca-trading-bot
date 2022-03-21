package main

import (
	"context"
	"dca-bot/conf"
	"dca-bot/service"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	binance "github.com/adshao/go-binance/v2"
)

func ListPrice(client *binance.Client, ticker <-chan time.Time) {
	for _ = range ticker {
		fmt.Println("receiver ticker")
		prices, err := client.NewListPricesService().Symbol("BNBUSDT").Do(context.Background())
		if err != nil {
			panic(err)
		}

		for _, p := range prices {
			fmt.Println(p)
		}
	}
}

func main() {
	config, err := conf.LoadConfig()
	if err != nil {
		panic(err)
	}

	b, _ := json.MarshalIndent(config, "", "\t")
	fmt.Println("start server with config", string(b))

	binance.UseTestnet = config.IsTestnet
	client := binance.NewClient(config.ApiKey, config.SecretKey)
	osSignal := make(chan os.Signal, 1)
	ticker30s := time.Tick(30 * time.Second)
	go ListPrice(client, ticker30s)
	s := service.NewDCAService(client)
	s.MakeAnOrder()
	time.Sleep(time.Second)
	s.GetAccountInfo()
	for {
		select {
		case sig := <-osSignal:
			if sig == syscall.SIGKILL || sig == syscall.SIGINT {
				log.Println("Receive signal kill, stop")
				os.Exit(0)
			}
		}
	}
}
