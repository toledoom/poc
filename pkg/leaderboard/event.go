package leaderboard

import (
	"errors"

	"github.com/toledoom/poc/pkg/framework"
	"github.com/toledoom/poc/pkg/player"
)

type PlayerScoreUpdateEventHandler struct {
	r Ranking
}

func NewPlayerScoreUpdateEventHandler(r Ranking) *PlayerScoreUpdateEventHandler {
	return &PlayerScoreUpdateEventHandler{
		r: r,
	}
}

func (eh *PlayerScoreUpdateEventHandler) Notify(evt framework.Event) error {
	pse, ok := evt.(player.PlayerScoreUpdated)
	if !ok {
		return errors.New("invalid error type. Want PlayerScoreUpdated")
	}
	err := eh.r.UpdateScore(pse.PlayerID(), pse.NewScore())
	return err
}
