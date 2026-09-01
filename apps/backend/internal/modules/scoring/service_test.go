package scoring

import (
	"context"
	"testing"
	"time"
)

// --- fakes ---------------------------------------------------------

type fakeRepo struct {
	wallets map[string]Wallet
	streaks map[string]Streak
	txs     []Transaction
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{wallets: map[string]Wallet{}, streaks: map[string]Streak{}}
}

func (f *fakeRepo) GetWallet(_ context.Context, userID string) (Wallet, error) {
	if w, ok := f.wallets[userID]; ok {
		return w, nil
	}
	return Wallet{UserID: userID, Level: 1}, nil
}

func (f *fakeRepo) SaveWallet(_ context.Context, w Wallet) error {
	f.wallets[w.UserID] = w
	return nil
}

func (f *fakeRepo) AddTransaction(_ context.Context, t Transaction) error {
	f.txs = append(f.txs, t)
	return nil
}

func (f *fakeRepo) GetStreak(_ context.Context, userID string) (Streak, error) {
	if s, ok := f.streaks[userID]; ok {
		return s, nil
	}
	return Streak{UserID: userID}, nil
}

func (f *fakeRepo) SaveStreak(_ context.Context, s Streak) error {
	f.streaks[s.UserID] = s
	return nil
}

func (f *fakeRepo) Leaderboard(_ context.Context, limit int) ([]Wallet, error) {
	var out []Wallet
	for _, w := range f.wallets {
		out = append(out, w)
	}
	// urut poin turun (repo asli ORDER BY); stabilkan lewat service.
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].TotalPoints > out[i].TotalPoints {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	if limit < len(out) {
		out = out[:limit]
	}
	return out, nil
}

type fakeDir struct{ names map[string]string }

func (d fakeDir) NamesByIDs(_ context.Context, ids []string) (map[string]string, error) {
	out := map[string]string{}
	for _, id := range ids {
		if n, ok := d.names[id]; ok {
			out[id] = n
		}
	}
	return out, nil
}

func newSvc(names map[string]string) (*service, *fakeRepo) {
	repo := newFakeRepo()
	return newService(repo, fakeDir{names: names}), repo
}

func dt(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// --- tests --------------------------------------------------------

func TestOnQuestCompleted(t *testing.T) {
	svc, repo := newSvc(nil)
	ctx := context.Background()

	if err := svc.OnQuestCompleted(ctx, "u1", "q1", 10, dt(2026, 9, 1)); err != nil {
		t.Fatalf("OnQuestCompleted: %v", err)
	}

	w := repo.wallets["u1"]
	if w.TotalPoints != 10 || w.XP != 10 || w.Level != 1 {
		t.Fatalf("wallet = %+v, mau {10,10,1}", w)
	}
	st := repo.streaks["u1"]
	if st.Current != 1 || st.Longest != 1 {
		t.Fatalf("streak = %+v, mau current/longest 1", st)
	}
	if len(repo.txs) != 1 || repo.txs[0].Points != 10 {
		t.Fatalf("txs = %+v, mau 1 transaksi +10", repo.txs)
	}

	// selesaikan lagi besok, cukup poin untuk naik level
	if err := svc.OnQuestCompleted(ctx, "u1", "q2", 95, dt(2026, 9, 2)); err != nil {
		t.Fatalf("complete kedua: %v", err)
	}
	w = repo.wallets["u1"]
	if w.XP != 105 || w.Level != 2 {
		t.Fatalf("wallet = %+v, mau XP 105 level 2", w)
	}
	if repo.streaks["u1"].Current != 2 {
		t.Fatalf("streak current = %d, mau 2", repo.streaks["u1"].Current)
	}
}

func TestOnQuestUncompleted_RollsBackPointsNotStreak(t *testing.T) {
	svc, repo := newSvc(nil)
	ctx := context.Background()
	d := dt(2026, 9, 1)

	svc.OnQuestCompleted(ctx, "u1", "q1", 20, d)
	streakBefore := repo.streaks["u1"]

	if err := svc.OnQuestUncompleted(ctx, "u1", "q1", 20, d); err != nil {
		t.Fatalf("OnQuestUncompleted: %v", err)
	}

	w := repo.wallets["u1"]
	if w.TotalPoints != 0 || w.XP != 0 || w.Level != 1 {
		t.Fatalf("wallet = %+v, mau kembali {0,0,1}", w)
	}
	if len(repo.txs) != 2 || repo.txs[1].Points != -20 {
		t.Fatalf("txs = %+v, mau transaksi kedua -20", repo.txs)
	}
	// ADR-009: streak tidak diputar balik
	if repo.streaks["u1"] != streakBefore {
		t.Fatalf("streak berubah setelah uncomplete: %+v -> %+v (harus tetap, ADR-009)",
			streakBefore, repo.streaks["u1"])
	}
}

func TestOnQuestUncompleted_ClampsAtZero(t *testing.T) {
	svc, repo := newSvc(nil)
	ctx := context.Background()

	svc.OnQuestCompleted(ctx, "u1", "q1", 5, dt(2026, 9, 1))
	// uncomplete dengan poin lebih besar dari saldo -> clamp di 0, bukan negatif
	if err := svc.OnQuestUncompleted(ctx, "u1", "q1", 50, dt(2026, 9, 1)); err != nil {
		t.Fatalf("uncomplete: %v", err)
	}
	w := repo.wallets["u1"]
	if w.TotalPoints != 0 || w.XP != 0 {
		t.Fatalf("wallet = %+v, mau clamp di 0", w)
	}
}

func TestLeaderboard_NamesRankAndFallback(t *testing.T) {
	svc, repo := newSvc(map[string]string{"u1": "Alice", "u2": "Bob"})
	ctx := context.Background()

	repo.wallets["u1"] = Wallet{UserID: "u1", TotalPoints: 30, Level: 1}
	repo.wallets["u2"] = Wallet{UserID: "u2", TotalPoints: 50, Level: 1}
	repo.wallets["u3"] = Wallet{UserID: "u3", TotalPoints: 10, Level: 1} // tak ada di directory

	board, err := svc.Leaderboard(ctx, 20)
	if err != nil {
		t.Fatalf("Leaderboard: %v", err)
	}
	if len(board) != 3 {
		t.Fatalf("len = %d, mau 3", len(board))
	}
	if board[0].UserID != "u2" || board[0].Rank != 1 || board[0].DisplayName != "Bob" || board[0].Points != 50 {
		t.Fatalf("rank 1 salah: %+v", board[0])
	}
	if board[1].DisplayName != "Alice" || board[1].Rank != 2 {
		t.Fatalf("rank 2 salah: %+v", board[1])
	}
	if board[2].UserID != "u3" || board[2].DisplayName != deletedUserName || board[2].Rank != 3 {
		t.Fatalf("rank 3 (fallback) salah: %+v", board[2])
	}
}

func TestGetScoreAndStreak_Defaults(t *testing.T) {
	svc, _ := newSvc(nil)
	ctx := context.Background()

	sc, err := svc.GetScore(ctx, "baru")
	if err != nil {
		t.Fatalf("GetScore: %v", err)
	}
	if sc.Level != 1 || sc.TotalPoints != 0 || sc.PointsToNextLevel != 100 {
		t.Fatalf("score default = %+v", sc)
	}

	st, err := svc.GetStreak(ctx, "baru")
	if err != nil {
		t.Fatalf("GetStreak: %v", err)
	}
	if st.Current != 0 || st.LastActive != nil {
		t.Fatalf("streak default = %+v, mau current 0 & last_active null", st)
	}
}
