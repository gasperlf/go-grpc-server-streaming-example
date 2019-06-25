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

	cli := pb.NewLiveScoreClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	new := &pb.News{Type: getType(action), Min: int32(min), Team: team, Details: details}
	r, err := cli.PublishNews(ctx, &pb.PublishNewsRequest{News: new, Game: game})

	if err != nil {
		log.Printf("Error en la petición %v ", err)
	}

	log.Printf("Respuesta del servidor %v ", r)

}

func getType(action string) pb.TypeNews {

	typeNews := pb.TypeNews_UNKOWN

	switch action {
	case "goal":
		typeNews = pb.TypeNews_GOAL
	case "offside":
		typeNews = pb.TypeNews_OFFSIDE
	case "yellowcard":
		typeNews = pb.TypeNews_YELLOW_CARD
	case "redcard":
		typeNews = pb.TypeNews_RED_CARD
	case "finished":
		typeNews = pb.TypeNews_FINISHED
	}

	return typeNews

}
