package player

import "github.com/toledoom/poc/pkg/framework"

type PlayerScoreUpdated interface {
	framework.Event
	PlayerID() string
	OldScore() int64
	NewScore() int64
}

type PlayerScoreDispatched struct {
	playerID string
	oldScore int64
	newScore int64
}

func NewPlayerScoreDispatched(playerID string, oldScore, newScore int64) PlayerScoreUpdated {
	return &PlayerScoreDispatched{
		playerID: playerID,
		oldScore: oldScore,
		newScore: newScore,
	}
}

func (evt *PlayerScoreDispatched) Name() string {
	return "PlayerScoreUpdated"
}

func (evt *PlayerScoreDispatched) PlayerID() string {
	return evt.playerID
}
func (evt *PlayerScoreDispatched) OldScore() int64 {
	return evt.oldScore
}
func (evt *PlayerScoreDispatched) NewScore() int64 {
	return evt.newScore
}
