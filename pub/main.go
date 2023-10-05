package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nats-io/stan.go"
)

type Order struct {
	OrderUID    string `json:"order_uid"`
	TrackNumber string `json:"track_number"`
}

func main() {
	// Настройки подключения к серверу NATS Streaming
	natsURL := "nats://nats-streaming:4222"
	clusterID := "mycluster"
	clientID := "publisher-client"

	// Создание соединения
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming: %v", err)
	}
	defer sc.Close()

	// Channel, в которую будем публиковать
	channel := "my-channel"

	// Обработчик HTTP-запросов
	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		// var order []byte
		// if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		// 	http.Error(w, "Couldn't decode json", http.StatusExpectationFailed)
		// }

		order, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Couldn't read request body", http.StatusInternalServerError)
			return
		}

		// // Чтение данных из тела HTTP-запроса
		// messageBody, err := json.Marshal(order)
		// if err != nil {
		// 	http.Error(w, "Couldn't read request body", http.StatusInternalServerError)
		// }

		log.Println(order)

		// Опубликовать сообщение в тему
		if err := sc.Publish(channel, order); err != nil {
			http.Error(w, fmt.Sprintf("Error publishing message: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Message published successfully.")
	})

	// Запуск HTTP-сервера
	serverPort := ":8081"
	fmt.Printf("HTTP server listening on %s\n", serverPort)
	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
