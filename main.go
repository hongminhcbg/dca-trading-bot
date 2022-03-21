package main

import (
	"context"
	"dca-bot/conf"
	"fmt"
	"time"

	binance "github.com/adshao/go-binance/v2"
)
func main(){
	config, err := conf.LoadConfig()
	if err != nil {
		panic(err)
	}

	binance.UseTestnet = true
	client := binance.NewClient(config.ApiKey, config.SecretKey)
	for {
		prices, err := client.NewListPricesService().Symbol("BNBUSDT").Do(context.Background())
		if err != nil {
			panic(err)
		}

		for _, p := range prices {
			fmt.Println(p)
		}

		time.Sleep(10*time.Second)
	}

}
