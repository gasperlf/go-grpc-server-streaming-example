package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "zetoslab.com/livescore/livescore"
)

const (
	port = ":50051"
)

/* var newsToSend []*pb.New
var indexNewToSend = -1 */
var games []*pb.Game                       // List of availables games
var newsGames = make(map[string][]*pb.New) // List of news for each game
var indexNewGames = make(map[string]int)   // Index of the last new

type server struct{} // Definir un struct donde mplementaremos todos los métodos

func (s *server) GetGamesList(ctx context.Context, in *pb.GetGamesListRequest) (*pb.GetGamesListResponse, error) {

	log.Printf("Req List of games for %v", in.Country)

	response := &pb.GetGamesListResponse{Games: games}
	return response, nil

}

func (s *server) PublishNew(ctx context.Context, in *pb.PublishNewRequest) (*pb.PublishNewResponse, error) {

	log.Printf("Payload: %v", in)

	/* newsToSend = append(newsToSend, in.New) */
	newsGames[in.Game] = append(newsGames[in.Game], in.New)
	log.Printf("Noticias del partido %v: ", newsGames[in.Game])

	response := &pb.PublishNewResponse{Ok: true}
	return response, nil

}

func (s *server) GetNewsGame(in *pb.GetNewsGameRequest, stream pb.Livescore_GetNewsGameServer) error {

	for {

		if len(newsGames[in.GameId]) > 0 {

			log.Printf("len(newsGames[%v]): %v", in.GameId, len(newsGames[in.GameId]))
			if len(newsGames[in.GameId]) > indexNewGames[in.GameId]+1 {

				indexNewGames[in.GameId]++
				err := stream.Send(&pb.GetNewsGameResponse{New: newsGames[in.GameId][indexNewGames[in.GameId]]})
				if err != nil {
					log.Printf("Error sending new %v ", err)
				}

				if newsGames[in.GameId][indexNewGames[in.GameId]].Type == 5 {
					log.Printf("Streaming connection has finished because the match has finished.")
					break
				}

			}

		}

	}

	return nil

}

func main() {

	games = append(games, &pb.Game{Id: "001", TeamLocal: "Python", TeamVisitor: "Node", Country: "Co"})
	indexNewGames["001"] = -1
	newsGames["001"] = append(newsGames["001"], &pb.New{Type: 0, Min: 0, Team: "", Details: "Game is starting"})

	games = append(games, &pb.Game{Id: "002", TeamLocal: "CPlusPlus", TeamVisitor: "C", Country: "Co"})
	indexNewGames["002"] = -1
	newsGames["002"] = append(newsGames["002"], &pb.New{Type: 0, Min: 0, Team: "", Details: "Game is starting"})

	games = append(games, &pb.Game{Id: "003", TeamLocal: "CSharp", TeamVisitor: "Java", Country: "Co"})
	indexNewGames["003"] = -1
	newsGames["003"] = append(newsGames["003"], &pb.New{Type: 0, Min: 0, Team: "", Details: "Game is starting"})

	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Error al exponer el puerto %v ", err)
	}

	s := grpc.NewServer()
	pb.RegisterLivescoreServer(s, &server{})

	log.Printf("Server listening by %v ", port)
	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("Error al iniciar servidor %v ", err)
	}

}
