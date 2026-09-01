package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

// pgUniqueViolation adalah SQLSTATE Postgres untuk pelanggaran UNIQUE.
const pgUniqueViolation = "23505"

// postgresRepository mengimplementasi Repository di atas *sql.DB (ADR-012).
type postgresRepository struct {
	db *sql.DB
}

func newPostgresRepository(db *sql.DB) *postgresRepository {
	return &postgresRepository{db: db}
}

// Create menyimpan user baru. ID sudah digenerate pemanggil (UUIDv7, ADR-015).
// Bentrok UNIQUE(email) dipetakan ke ErrEmailTaken.
func (r *postgresRepository) Create(ctx context.Context, u User) error {
	const q = `
		INSERT INTO users (id, email, password_hash, display_name, timezone)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, q, u.ID, u.Email, u.PasswordHash, u.DisplayName, u.Timezone)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return ErrEmailTaken
		}
		return fmt.Errorf("user: create: %w", err)
	}
	return nil
}

// GetByEmail mengambil user berdasarkan email. Kosong -> ErrUserNotFound.
func (r *postgresRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	const q = `
		SELECT id, email, password_hash, display_name, timezone, created_at
		FROM users WHERE email = $1`
	return r.scanOne(r.db.QueryRowContext(ctx, q, email))
}

// GetByID mengambil user berdasarkan id. Kosong -> ErrUserNotFound.
func (r *postgresRepository) GetByID(ctx context.Context, id string) (User, error) {
	const q = `
		SELECT id, email, password_hash, display_name, timezone, created_at
		FROM users WHERE id = $1`
	return r.scanOne(r.db.QueryRowContext(ctx, q, id))
}

// Update menyimpan perubahan display_name & timezone. Email tak ikut diubah.
func (r *postgresRepository) Update(ctx context.Context, u User) error {
	const q = `
		UPDATE users SET display_name = $2, timezone = $3
		WHERE id = $1`

	res, err := r.db.ExecContext(ctx, q, u.ID, u.DisplayName, u.Timezone)
	if err != nil {
		return fmt.Errorf("user: update: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("user: update rows: %w", err)
	}
	if n == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *postgresRepository) scanOne(row *sql.Row) (User, error) {
	var u User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.Timezone, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("user: scan: %w", err)
	}
	return u, nil
}
