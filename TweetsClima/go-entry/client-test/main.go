package main

import (
	"context"
	"fmt"
	"log"

	"github.com/SaulCerezo/TweetsClima/go-entry/github.com/SaulCerezo/TweetsClima/go-entry/proto" // ← usa el import local correcto

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor gRPC: %v", err)
	}
	defer conn.Close()

	client := proto.NewWeatherServiceClient(conn)

	req := &proto.TweetBatch{
		Tweets: []*proto.Tweet{
			{Description: "Desde el cliente", Country: "GT", Weather: "Nublado"},
		},
	}

	res, err := client.SendTweets(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al enviar los tweets: %v", err)
	}

	fmt.Printf("Respuesta: %v\n", res)
}
