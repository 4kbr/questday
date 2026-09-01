package scoring

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// postgresRepository = implementasi Repository di atas *sql.DB (ADR-012).
// SATU-SATUNYA tempat SQL module scoring. TIDAK pernah menyentuh tabel `users`
// (ADR-014) — nama untuk leaderboard diisi service lewat port UserDirectory.
type postgresRepository struct {
	db *sql.DB
}

func newPostgresRepository(db *sql.DB) *postgresRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetWallet(ctx context.Context, userID string) (Wallet, error) {
	const q = `SELECT user_id, total_points, xp, level FROM wallets WHERE user_id = $1`
	var w Wallet
	err := r.db.QueryRowContext(ctx, q, userID).Scan(&w.UserID, &w.TotalPoints, &w.XP, &w.Level)
	if errors.Is(err, sql.ErrNoRows) {
		return Wallet{UserID: userID, TotalPoints: 0, XP: 0, Level: 1}, nil
	}
	if err != nil {
		return Wallet{}, fmt.Errorf("scoring: get wallet: %w", err)
	}
	return w, nil
}

func (r *postgresRepository) SaveWallet(ctx context.Context, w Wallet) error {
	const stmt = `
		INSERT INTO wallets (user_id, total_points, xp, level)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE
		SET total_points = EXCLUDED.total_points, xp = EXCLUDED.xp, level = EXCLUDED.level`
	if _, err := r.db.ExecContext(ctx, stmt, w.UserID, w.TotalPoints, w.XP, w.Level); err != nil {
		return fmt.Errorf("scoring: save wallet: %w", err)
	}
	return nil
}

func (r *postgresRepository) AddTransaction(ctx context.Context, t Transaction) error {
	const stmt = `
		INSERT INTO point_transactions (id, user_id, quest_id, points, date)
		VALUES ($1, $2, $3, $4, $5)`
	var questID any
	if t.QuestID != "" {
		questID = t.QuestID
	}
	if _, err := r.db.ExecContext(ctx, stmt, t.ID, t.UserID, questID, t.Points, t.Date); err != nil {
		return fmt.Errorf("scoring: add transaction: %w", err)
	}
	return nil
}

func (r *postgresRepository) GetStreak(ctx context.Context, userID string) (Streak, error) {
	const q = `SELECT user_id, current, longest, last_active FROM streaks WHERE user_id = $1`
	var (
		s          Streak
		lastActive sql.NullTime
	)
	err := r.db.QueryRowContext(ctx, q, userID).Scan(&s.UserID, &s.Current, &s.Longest, &lastActive)
	if errors.Is(err, sql.ErrNoRows) {
		return Streak{UserID: userID}, nil
	}
	if err != nil {
		return Streak{}, fmt.Errorf("scoring: get streak: %w", err)
	}
	if lastActive.Valid {
		s.LastActive = lastActive.Time
	}
	return s, nil
}

func (r *postgresRepository) SaveStreak(ctx context.Context, s Streak) error {
	const stmt = `
		INSERT INTO streaks (user_id, current, longest, last_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE
		SET current = EXCLUDED.current, longest = EXCLUDED.longest, last_active = EXCLUDED.last_active`
	var lastActive any
	if !s.LastActive.IsZero() {
		lastActive = s.LastActive
	}
	if _, err := r.db.ExecContext(ctx, stmt, s.UserID, s.Current, s.Longest, lastActive); err != nil {
		return fmt.Errorf("scoring: save streak: %w", err)
	}
	return nil
}

func (r *postgresRepository) Leaderboard(ctx context.Context, limit int) ([]Wallet, error) {
	const q = `
		SELECT user_id, total_points
		FROM wallets
		ORDER BY total_points DESC, user_id
		LIMIT $1`
	rows, err := r.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, fmt.Errorf("scoring: leaderboard: %w", err)
	}
	defer rows.Close()

	var out []Wallet
	for rows.Next() {
		var w Wallet
		if err := rows.Scan(&w.UserID, &w.TotalPoints); err != nil {
			return nil, fmt.Errorf("scoring: scan leaderboard: %w", err)
		}
		out = append(out, w)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scoring: leaderboard rows: %w", err)
	}
	return out, nil
}
