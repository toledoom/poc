package leaderboard

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Member struct {
	PlayerID string
	Score    int64
}

type Ranking interface {
	GetRank(playerID string) (int64, error)
	GetTopPlayers(limit int64) ([]Member, error)
	UpdateScore(playerID string, score int64) error
}

type RedisRanking struct {
	c    *redis.Client
	name string
}

func NewRedisRanking(c *redis.Client, name string) *RedisRanking {
	return &RedisRanking{
		c:    c,
		name: name,
	}
}

func (r *RedisRanking) GetRank(playerID string) (int64, error) {
	zRank := r.c.ZRank(context.TODO(), r.name, playerID)
	if zRank.Err() != redis.Nil {
		return 0, zRank.Err()
	}

	return zRank.Val(), nil
}

func (r *RedisRanking) GetTopPlayers(limit int64) ([]Member, error) {
	zRevRange, err := r.c.ZRevRangeWithScores(context.TODO(), r.name, 0, limit).Result()
	if err != nil {
		return []Member{}, err
	}

	var members []Member
	for _, item := range zRevRange {
		m := Member{
			PlayerID: fmt.Sprintf("%v", item.Member),
			Score:    int64(item.Score),
		}
		members = append(members, m)
	}

	return members, nil
}

func (r *RedisRanking) UpdateScore(playerID string, score int64) error {
	member := &redis.Z{
		Member: playerID,
		Score:  float64(score),
	}
	zAdd := r.c.ZAdd(context.TODO(), r.name, member)

	if zAdd.Err() != redis.Nil {
		return zAdd.Err()
	}

	return nil
}
