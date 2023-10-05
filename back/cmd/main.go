package main

import (
	"back"
	"back/pkg/cashe"
	"back/pkg/handler"
	order "back/pkg/order"
	"back/pkg/repository"
	"back/pkg/service"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"

	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error occured while initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error occured while loading env variables: %s", err.Error())
	}

	// Ожидание завершения приложения

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("error occured while connecting to db: %s", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("error occured while pinging db: %s", err.Error())
	}

	log.Printf("Connected to DB!")

	repos := repository.NewRepository(db)
	cache := cashe.NewCache(2*time.Hour, 2*time.Hour)
	services := service.NewService(repos, cache)
	handlers := handler.NewHandler(services)

	natsURL := "nats://nats-streaming:4222"
	clusterID := "mycluster"
	clientID := "subscriber-client"

	// Создание соединения
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming: %v", err)
	}
	defer sc.Close()

	// Название канала, на который вы хотите подписаться
	channel := "my-channel"

	// Обработчик сообщений
	handlerStan := func(msg *stan.Msg) {
		var orderData order.OrderData
		if err := json.Unmarshal(msg.Data, &orderData); err != nil {
			log.Fatalf("Error unmarshalling data: %v", err)
		}

		cache.Set(orderData.OrderUID, orderData, 1*time.Hour)
		orderDataFromCache, ok := cache.Get(orderData.OrderUID)
		if ok {
			fmt.Printf("Cache: %v", orderDataFromCache)
			log.Printf("Cache: %v", orderDataFromCache)
		} else {
			fmt.Println("Cache: not found")
			log.Println("Cache: not found")
		}

		if err := order.InsertOrder(db, orderData); err != nil {
			log.Fatalf("Error inserting order: %v", err)
		}

		if err := order.InsertDelivery(db, orderData); err != nil {
			log.Fatalf("Error inserting delivery: %v", err)
		}

		if err := order.InsertPayment(db, orderData); err != nil {
			log.Fatalf("Error inserting payment: %v", err)
		}

		if err := order.InsertItems(db, orderData); err != nil {
			log.Fatalf("Error inserting items: %v", err)
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	// Подписка на канал
	_, err = sc.Subscribe(channel, handlerStan)
	if err != nil {
		log.Fatalf("Error subscribing to channel: %v", err)
	}

	fmt.Printf("Subscribed to channel %s. Waiting for messages...\n", channel)

	srv := new(back.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
