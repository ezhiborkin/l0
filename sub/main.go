package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/nats-io/stan.go"
)

func main() {
	// Настройки подключения к серверу NATS Streaming
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
	handler := func(msg *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(msg.Data))
	}

	// Подписка на канал
	_, err = sc.Subscribe(channel, handler)
	if err != nil {
		log.Fatalf("Error subscribing to channel: %v", err)
	}

	fmt.Printf("Subscribed to channel %s. Waiting for messages...\n", channel)

	// Ожидание завершения приложения
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Application is shutting down.")
}
