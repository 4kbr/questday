package quest

import (
	"context"
	"time"
)

// Repository = kontrak penyimpanan module quest. Implementasi Postgres ada di
// repository_postgres.go — satu-satunya tempat SQL module ini. Setiap query
// difilter user_id (pertahanan kedua setelah cek kepemilikan di service).
type Repository interface {
	// Quest (definisi)
	CreateQuest(ctx context.Context, q Quest) error
	GetQuest(ctx context.Context, userID, id string) (Quest, error) // ErrQuestNotFound bila kosong
	ListQuestsByUser(ctx context.Context, userID string) ([]Quest, error)
	UpdateQuest(ctx context.Context, q Quest) error
	ArchiveQuest(ctx context.Context, userID, id string) error // soft delete: active=false

	// QuestLog (instance harian)
	CreateLog(ctx context.Context, l QuestLog) error // ErrAlreadyCompleted bila UNIQUE(quest_id,date) bentrok
	GetLog(ctx context.Context, userID, questID string, date time.Time) (QuestLog, error)
	DeleteLog(ctx context.Context, userID, questID string, date time.Time) error
	ListLogsByUserAndDate(ctx context.Context, userID string, date time.Time) ([]QuestLog, error)
	ListActiveDates(ctx context.Context, userID string, from, to time.Time) ([]time.Time, error)
}
