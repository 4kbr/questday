package quest

// repository.go = KONTRAK penyimpanan (interface). Service bergantung ke
// interface ini, bukan ke Postgres langsung — jadi mudah diganti/di-mock.
// Implementasi konkret ada di repository_postgres.go.
//
// TODO:
//   type Repository interface {
//       // Quest (definisi)
//       CreateQuest(ctx, *Quest) error
//       GetQuest(ctx, id string) (*Quest, error)          // ErrQuestNotFound bila kosong
//       ListQuestsByUser(ctx, userID string) ([]Quest, error)
//       UpdateQuest(ctx, *Quest) error
//       ArchiveQuest(ctx, id string) error                 // soft delete
//
//       // QuestLog (instance harian)
//       CreateLog(ctx, *QuestLog) error
//       GetLog(ctx, questID, userID string, date time.Time) (*QuestLog, error)
//       DeleteLog(ctx, questID, userID string, date time.Time) error
//       ListLogsByUserAndDate(ctx, userID string, date time.Time) ([]QuestLog, error)
//       // untuk streak: tanggal-tanggal user pernah menyelesaikan sesuatu
//       ListActiveDates(ctx, userID string, from, to time.Time) ([]time.Time, error)
//   }
