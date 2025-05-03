package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/SaulCerezo/TweetsClima/go-entry/github.com/SaulCerezo/TweetsClima/go-entry/proto"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedWeatherServiceServer
}

func (s *server) SendTweets(ctx context.Context, req *proto.TweetBatch) (*proto.Ack, error) {
	log.Printf("Received %d tweets", len(req.Tweets))

	// Variables de entorno
	rabbitUser := os.Getenv("RABBITMQ_USER")
	rabbitPass := os.Getenv("RABBITMQ_PASS")
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitPort := os.Getenv("RABBITMQ_PORT")

	// URL de conexión a RabbitMQ
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUser, rabbitPass, rabbitHost, rabbitPort)
	log.Println("🔌 Conectando a RabbitMQ en:", connStr)

	// Conectar a RabbitMQ
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Printf("❌ Error al conectar a RabbitMQ: %v", err)
		return nil, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("❌ Error al abrir canal: %v", err)
		return nil, err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare("tweets", false, false, false, false, nil)
	if err != nil {
		log.Printf("❌ Error al declarar cola: %v", err)
		return nil, err
	}

	for _, t := range req.Tweets {
		body := fmt.Sprintf("📨 Tweet => %s (%s) - %s", t.Description, t.Country, t.Weather)
		err = ch.Publish("", "tweets", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
		if err != nil {
			log.Printf("❌ Error al publicar tweet: %v", err)
		}
		log.Println("✅ Enviado a RabbitMQ:", body)
	}

	return &proto.Ack{
		Status: "received",
		Count:  int32(len(req.Tweets)),
	}, nil
}

func main() {
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	addr := fmt.Sprintf(":%s", grpcPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterWeatherServiceServer(s, &server{})

	fmt.Printf("🚀 Servidor gRPC escuchando en %s...\n", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
