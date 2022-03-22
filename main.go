package main

import (
	"dca-bot/conf"
	"dca-bot/service"
	"encoding/json"
	"fmt"
	binance "github.com/adshao/go-binance/v2"
	"log"
	"os"
	"syscall"
)

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
	s := service.NewDCAService(client, config)
	go s.StartConsumerCheckTp()
	go s.StartConsumerCheckBuy()
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
