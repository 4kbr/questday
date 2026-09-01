package quest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

const pgUniqueViolation = "23505"

// postgresRepository = implementasi Repository di atas *sql.DB (ADR-012).
// SATU-SATUNYA tempat SQL module quest.
type postgresRepository struct {
	db *sql.DB
}

func newPostgresRepository(db *sql.DB) *postgresRepository {
	return &postgresRepository{db: db}
}

// --- Quest ---------------------------------------------------------------

func (r *postgresRepository) CreateQuest(ctx context.Context, q Quest) error {
	const stmt = `
		INSERT INTO quests (id, user_id, title, note, category, difficulty, recurrence, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, stmt,
		q.ID, q.UserID, q.Title, q.Note, string(q.Category), string(q.Difficulty), string(q.Recurrence), q.Active)
	if err != nil {
		return fmt.Errorf("quest: create quest: %w", err)
	}
	return nil
}

func (r *postgresRepository) GetQuest(ctx context.Context, userID, id string) (Quest, error) {
	const q = `
		SELECT id, user_id, title, note, category, difficulty, recurrence, active, created_at
		FROM quests WHERE id = $1 AND user_id = $2`
	return scanQuest(r.db.QueryRowContext(ctx, q, id, userID))
}

func (r *postgresRepository) ListQuestsByUser(ctx context.Context, userID string) ([]Quest, error) {
	const q = `
		SELECT id, user_id, title, note, category, difficulty, recurrence, active, created_at
		FROM quests WHERE user_id = $1 ORDER BY created_at`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("quest: list quests: %w", err)
	}
	defer rows.Close()

	var out []Quest
	for rows.Next() {
		qt, err := scanQuest(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, qt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("quest: list quests rows: %w", err)
	}
	return out, nil
}

func (r *postgresRepository) UpdateQuest(ctx context.Context, q Quest) error {
	const stmt = `
		UPDATE quests
		SET title = $3, note = $4, category = $5, difficulty = $6, recurrence = $7, active = $8
		WHERE id = $1 AND user_id = $2`
	res, err := r.db.ExecContext(ctx, stmt,
		q.ID, q.UserID, q.Title, q.Note, string(q.Category), string(q.Difficulty), string(q.Recurrence), q.Active)
	if err != nil {
		return fmt.Errorf("quest: update quest: %w", err)
	}
	return rowsAffectedOrNotFound(res)
}

func (r *postgresRepository) ArchiveQuest(ctx context.Context, userID, id string) error {
	const stmt = `UPDATE quests SET active = false WHERE id = $1 AND user_id = $2`
	res, err := r.db.ExecContext(ctx, stmt, id, userID)
	if err != nil {
		return fmt.Errorf("quest: archive quest: %w", err)
	}
	return rowsAffectedOrNotFound(res)
}

// --- QuestLog ----------------------------------------------------------

func (r *postgresRepository) CreateLog(ctx context.Context, l QuestLog) error {
	const stmt = `
		INSERT INTO quest_logs (id, quest_id, user_id, date, status, points_awarded)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, stmt,
		l.ID, l.QuestID, l.UserID, l.Date, string(l.Status), l.PointsAwarded)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return ErrAlreadyCompleted
		}
		return fmt.Errorf("quest: create log: %w", err)
	}
	return nil
}

func (r *postgresRepository) GetLog(ctx context.Context, userID, questID string, date time.Time) (QuestLog, error) {
	const q = `
		SELECT id, quest_id, user_id, date, status, points_awarded, completed_at
		FROM quest_logs WHERE quest_id = $1 AND user_id = $2 AND date = $3`
	return scanLog(r.db.QueryRowContext(ctx, q, questID, userID, date))
}

func (r *postgresRepository) DeleteLog(ctx context.Context, userID, questID string, date time.Time) error {
	const stmt = `DELETE FROM quest_logs WHERE quest_id = $1 AND user_id = $2 AND date = $3`
	res, err := r.db.ExecContext(ctx, stmt, questID, userID, date)
	if err != nil {
		return fmt.Errorf("quest: delete log: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("quest: delete log rows: %w", err)
	}
	if n == 0 {
		return ErrNotCompleted
	}
	return nil
}

func (r *postgresRepository) ListLogsByUserAndDate(ctx context.Context, userID string, date time.Time) ([]QuestLog, error) {
	const q = `
		SELECT id, quest_id, user_id, date, status, points_awarded, completed_at
		FROM quest_logs WHERE user_id = $1 AND date = $2`
	rows, err := r.db.QueryContext(ctx, q, userID, date)
	if err != nil {
		return nil, fmt.Errorf("quest: list logs: %w", err)
	}
	defer rows.Close()

	var out []QuestLog
	for rows.Next() {
		l, err := scanLog(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("quest: list logs rows: %w", err)
	}
	return out, nil
}

func (r *postgresRepository) ListActiveDates(ctx context.Context, userID string, from, to time.Time) ([]time.Time, error) {
	const q = `
		SELECT DISTINCT date FROM quest_logs
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date`
	rows, err := r.db.QueryContext(ctx, q, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("quest: list active dates: %w", err)
	}
	defer rows.Close()

	var out []time.Time
	for rows.Next() {
		var d time.Time
		if err := rows.Scan(&d); err != nil {
			return nil, fmt.Errorf("quest: scan active date: %w", err)
		}
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("quest: list active dates rows: %w", err)
	}
	return out, nil
}

// --- helper scan ------------------------------------------------------

type rowScanner interface {
	Scan(dest ...any) error
}

func scanQuest(s rowScanner) (Quest, error) {
	var (
		q                                Quest
		category, difficulty, recurrence string
	)
	err := s.Scan(&q.ID, &q.UserID, &q.Title, &q.Note, &category, &difficulty, &recurrence, &q.Active, &q.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Quest{}, ErrQuestNotFound
	}
	if err != nil {
		return Quest{}, fmt.Errorf("quest: scan quest: %w", err)
	}
	q.Category = Category(category)
	q.Difficulty = Difficulty(difficulty)
	q.Recurrence = Recurrence(recurrence)
	return q, nil
}

func scanLog(s rowScanner) (QuestLog, error) {
	var (
		l      QuestLog
		status string
	)
	err := s.Scan(&l.ID, &l.QuestID, &l.UserID, &l.Date, &status, &l.PointsAwarded, &l.CompletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return QuestLog{}, ErrNotCompleted
	}
	if err != nil {
		return QuestLog{}, fmt.Errorf("quest: scan log: %w", err)
	}
	l.Status = LogStatus(status)
	return l, nil
}

func rowsAffectedOrNotFound(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("quest: rows affected: %w", err)
	}
	if n == 0 {
		return ErrQuestNotFound
	}
	return nil
}
