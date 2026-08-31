package quest

// service.go = use case / logika aplikasi module quest. Mengorkestrasi
// repository + aturan domain. TIDAK tahu soal HTTP maupun SQL.
//
// GAMIFIKASI: saat quest diselesaikan, poin harus bertambah. Supaya module
// quest tidak bergantung ke module scoring, kita definisikan PORT di sini dan
// scoring yang mengimplementasikannya (disuntik dari server). Lihat DECISIONS
// "Orkestrasi vs event bus".
//
// TODO:
//   // ScoreAwarder = port keluar. scoring.Service akan memenuhinya.
//   type ScoreAwarder interface {
//       OnQuestCompleted(ctx, userID, questID string, points int, date time.Time) error
//       OnQuestUncompleted(ctx, userID, questID string, points int, date time.Time) error
//   }
//
//   type service struct { repo Repository; scorer ScoreAwarder }
//   func newService(repo Repository, scorer ScoreAwarder) *service
//
//   // Use case:
//   func (s *service) CreateQuest(ctx, userID string, in CreateQuestRequest) (*Quest, error)
//   func (s *service) ListQuests(ctx, userID string) ([]Quest, error)
//   func (s *service) UpdateQuest(ctx, userID, id string, in UpdateQuestRequest) (*Quest, error)
//   func (s *service) ArchiveQuest(ctx, userID, id string) error
//
//   // GetToday: gabungkan quest aktif user dengan log hari ini -> selesai/belum.
//   func (s *service) GetToday(ctx, userID string, today time.Time) (TodayQuestsResponse, error)
//
//   // CompleteQuest: cek kepemilikan & belum selesai -> buat log -> panggil
//   // scorer.OnQuestCompleted. Pertimbangkan transaksi agar log & poin atomik.
//   func (s *service) CompleteQuest(ctx, userID, questID string, date time.Time) error
//
//   // UncompleteQuest: hapus log hari ini -> scorer.OnQuestUncompleted (rollback).
//   func (s *service) UncompleteQuest(ctx, userID, questID string, date time.Time) error
