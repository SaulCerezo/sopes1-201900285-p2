package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

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

func main() {
	conn := connectRabbitMQ()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Error al abrir canal: %v", err)
	}
	defer ch.Close()

	// Asegurar que la cola exista
	_, err = ch.QueueDeclare(
		"tweets", // nombre de la cola
		false,    // durable
		false,    // auto-delete
		false,    // exclusive
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		log.Fatalf("❌ Error al declarar la cola: %v", err)
	}

	msgs, err := ch.Consume(
		"tweets", // nombre de la cola
		"",       // consumer tag
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		log.Fatalf("❌ Error al consumir de la cola: %v", err)
	}

	log.Println("📡 Esperando mensajes de tweets...")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("🟢 Recibido: %s", d.Body)
		}
	}()
	<-forever
}
