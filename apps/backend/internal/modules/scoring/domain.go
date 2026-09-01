// Package scoring memegang SEMUA logika gamifikasi: poin, XP, level, streak,
// dan leaderboard. Bereaksi terhadap penyelesaian quest lewat port
// quest.ScoreAwarder (diimplementasi service, disuntik server).
package scoring

import "time"

// Wallet = akumulasi poin/XP/level milik satu user. Untuk MVP XP == TotalPoints
// (ADR-007); dipisah supaya kurva XP bisa berubah tanpa mengutak-atik poin.
type Wallet struct {
	UserID      string
	TotalPoints int
	XP          int
	Level       int
}

// Streak = rentetan hari aktif berturut-turut (tanggal lokal user).
type Streak struct {
	UserID     string
	Current    int
	Longest    int
	LastActive time.Time // zero = belum pernah aktif
}

// Transaction = ledger perubahan poin (audit + memungkinkan rollback).
type Transaction struct {
	ID      string
	UserID  string
	QuestID string
	Points  int       // + saat complete, - saat uncomplete
	Date    time.Time // tanggal lokal user
	At      time.Time
}

// --- Kurva level (SATU-SATUNYA sumber, ADR-007) ---
//
// Ambang XP kumulatif: 0 / 100 / 250 / 450 / 700 / 1000 / ...
// Kenaikan +100, +150, +200, +250, ... (selisih +50 tiap naik level).
// XPForLevel(n) = 25*n*(n+1) - 50 untuk n >= 1.

// XPForLevel mengembalikan XP minimum untuk mencapai level n (n >= 1).
func XPForLevel(n int) int {
	if n < 1 {
		n = 1
	}
	return 25*n*(n+1) - 50
}

// LevelForXP mengembalikan level tertinggi yang dicapai dengan xp poin.
// Tanpa cap untuk MVP.
func LevelForXP(xp int) int {
	if xp < 0 {
		xp = 0
	}
	level := 1
	for XPForLevel(level+1) <= xp {
		level++
	}
	return level
}

// PointsToNextLevel = sisa XP menuju level berikutnya.
func PointsToNextLevel(xp int) int {
	next := XPForLevel(LevelForXP(xp) + 1)
	return next - xp
}

// --- Aturan streak (SATU-SATUNYA sumber, ADR-008) ---

// NextStreak menghitung streak baru dari streak lama + tanggal aktif.
// Murni: tak menyentuh DB, tak memanggil time.Now(). activeDate = tanggal lokal
// user (UTC wall-clock, jam 00:00), konsisten dengan kolom DATE.
//
//	activeDate == LastActive        -> tidak berubah
//	activeDate == LastActive + 1hr  -> Current + 1
//	selisih > 1 hari / belum pernah -> Current = 1
//
// Longest tak pernah turun (ADR-008). "Freeze"/grace day = fitur nanti.
func NextStreak(cur Streak, activeDate time.Time) Streak {
	activeDate = dateOnly(activeDate)
	next := cur
	next.LastActive = activeDate

	switch {
	case cur.LastActive.IsZero():
		next.Current = 1
	case activeDate.Equal(dateOnly(cur.LastActive)):
		return cur // hari yang sama: tak ada yang berubah
	case activeDate.Equal(dateOnly(cur.LastActive).AddDate(0, 0, 1)):
		next.Current = cur.Current + 1
	default:
		next.Current = 1
	}

	if next.Current > next.Longest {
		next.Longest = next.Current
	}
	return next
}

func dateOnly(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
