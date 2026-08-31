package scoring

// domain.go = entitas & aturan gamifikasi.
//
// TODO:
//   // Wallet = akumulasi poin/XP/level milik satu user.
//   type Wallet struct {
//       UserID      string
//       TotalPoints int
//       XP          int      // untuk MVP bisa == TotalPoints; pisahkan bila mau
//       Level       int
//   }
//
//   // Streak = rentetan hari aktif berturut-turut.
//   type Streak struct {
//       UserID     string
//       Current    int
//       Longest    int
//       LastActive time.Time   // tanggal lokal terakhir user menyelesaikan quest
//   }
//
//   // Transaction = ledger perubahan poin (audit + memungkinkan rollback).
//   type Transaction struct {
//       ID      string
//       UserID  string
//       QuestID string
//       Points  int          // + saat complete, - saat uncomplete
//       Date    time.Time
//       At      time.Time
//   }
//
//   // LevelForXP menghitung level dari XP. Simpan kurva di sini (mis. ambang
//   // 100/250/500/... atau rumus). Satu sumber kebenaran.
//   func LevelForXP(xp int) int
//
//   // NextStreak menghitung streak baru berdasarkan streak lama + jarak hari.
//   // Aturan (lihat DECISIONS): +1 kalau hari berurutan, tetap kalau hari sama,
//   // reset ke 1 kalau ada hari bolong. (Grace/freeze = fitur nanti.)
//   func NextStreak(cur Streak, activeDate time.Time) Streak
