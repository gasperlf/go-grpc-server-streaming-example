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

	conn, err := grpc.Dial(addressCnx, grpc.WithInsecure()) // Establecemos el canal para comunicarnos

	if err != nil {

		log.Fatalf("Error conectando con servidor %v ", err)

	}
	defer conn.Close()

	cli := pb.NewLiveScoreClient(conn) // Creamos nuestro servicios / stub gRPC

	ctx, cancel := context.WithTimeout(context.Background(), time.Second) // Definimos contexto con timeout
	defer cancel()

	r, err := cli.GetGamesList(ctx, &pb.GetGamesListRequest{Country: "CO"}) // Usamos nuestro cliente

	if err != nil {
		log.Printf("Error en la petición %v ", err)
	}

	log.Printf("Respuesta del servidor %v ", r)

}
