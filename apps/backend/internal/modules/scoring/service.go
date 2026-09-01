package scoring

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
)

// UserDirectory = port KELUAR milik scoring (ADR-014). Diimplementasi module
// user, disuntik server. scoring TIDAK meng-import user — leaderboard butuh
// nama tampilan tapi SQL scoring tak boleh JOIN ke tabel `users`.
type UserDirectory interface {
	NamesByIDs(ctx context.Context, ids []string) (map[string]string, error)
}

// deletedUserName dipakai bila UserDirectory tak menemukan sebuah userID
// (mis. user sudah dihapus) — entri leaderboard tetap ditampilkan.
const deletedUserName = "(pengguna dihapus)"

// service = logika gamifikasi + implementasi port quest.ScoreAwarder.
// Kecocokan signature dengan quest.ScoreAwarder diverifikasi compile-time saat
// server menyuntik (server tak akan build kalau tak cocok).
type service struct {
	repo Repository
	dir  UserDirectory
}

func newService(repo Repository, dir UserDirectory) *service {
	return &service{repo: repo, dir: dir}
}

// OnQuestCompleted menambah poin/XP, menghitung ulang level, menaikkan streak,
// dan mencatat transaksi (+points). Memenuhi quest.ScoreAwarder.
func (s *service) OnQuestCompleted(ctx context.Context, userID, questID string, points int, date time.Time) error {
	w, err := s.repo.GetWallet(ctx, userID)
	if err != nil {
		return err
	}
	w.TotalPoints += points
	w.XP += points
	w.Level = LevelForXP(w.XP)
	if err := s.repo.SaveWallet(ctx, w); err != nil {
		return err
	}

	st, err := s.repo.GetStreak(ctx, userID)
	if err != nil {
		return err
	}
	if err := s.repo.SaveStreak(ctx, NextStreak(st, date)); err != nil {
		return err
	}

	return s.addTx(ctx, userID, questID, points, date)
}

// OnQuestUncompleted mengembalikan poin/XP (clamp >= 0), menghitung ulang level,
// dan mencatat transaksi (-points). Streak SENGAJA tidak diputar balik (ADR-009):
// rekonstruksi urutan hari yang benar terlalu rumit untuk MVP.
func (s *service) OnQuestUncompleted(ctx context.Context, userID, questID string, points int, date time.Time) error {
	w, err := s.repo.GetWallet(ctx, userID)
	if err != nil {
		return err
	}
	w.TotalPoints = max0(w.TotalPoints - points)
	w.XP = max0(w.XP - points)
	w.Level = LevelForXP(w.XP)
	if err := s.repo.SaveWallet(ctx, w); err != nil {
		return err
	}

	return s.addTx(ctx, userID, questID, -points, date)
}

func (s *service) addTx(ctx context.Context, userID, questID string, points int, date time.Time) error {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("scoring: generate tx id: %w", err)
	}
	return s.repo.AddTransaction(ctx, Transaction{
		ID:      id.String(),
		UserID:  userID,
		QuestID: questID,
		Points:  points,
		Date:    date,
	})
}

// GetScore mengembalikan ringkasan skor user.
func (s *service) GetScore(ctx context.Context, userID string) (ScoreResponse, error) {
	w, err := s.repo.GetWallet(ctx, userID)
	if err != nil {
		return ScoreResponse{}, err
	}
	return toScoreResponse(w), nil
}

// GetStreak mengembalikan data streak user.
func (s *service) GetStreak(ctx context.Context, userID string) (StreakResponse, error) {
	st, err := s.repo.GetStreak(ctx, userID)
	if err != nil {
		return StreakResponse{}, err
	}
	return toStreakResponse(st), nil
}

// Leaderboard mengembalikan peringkat berdasarkan poin, lengkap dengan nama
// tampilan (dari UserDirectory) dan rank 1-based.
func (s *service) Leaderboard(ctx context.Context, limit int) ([]LeaderboardEntry, error) {
	wallets, err := s.repo.Leaderboard(ctx, limit)
	if err != nil {
		return nil, err
	}
	if len(wallets) == 0 {
		return []LeaderboardEntry{}, nil
	}

	ids := make([]string, 0, len(wallets))
	for _, w := range wallets {
		ids = append(ids, w.UserID)
	}
	names, err := s.dir.NamesByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("scoring: resolve names: %w", err)
	}

	// repo sudah ORDER BY total_points DESC; jaga urutan tetap stabil.
	sort.SliceStable(wallets, func(i, j int) bool {
		return wallets[i].TotalPoints > wallets[j].TotalPoints
	})

	out := make([]LeaderboardEntry, 0, len(wallets))
	for i, w := range wallets {
		name := names[w.UserID]
		if name == "" {
			name = deletedUserName
		}
		out = append(out, LeaderboardEntry{
			Rank:        i + 1,
			UserID:      w.UserID,
			DisplayName: name,
			Points:      w.TotalPoints,
		})
	}
	return out, nil
}

func max0(n int) int {
	if n < 0 {
		return 0
	}
	return n
}
