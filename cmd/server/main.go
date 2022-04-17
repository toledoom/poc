package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/toledoom/poc/proto/battle"
	"github.com/toledoom/poc/proto/leaderboard"
	"github.com/toledoom/poc/proto/player"
	"github.com/toledoom/poc/server"
)

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gameServer := server.NewGrpc(os.Getenv("REDIS_ADDR"), os.Getenv("DYNAMO_ADDR"))

	grpcServer := grpc.NewServer()
	battle.RegisterBattleServer(grpcServer, gameServer)
	leaderboard.RegisterLeaderboardServer(grpcServer, gameServer)
	player.RegisterPlayerServer(grpcServer, gameServer)
	log.Print("Server started")
	grpcServer.Serve(lis)
}
