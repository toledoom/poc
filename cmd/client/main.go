package main

import (
	"context"
	"fmt"
	"os"

	"github.com/toledoom/poc/proto/battle"
	"github.com/toledoom/poc/proto/leaderboard"
	"github.com/toledoom/poc/proto/player"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	serverAddr := os.Getenv("REMOTE_ADDR")
	if serverAddr == "" {
		serverAddr = "localhost:50051"
	}
	conn, err := grpc.Dial(serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ctx := context.TODO()

	playerClient := player.NewPlayerClient(conn)
	playerID := "player-toledo-id"
	cpReq := &player.CreatePlayerRequest{
		Id:   playerID,
		Name: "Toledoom",
	}
	_, err = playerClient.CreatePlayer(ctx, cpReq)
	if err != nil {
		panic(err)
	}

	opponentID := "evil-guy-id"
	coReq := &player.CreatePlayerRequest{
		Id:   opponentID,
		Name: "Evil guy",
	}
	_, err = playerClient.CreatePlayer(ctx, coReq)
	if err != nil {
		panic(err)
	}

	battleClient := battle.NewBattleClient(conn)
	sbReq := &battle.StartBattleRequest{
		PlayerId1: playerID,
		PlayerId2: opponentID,
	}
	sbResp, err := battleClient.StartBattle(ctx, sbReq)
	if err != nil {
		panic(err)
	}
	battleID := sbResp.BattleId

	fbReq := &battle.FinishBattleRequest{
		BattleId: battleID,
		WinnerId: playerID,
	}
	fbResp, err := battleClient.FinishBattle(ctx, fbReq)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Player score: %d\n", fbResp.Player1Score)
	fmt.Printf("Opponent score: %d\n", fbResp.Player2Score)

	leaderboardClient := leaderboard.NewLeaderboardClient(conn)
	gtpReq := &leaderboard.GetTopPlayersRequest{
		NumPlayers: 3,
	}
	gtpResp, err := leaderboardClient.GetTopPlayers(ctx, gtpReq)
	if err != nil {
		panic(err)
	}

	fmt.Println("LEADERBOARDS")
	for _, m := range gtpResp.MemberList {
		fmt.Printf("PlayerID: %s - Score:%d\n", m.Id, m.Score)
	}
}
