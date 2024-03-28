package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/basedalex/nats-service/internal/cache"
	"github.com/basedalex/nats-service/internal/config"
	"github.com/basedalex/nats-service/internal/db"
	"github.com/basedalex/nats-service/internal/stream"
	"github.com/basedalex/nats-service/model"
	"github.com/nats-io/stan.go"
)

func main() {
	// Load Config 
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	// Open DB
	dbConnect := cfg.DBConnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	pgConn, err := db.NewPostgres(ctx, dbConnect)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to db")

	// init cache
	cache := cache.New()
	orderCh := pgConn.GetAllOrders(ctx)
	for order := range orderCh {
		if order.Err != nil {
			log.Fatal(order.Err)
		}
		cache.Set(order.Data.ID, order.Data.OrderData)
	}

	// connect to nats streaming
	stanConn, err := stream.NewNats(cfg.NatsUrl, "server")
	if err != nil {
		log.Fatal(err)
	}
	// subscribe to channel in nats streaming
	_, err = stanConn.Subscribe("orders", func(msg *stan.Msg) {
		data := model.OrderData{}
		err = json.Unmarshal(msg.Data, &data)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("found message %s", string(msg.Data))
		err = pgConn.CreateOrder(ctx, data)
		if err != nil {
			log.Println(err)
			return
		}
	}, stan.StartWithLastReceived())

	if err != nil {
		log.Fatal(err)
	}
	// start http server
	serve(cache, cfg)
}

