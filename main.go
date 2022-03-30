package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"syscall"

	binance "github.com/adshao/go-binance/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"dca-bot/conf"
	"dca-bot/noti"
	"dca-bot/service"
	"dca-bot/store"
)

func init() {
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
	db := mustConnectPsqlv2(config)

	binance.UseTestnet = config.IsTestnet
	client := binance.NewClient(config.ApiKey, config.SecretKey)
	osSignal := make(chan os.Signal, 1)
	notiService, err := noti.NewTelegramNoti(config)
	if err != nil {
		panic(err)
	}

	err = notiService.Send("server will start now")
	if err != nil {
		if !config.IsTestnet {
			panic(err)
		}

		log.Println("[ERROR] send noti error")
	}

	//orderTrackingStore := store.NewOrderCheckingStoreMapImpl()
	orderTrackingStore := store.NewOrderTrackingPostgresImpl(db)
	s := service.NewDCAService(client, notiService, orderTrackingStore, config)
	go s.StartConsumerCheckTp()
	go s.StartConsumerCheckBuy()
	go s.StartWebServer()
	for {
		select {
		case sig := <-osSignal:
			if sig == syscall.SIGKILL || sig == syscall.SIGINT {
				log.Println("Received signal kill, stop")
				_ = orderTrackingStore.Close()
				os.Exit(0)
			}
		}
	}
}

func mustConnectPostgres(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = sqlDb.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func mustConnectPsqlv2(cfg *conf.Config) *gorm.DB {
	dbPsql, err := sql.Open("postgres", cfg.DatabaseUrl)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbPsql}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = sqlDb.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
