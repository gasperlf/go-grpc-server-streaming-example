package main

import (
	"context"
	"log"
	"net"
	"time"

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
var indexNewGamesByConn = make([]int, 1)   // Array of Index. For each connection.

type server struct{} // Definir una interfaz donde implementaremos todos los métodos de nuestra definición

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

	log.Printf("Connection to game %v ", in.GameId)
	indexNewGamesByConn = append(indexNewGamesByConn, indexNewGames[in.GameId])
	internalIndex := len(indexNewGamesByConn) - 1

	// log.Printf("indexNewGamesByConn: %v ", indexNewGamesByConn)

	for {

		if len(newsGames[in.GameId]) > 0 {

			time.Sleep(2 * time.Second)
			if len(newsGames[in.GameId]) > indexNewGames[in.GameId]+1 {

				indexNewGames[in.GameId]++ // Updating the general index of a new-game

			}

			if len(newsGames[in.GameId]) > indexNewGamesByConn[internalIndex]+1 {

				indexNewGamesByConn[internalIndex]++
				err := stream.Send(&pb.GetNewsGameResponse{New: newsGames[in.GameId][indexNewGamesByConn[internalIndex]]})
				if err != nil {
					log.Printf("Error sending new %v ", err)
				}

				if newsGames[in.GameId][indexNewGamesByConn[internalIndex]].Type == 5 {
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

	listener, err := net.Listen("tcp", port) // Definimos que puerto usaremos para epxoner nuestro servicio

	if err != nil {
		log.Fatalf("Error al exponer el puerto %v ", err)
	}

	s := grpc.NewServer()                    // Creamos servidor gRPC
	pb.RegisterLivescoreServer(s, &server{}) // Registramos nuestro servidor para Livescore pasando los métodos implementados

	log.Printf("Server listening by %v ", port)
	err = s.Serve(listener) // Exponemos nuestros servicios
	if err != nil {
		log.Fatalf("Error al iniciar servidor %v ", err)
	}

}
