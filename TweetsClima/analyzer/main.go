package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

type Tweet struct {
	Body string    `json:"body"`
	Time time.Time `json:"time"`
}

func connectRabbitMQ() *amqp.Connection {
	var conn *amqp.Connection
	var err error

	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial("amqp://admin:admin@rabbitmq:5672/")
		if err == nil {
			log.Println("✅ Conectado a RabbitMQ")
			return conn
		}
		log.Printf("⏳ Intentando conectar a RabbitMQ... (%d/10)", i+1)
		time.Sleep(3 * time.Second)
	}

	log.Fatalf("❌ No se pudo conectar a RabbitMQ después de varios intentos: %v", err)
	return nil
}

func appendTweetToFile(tweet Tweet) {
	log.Println("📁 Intentando guardar tweet en tweets.json...")

	file, err := os.OpenFile("tweets.json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("❌ Error abriendo archivo: %v", err)
		return
	}
	defer file.Close()

	data, err := json.Marshal(tweet)
	if err != nil {
		log.Printf("❌ Error convirtiendo a JSON: %v", err)
		return
	}

	_, err = file.Write(append(data, '\n'))
	if err != nil {
		log.Printf("❌ Error escribiendo al archivo: %v", err)
	} else {
		log.Println("✅ Tweet guardado en tweets.json")
	}
}

func main() {
	conn := connectRabbitMQ()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Error al abrir canal: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare("tweets", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ Error al declarar cola: %v", err)
	}

	msgs, err := ch.Consume("tweets", "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ Error al consumir de la cola: %v", err)
	}

	log.Println("📡 Esperando mensajes de tweets...")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("🟢 Recibido: %s", d.Body)
			tweet := Tweet{
				Body: string(d.Body),
				Time: time.Now(),
			}
			appendTweetToFile(tweet)
		}
	}()
	<-forever
}
