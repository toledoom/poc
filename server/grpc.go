package server

import (
	"context"

	battlemodel "github.com/toledoom/poc/pkg/battle"
	"github.com/toledoom/poc/pkg/framework"
	leaderboardmodel "github.com/toledoom/poc/pkg/leaderboard"
	playermodel "github.com/toledoom/poc/pkg/player"
	"github.com/toledoom/poc/proto/battle"
	"github.com/toledoom/poc/proto/leaderboard"
	"github.com/toledoom/poc/proto/player"
)

type GameServer struct {
	br             battlemodel.BattleRepository
	pr             playermodel.PlayerRepository
	calculator     battlemodel.ScoreCalculator
	ranking        leaderboardmodel.Ranking
	eventPublisher *framework.EventPublisher

	battle.UnimplementedBattleServer
	leaderboard.UnimplementedLeaderboardServer
	player.UnimplementedPlayerServer
}

func (s *GameServer) StartBattle(ctx context.Context, sbr *battle.StartBattleRequest) (*battle.StartBattleResponse, error) {
	b := battlemodel.New(sbr.PlayerId1, sbr.PlayerId2)
	s.br.Add(b)
	return &battle.StartBattleResponse{
		BattleId: b.Id(),
	}, nil
}

func (s *GameServer) FinishBattle(ctx context.Context, fbr *battle.FinishBattleRequest) (*battle.FinishBattleResponse, error) {
	battleID := fbr.BattleId
	winnerID := fbr.WinnerId

	b, err := s.br.GetByID(battleID)
	if err != nil {
		return &battle.FinishBattleResponse{}, err
	}
	player1ID := b.Player1ID
	player2ID := b.Player2ID

	player1, err := s.pr.GetByID(player1ID)
	if err != nil {
		return &battle.FinishBattleResponse{}, err
	}
	player2, err := s.pr.GetByID(player2ID)
	if err != nil {
		return &battle.FinishBattleResponse{}, err
	}

	player1OldScore := player1.Score
	player2OldScore := player2.Score

	if winnerID == player1ID {
		player1.Score = s.calculator.Calculate(player1.Score, player2.Score, battlemodel.Win)
		player2.Score = s.calculator.Calculate(player2.Score, player1.Score, battlemodel.Lose)
	} else {
		player1.Score = s.calculator.Calculate(player1.Score, player2.Score, battlemodel.Lose)
		player2.Score = s.calculator.Calculate(player2.Score, player1.Score, battlemodel.Win)
	}

	err = s.pr.Update(player1)
	if err != nil {
		return nil, err
	}
	err = s.pr.Update(player2)
	if err != nil {
		return nil, err
	}

	err = s.eventPublisher.Notify(playermodel.NewPlayerScoreDispatched(player1ID, player1OldScore, player1.Score))
	if err != nil {
		return nil, err
	}
	err = s.eventPublisher.Notify(playermodel.NewPlayerScoreDispatched(player2ID, player2OldScore, player2.Score))
	if err != nil {
		return nil, err
	}

	return &battle.FinishBattleResponse{
		Player1Score: player1.Score,
		Player2Score: player2.Score,
	}, nil
}

func (s *GameServer) GetRank(ctx context.Context, grr *leaderboard.GetRankRequest) (*leaderboard.GetRankResponse, error) {
	playerID := grr.PlayerId
	rank, err := s.ranking.GetRank(playerID)
	return &leaderboard.GetRankResponse{
		Rank: rank,
	}, err
}

func (s *GameServer) GetTopPlayers(ctx context.Context, gtp *leaderboard.GetTopPlayersRequest) (*leaderboard.GetTopPlayersResponse, error) {
	limit := gtp.NumPlayers
	membersModel, err := s.ranking.GetTopPlayers(limit)
	if err != nil {
		return nil, err
	}
	var members []*leaderboard.Member
	for _, mm := range membersModel {
		m := &leaderboard.Member{
			Id:    mm.PlayerID,
			Score: mm.Score,
		}
		members = append(members, m)
	}
	return &leaderboard.GetTopPlayersResponse{
		MemberList: members,
	}, err
}

func (s *GameServer) CreatePlayer(ctx context.Context, cpr *player.CreatePlayerRequest) (*player.CreatePlayerResponse, error) {
	p := playermodel.New(cpr.Id, cpr.Name)
	err := s.pr.Add(p)
	return &player.CreatePlayerResponse{}, err
}

func (s *GameServer) GetPlayerById(ctx context.Context, cpr *player.GetPlayerByIdRequest) (*player.GetPlayerByIdResponse, error) {
	p, err := s.pr.GetByID(cpr.Id)
	return &player.GetPlayerByIdResponse{
		Name:  p.Name,
		Score: p.Score,
	}, err
}
