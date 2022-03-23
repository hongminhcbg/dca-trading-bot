package main

import (
	"crypto/rand"
	"dca-bot/conf"
	"dca-bot/noti"
	"dca-bot/service"
	"dca-bot/store"
	"encoding/json"
	"fmt"

	binance "github.com/adshao/go-binance/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"log"
	"os"
	"syscall"
)

func init()  {
	b := make([]byte, 128)
	rand.Read(b)
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
	notiService, err := noti.NewTelegramNoti(config)
	if err != nil {
		panic(err)
	}

	err = notiService.Send("server will start now")
	if err != nil {
		panic(err)
	}

	db := newDB(config.MysqlDsn)
	orderTrackingStore := store.NewOrderTrackingStore(db)
	s := service.NewDCAService(client, notiService, orderTrackingStore, config)
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

func newDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// force a connection and test that it worked
	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
