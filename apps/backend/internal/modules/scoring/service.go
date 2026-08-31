package scoring

// service.go = logika gamifikasi + implementasi port quest.ScoreAwarder.
//
// TODO:
//   type service struct { repo Repository }
//   func newService(repo Repository) *service
//
//   // --- Memenuhi quest.ScoreAwarder ---
//   // OnQuestCompleted: tambah poin ke wallet, hitung ulang level, update streak
//   // (NextStreak), catat Transaction (+points). Idealnya dalam satu transaksi DB.
//   func (s *service) OnQuestCompleted(ctx, userID, questID string, points int, date time.Time) error
//
//   // OnQuestUncompleted: kurangi poin, hitung ulang level, catat Transaction
//   // (-points). Catatan: memutar-balik streak itu rumit — untuk MVP boleh
//   // hanya sesuaikan poin dan biarkan streak apa adanya (catat di DECISIONS).
//   func (s *service) OnQuestUncompleted(ctx, userID, questID string, points int, date time.Time) error
//
//   // --- Query untuk endpoint ---
//   func (s *service) GetScore(ctx, userID string) (ScoreResponse, error)
//   func (s *service) GetStreak(ctx, userID string) (StreakResponse, error)
//   func (s *service) Leaderboard(ctx, limit int) ([]LeaderboardEntry, error)
