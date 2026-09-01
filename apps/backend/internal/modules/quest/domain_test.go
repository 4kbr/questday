package quest

import "testing"

func TestQuestPoints(t *testing.T) {
	cases := []struct {
		diff Difficulty
		want int
	}{
		{DifficultyEasy, 5},
		{DifficultyMedium, 10},
		{DifficultyHard, 20},
		{Difficulty("ngawur"), 0},
		{Difficulty(""), 0},
	}
	for _, c := range cases {
		got := Quest{Difficulty: c.diff}.Points()
		if got != c.want {
			t.Errorf("Points(%q) = %d, mau %d", c.diff, got, c.want)
		}
	}
}

func TestDifficultyValid(t *testing.T) {
	valid := []Difficulty{DifficultyEasy, DifficultyMedium, DifficultyHard}
	for _, d := range valid {
		if !d.Valid() {
			t.Errorf("%q harus valid", d)
		}
	}
	invalid := []Difficulty{"", "EASY", "sedang", "hardcore"}
	for _, d := range invalid {
		if d.Valid() {
			t.Errorf("%q harus tidak valid", d)
		}
	}
}
