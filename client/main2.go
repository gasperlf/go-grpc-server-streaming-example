package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "zetoslab.com/livescore/livescore"
)

const (
	addressCnx = "localhost:50051"
)

func main() {

	var game string

	flag.StringVar(&game, "game", "", "001")
	flag.Parse()

	conn, err := grpc.Dial(addressCnx, grpc.WithInsecure())

	if err != nil {

		log.Fatalf("Error conectando con servidor %v ", err)

	}
	defer conn.Close()

	cli := pb.NewLivescoreClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	log.Printf("Will be subscribed to game: ", game)
	newReq := &pb.GetNewsGameRequest{GameId: game}

	stream, err := cli.GetNewsGame(ctx, newReq) // Creamos nuestro cliente de streaming

	if err != nil {
		log.Printf("Error en la petición %v ", err)
	}

	for {

		new, err := stream.Recv() // Recibimos los datos que lleguen dle servidor

		if err != nil {

			if err == io.EOF {
				log.Printf("Game finalizado")
				break
			} else {
				log.Fatalf("Error recibiendo new %v ", err)
			}

		}

		log.Printf("New recibida %v ", new)
	}

}
