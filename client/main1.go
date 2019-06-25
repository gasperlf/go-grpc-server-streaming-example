package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "zetoslab.com/livescore/livescore"
)

const (
	addressCnx = "localhost:50051"
)

func main() {

	conn, err := grpc.Dial(addressCnx, grpc.WithInsecure())

	if err != nil {

		log.Fatalf("Error conectando con servidor %v ", err)

	}
	defer conn.Close()

	cli := pb.NewLivescoreClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := cli.GetGamesList(ctx, &pb.GetGamesListRequest{Country: "CO"})

	if err != nil {
		log.Printf("Error en la petición %v ", err)
	}

	log.Printf("Respuesta del servidor %v ", r)

}
