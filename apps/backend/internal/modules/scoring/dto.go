package scoring

// dto.go = bentuk data HTTP module scoring. Acuan: contracts/openapi.yaml
// (semua response sukses dibungkus {"data": ...}, ADR-025).

const dateLayout = "2006-01-02"

// ScoreResponse = GET /me/score.
type ScoreResponse struct {
	TotalPoints       int `json:"total_points"`
	XP                int `json:"xp"`
	Level             int `json:"level"`
	PointsToNextLevel int `json:"points_to_next_level"`
}

// StreakResponse = GET /me/streak. LastActive null bila user belum pernah aktif.
type StreakResponse struct {
	Current    int     `json:"current"`
	Longest    int     `json:"longest"`
	LastActive *string `json:"last_active"`
}

// LeaderboardEntry = satu baris GET /leaderboard.
type LeaderboardEntry struct {
	Rank        int    `json:"rank"`
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Points      int    `json:"points"`
}

func toScoreResponse(w Wallet) ScoreResponse {
	return ScoreResponse{
		TotalPoints:       w.TotalPoints,
		XP:                w.XP,
		Level:             w.Level,
		PointsToNextLevel: PointsToNextLevel(w.XP),
	}
}

func toStreakResponse(s Streak) StreakResponse {
	out := StreakResponse{Current: s.Current, Longest: s.Longest}
	if !s.LastActive.IsZero() {
		d := s.LastActive.Format(dateLayout)
		out.LastActive = &d
	}
	return out
}
