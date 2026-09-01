package server

import (
	"context"
	"time"
)

// noopScorer memenuhi quest.ScoreAwarder tanpa efek apa pun. Dipakai sementara
// di Phase 2 supaya module quest bisa dirakit & diuji tanpa module scoring.
//
// TODO(T3.9): hapus, ganti dengan scoringMod.AsScoreAwarder() di server.New.
type noopScorer struct{}

func (noopScorer) OnQuestCompleted(_ context.Context, _, _ string, _ int, _ time.Time) error {
	return nil
}

func (noopScorer) OnQuestUncompleted(_ context.Context, _, _ string, _ int, _ time.Time) error {
	return nil
}
