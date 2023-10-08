package main

import (
	"fmt"
	"pub"
	"pub/pkg/handler"

	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
)

type Order struct {
	OrderUID    string `json:"order_uid"`
	TrackNumber string `json:"track_number"`
}

func main() {

	natsURL := "nats://nats-streaming:4222"
	clusterID := "mycluster"
	clientID := "publisher-client"

	if err := initConfig(); err != nil {
		fmt.Printf("error initializing configs: %s", err.Error())
	}

	// Создание соединения
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		fmt.Printf("Error connecting to NATS Streaming: %v", err)
	}
	defer sc.Close()

	handlers := handler.NewHandler(sc)

	// Запуск HTTP-сервера
	srv := new(pub.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		fmt.Printf("error occured while running http server: %s", err.Error())
	}
	fmt.Printf("HTTP server listening on %s\n", viper.GetString("port"))

	// if err := sc.Publish(channel, []byte("kekus")); err != nil {
	// 	fmt.Printf("Error publishing message: %v", err)
	// 	return
	// }
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
