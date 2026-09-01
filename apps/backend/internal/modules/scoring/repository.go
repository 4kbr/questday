package scoring

import "context"

// Repository = kontrak penyimpanan scoring. Implementasi Postgres ada di
// repository_postgres.go — satu-satunya tempat SQL module ini. Semua query
// terikat user_id kecuali Leaderboard.
type Repository interface {
	GetWallet(ctx context.Context, userID string) (Wallet, error) // default {0,0,1} bila belum ada
	SaveWallet(ctx context.Context, w Wallet) error               // UPSERT
	AddTransaction(ctx context.Context, t Transaction) error
	GetStreak(ctx context.Context, userID string) (Streak, error) // default kosong bila belum ada
	SaveStreak(ctx context.Context, s Streak) error               // UPSERT
	Leaderboard(ctx context.Context, limit int) ([]Wallet, error) // user_id + total_points saja, TANPA JOIN users (ADR-014)
}
