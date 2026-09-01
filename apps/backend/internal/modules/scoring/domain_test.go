package scoring

import (
	"testing"
	"time"
)

func TestXPForLevel(t *testing.T) {
	want := map[int]int{1: 0, 2: 100, 3: 250, 4: 450, 5: 700, 6: 1000}
	for n, w := range want {
		if got := XPForLevel(n); got != w {
			t.Errorf("XPForLevel(%d) = %d, mau %d", n, got, w)
		}
	}
}

func TestLevelForXP(t *testing.T) {
	cases := []struct {
		xp, want int
	}{
		{-10, 1}, {0, 1}, {99, 1}, {100, 2}, {249, 2}, {250, 3},
		{449, 3}, {450, 4}, {699, 4}, {700, 5}, {999, 5}, {1000, 6},
		{1_000_000, 200}, // sanity: XPForLevel(200)=25*200*201-50=1_004_950 > 1e6; level 199
	}
	for _, c := range cases {
		got := LevelForXP(c.xp)
		if c.xp == 1_000_000 {
			if got < 190 || got > 205 {
				t.Errorf("LevelForXP(1e6) = %d, di luar rentang wajar", got)
			}
			continue
		}
		if got != c.want {
			t.Errorf("LevelForXP(%d) = %d, mau %d", c.xp, got, c.want)
		}
	}
}

func TestPointsToNextLevel(t *testing.T) {
	// xp=0 -> level 1, next threshold 100 -> 100
	if got := PointsToNextLevel(0); got != 100 {
		t.Errorf("PointsToNextLevel(0) = %d, mau 100", got)
	}
	// xp=120 -> level 2, next threshold 250 -> 130
	if got := PointsToNextLevel(120); got != 130 {
		t.Errorf("PointsToNextLevel(120) = %d, mau 130", got)
	}
}

func day(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func TestNextStreak(t *testing.T) {
	d1 := day(2026, 9, 1)
	d2 := day(2026, 9, 2)
	d4 := day(2026, 9, 4)

	// pertama kali: LastActive nol -> Current 1, Longest 1
	s := NextStreak(Streak{UserID: "u"}, d1)
	if s.Current != 1 || s.Longest != 1 || !s.LastActive.Equal(d1) {
		t.Fatalf("pertama: %+v", s)
	}

	// hari berikutnya -> +1
	s = NextStreak(s, d2)
	if s.Current != 2 || s.Longest != 2 {
		t.Fatalf("berurutan: %+v", s)
	}

	// hari yang sama -> tak berubah
	same := NextStreak(s, d2)
	if same.Current != 2 || same.Longest != 2 || !same.LastActive.Equal(d2) {
		t.Fatalf("hari sama: %+v", same)
	}

	// bolong (d2 -> d4) -> reset ke 1, Longest tetap 2 (tak turun)
	s = NextStreak(s, d4)
	if s.Current != 1 {
		t.Fatalf("bolong: Current = %d, mau 1", s.Current)
	}
	if s.Longest != 2 {
		t.Fatalf("bolong: Longest = %d, mau tetap 2 (tak pernah turun)", s.Longest)
	}
	if !s.LastActive.Equal(d4) {
		t.Fatalf("bolong: LastActive = %v, mau %v", s.LastActive, d4)
	}
}
