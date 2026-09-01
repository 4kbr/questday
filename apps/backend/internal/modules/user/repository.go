package user

import "context"

// Repository adalah kontrak penyimpanan user. Implementasi Postgres ada di
// repository_postgres.go — satu-satunya tempat SQL module ini.
type Repository interface {
	Create(ctx context.Context, u User) error
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, u User) error
}
