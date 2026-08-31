package quest

// domain.go = entitas inti, enum, aturan invarian, dan error domain module ini.
// Tidak ada dependensi ke HTTP maupun SQL di sini — Go murni.
//
// TODO:
//   // --- Entitas ---
//   type Quest struct {
//       ID         string
//       UserID     string
//       Title      string
//       Note       string
//       Category   Category        // health, learning, coding, sleep, ...
//       Difficulty Difficulty       // menentukan poin (lihat Points())
//       Recurrence Recurrence       // daily untuk MVP
//       Active     bool             // false = diarsipkan, bukan dihapus keras
//       CreatedAt  time.Time
//   }
//
//   type QuestLog struct {
//       ID            string
//       QuestID       string
//       UserID        string
//       Date          time.Time    // tanggal lokal user (batas hari, lihat DECISIONS)
//       Status        LogStatus    // completed | (nanti: partial)
//       PointsAwarded int          // dicatat saat selesai (audit & rollback)
//       CompletedAt   time.Time
//   }
//
//   // --- Enum (pakai tipe string biar jelas di DB & JSON) ---
//   type Category string
//   type Difficulty string   // easy | medium | hard
//   type Recurrence string   // daily (MVP); weekly/custom = nanti
//   type LogStatus string    // completed (MVP); partial = nanti
//
//   // Points memetakan Difficulty -> poin. Simpan aturan poin DI SINI supaya
//   // satu sumber kebenaran (jangan sebar angka poin di banyak tempat).
//   func (q Quest) Points() int
//
//   // --- Error domain (dipetakan ke HTTP di handler/httpx) ---
//   var (
//       ErrQuestNotFound     = errors.New("quest not found")
//       ErrNotOwner          = errors.New("quest bukan milik user")
//       ErrAlreadyCompleted  = errors.New("quest sudah selesai hari ini")
//       ErrNotCompleted      = errors.New("quest belum selesai hari ini")
//   )
