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

	var game, team, details, action string
	var min int

	flag.StringVar(&game, "game", "", "001")
	flag.StringVar(&team, "team", "", "Local")
	flag.StringVar(&details, "details", "", "")
	flag.StringVar(&action, "action", "", "")
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

	new := &pb.New{Type: getType(action), Min: int32(min), Team: team, Details: details}
	r, err := cli.PublishNew(ctx, &pb.PublishNewRequest{New: new, Game: game})

	if err != nil {
		log.Printf("Error en la petición %v ", err)
	}

	log.Printf("Respuesta del servidor %v ", r)

}

func getType(action string) pb.TypeNew {

	typeNews := pb.TypeNew_UNKOWN

	switch action {
	case "goal":
		typeNews = pb.TypeNew_GOAL
	case "offside":
		typeNews = pb.TypeNew_OFFSIDE
	case "yellowcard":
		typeNews = pb.TypeNew_YELLOW_CARD
	case "redcard":
		typeNews = pb.TypeNew_RED_CARD
	case "finished":
		typeNews = pb.TypeNew_FINISHED
	}

	return typeNews

}
