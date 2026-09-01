// Package quest mengelola definisi quest (template task harian) dan log
// penyelesaiannya per hari.
//
// KEPUTUSAN KUNCI (ADR-004): Quest (definisi) dan QuestLog (instance harian)
// dipisah. Semua perhitungan "hari ini selesai apa" & streak berbasis QuestLog.
package quest

import (
	"errors"
	"time"
)

// --- Enum (tipe string biar jelas di DB & JSON) ---

type Category string   // olahraga, belajar, ngoding, tidur, ...
type Difficulty string // easy | medium | hard
type Recurrence string // daily (MVP)
type LogStatus string  // completed (MVP)

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"

	RecurrenceDaily Recurrence = "daily"

	LogStatusCompleted LogStatus = "completed"
)

// Valid melaporkan apakah d salah satu difficulty yang dikenal.
func (d Difficulty) Valid() bool {
	switch d {
	case DifficultyEasy, DifficultyMedium, DifficultyHard:
		return true
	default:
		return false
	}
}

// --- Aturan poin (SATU-SATUNYA sumber, ADR-007) ---
//
// Poin quest ditentukan difficulty-nya. Jangan menaruh angka poin di service,
// handler, SQL, atau module scoring — semuanya lewat Quest.Points().
const (
	pointsEasy   = 5
	pointsMedium = 10
	pointsHard   = 20
)

// --- Entitas ---

// Quest adalah definisi/template task harian milik seorang user.
type Quest struct {
	ID         string
	UserID     string
	Title      string
	Note       string
	Category   Category
	Difficulty Difficulty
	Recurrence Recurrence
	Active     bool // false = diarsipkan (soft delete), bukan dihapus keras
	CreatedAt  time.Time
}

// Points memetakan Difficulty -> poin. Difficulty tak dikenal -> 0.
func (q Quest) Points() int {
	switch q.Difficulty {
	case DifficultyEasy:
		return pointsEasy
	case DifficultyMedium:
		return pointsMedium
	case DifficultyHard:
		return pointsHard
	default:
		return 0
	}
}

// QuestLog adalah penyelesaian satu Quest pada satu tanggal lokal user.
type QuestLog struct {
	ID            string
	QuestID       string
	UserID        string
	Date          time.Time // tanggal LOKAL user (ADR-006), bukan instant UTC
	Status        LogStatus
	PointsAwarded int
	CompletedAt   time.Time
}

// --- Error domain (dipetakan ke HTTP lewat httpx.RegisterErrorMapping, ADR-023) ---
var (
	ErrQuestNotFound    = errors.New("quest tidak ditemukan")
	ErrNotOwner         = errors.New("quest bukan milik user")
	ErrAlreadyCompleted = errors.New("quest sudah selesai hari ini")
	ErrNotCompleted     = errors.New("quest belum selesai hari ini")
)
