package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "zetoslab.com/livescore/livescore"
)

const (
	addressCnx = "localhost:50051"
)

func main() {

	var game, team, details string
	var typeNew, min int

	flag.StringVar(&game, "game", "", "001")
	flag.StringVar(&team, "team", "", "Local")
	flag.StringVar(&details, "details", "", "")
	flag.IntVar(&typeNew, "typeNew", 0, "")
	flag.IntVar(&min, "min", 0, "")
	flag.Parse()

	conn, err := grpc.Dial(addressCnx, grpc.WithInsecure())

	if err != nil {

		log.Fatalf("Error conectando con servidor %v ", err)

	}
	defer conn.Close()

	cli := pb.NewLivescoreClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	new := &pb.New{Type: 1, Min: int32(min), Team: team, Details: details}
	r, err := cli.PublishNew(ctx, &pb.PublishNewRequest{New: new, Game: game})

	if err != nil {
		log.Printf("Error en la petición %v ", err)
	}

	log.Printf("Respuesta del servidor %v ", r)

}
