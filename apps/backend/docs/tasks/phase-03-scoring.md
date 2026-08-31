# Phase 3 — Scoring

**Tujuan:** module `scoring` lengkap (poin, XP, level, streak, leaderboard) dan
menyambungkannya ke `quest` lewat port `ScoreAwarder`. Setelah phase ini
gamifikasi benar-benar hidup: menyelesaikan quest menaikkan poin & streak.

**Prasyarat:** Phase 2 selesai (port `quest.ScoreAwarder` sudah terdefinisi,
tabel `wallets` / `streaks` / `point_transactions` sudah ada).

**Rujukan:** ADR-005 (port), ADR-007 (poin & level), ADR-008 (streak),
ADR-009 (uncomplete), ADR-014 (leaderboard lewat `UserDirectory`).

---

## T3.1 — `scoring/domain.go`

- **Sentuh:** `internal/modules/scoring/domain.go`
- **Isi:**
  ```go
  type Wallet struct{ UserID string; TotalPoints, XP, Level int }
  type Streak struct{ UserID string; Current, Longest int; LastActive time.Time }
  type Transaction struct {
      ID, UserID, QuestID string
      Points int          // boleh negatif
      Date   time.Time    // tanggal lokal user
      At     time.Time
  }

  func LevelForXP(xp int) int
  func NextStreak(cur Streak, activeDate time.Time) Streak
  ```
- **Aturan keras:**
  - `LevelForXP` adalah **satu-satunya** definisi kurva level (ADR-007). Jangan
    ada rumus level di service/handler/SQL.
  - `NextStreak` adalah **satu-satunya** aturan streak (ADR-008):
    | Kondisi | Hasil |
    |---|---|
    | `activeDate == LastActive` | tidak berubah |
    | `activeDate == LastActive + 1 hari` | `Current + 1` |
    | selisih > 1 hari, atau belum pernah aktif | `Current = 1` |
    `Longest = max(Longest, Current)` — **tak pernah turun**, supaya pencapaian
    tak hilang saat reset.
  - Fungsi ini **murni**: tak menyentuh DB, tak memanggil `time.Now()`.
    `activeDate` datang dari luar (tanggal lokal user).
- **DoD:** 3 entitas + 2 fungsi murni.
- **Verifikasi:** test T3.11.

## T3.2 — `scoring/dto.go`

- **Sentuh:** `internal/modules/scoring/dto.go`
- **Isi:**
  ```go
  type ScoreResponse struct {
      TotalPoints, XP, Level, PointsToNextLevel int
  }
  type StreakResponse struct {
      Current, Longest int
      LastActive       string   // "2026-08-31"
  }
  type LeaderboardEntry struct {
      Rank   int
      UserID string
      DisplayName string
      Points int
  }
  ```
  Plus mapper dari entitas domain.
- **Catatan:** `PointsToNextLevel` diturunkan dari `LevelForXP` — jangan
  hardcode ambangnya di sini.
- **DoD:** 3 DTO + mapper.
- **Verifikasi:** `make build`

## T3.3 — `scoring/repository.go` + `repository_postgres.go`

- **Sentuh:** `internal/modules/scoring/repository.go`,
  `internal/modules/scoring/repository_postgres.go`
- **Isi:** interface `Repository` 6 method:
  ```go
  GetWallet(ctx, userID) (Wallet, error)     // buat default {0,0,1} kalau belum ada
  SaveWallet(ctx, w Wallet) error
  AddTransaction(ctx, t Transaction) error
  GetStreak(ctx, userID) (Streak, error)     // default kosong kalau belum ada
  SaveStreak(ctx, s Streak) error
  Leaderboard(ctx, limit int) ([]Wallet, error)   // userID + points saja
  ```
- **Aturan keras:**
  - `Leaderboard` **tidak boleh JOIN ke tabel `users`** (ADR-014). Nama diisi
    di service lewat port `UserDirectory` (T3.4). Scoring tak boleh tahu skema
    module lain — walaupun satu database.
  - `SaveWallet` & `SaveStreak` pakai **UPSERT** (`INSERT ... ON CONFLICT
    (user_id) DO UPDATE`) supaya user baru tak perlu inisialisasi terpisah.
  - SQL hanya di file ini; semua query terikat `user_id` kecuali `Leaderboard`.
- **DoD:** 6 method; wallet & streak idempoten untuk user baru.
- **Verifikasi:** `make vet && make build`

## T3.4 — Port `UserDirectory` + `user.AsUserDirectory()`

- **Sentuh:** `internal/modules/scoring/service.go` (definisi port),
  `internal/modules/user/module.go` + `service.go` (implementasi),
  `internal/server/server.go` (penjahitan)
- **Isi:**
  ```go
  // scoring/service.go — port MILIK scoring
  type UserDirectory interface {
      NamesByIDs(ctx context.Context, ids []string) (map[string]string, error)
  }

  // user/service.go
  func (s *service) NamesByIDs(ctx, ids []string) (map[string]string, error)
  // user/module.go
  func (m *Module) AsUserDirectory() *service

  // server/server.go
  scoringMod := scoring.New(db, userMod.AsUserDirectory())
  ```
  `NamesByIDs` butuh method repository baru di `user`:
  `ListNamesByIDs(ctx, ids) (map[string]string, error)` (satu query
  `WHERE id = ANY($1)`).
- **Aturan keras (GUIDES §5):** **interface milik peminta**, di sini `scoring`.
  `user` hanya kebetulan cocok. `scoring` tak boleh `import .../modules/user`,
  dan `user` tak boleh `import .../modules/scoring`.
- **Kenapa begini:** `LeaderboardEntry.DisplayName` berasal dari tabel `users`.
  Jalan pintasnya JOIN, tapi itu membuat SQL scoring bergantung pada skema user —
  pelanggaran batas module. ADR-014 memilih port; biayanya satu query tambahan.
- **DoD:** leaderboard menampilkan nama tanpa satu pun JOIN ke `users` di SQL scoring.
- **Verifikasi:** `grep -n "users" internal/modules/scoring/repository_postgres.go` → kosong.

## T3.5 — `scoring/service.go`: implementasi port `ScoreAwarder`

- **Sentuh:** `internal/modules/scoring/service.go`
- **Isi:**
  ```go
  type service struct{ repo Repository; dir UserDirectory }
  func newService(repo Repository, dir UserDirectory) *service

  func (s *service) OnQuestCompleted(ctx, userID, questID string, points int, date time.Time) error
  func (s *service) OnQuestUncompleted(ctx, userID, questID string, points int, date time.Time) error
  ```
  `OnQuestCompleted`: `GetWallet` → `TotalPoints += points`, `XP += points`,
  `Level = LevelForXP(XP)` → `SaveWallet` → `GetStreak` →
  `NextStreak(cur, date)` → `SaveStreak` → `AddTransaction(+points)`.
  `OnQuestUncompleted`: kebalikan poin & XP, level dihitung ulang,
  `AddTransaction(-points)`.
- **Aturan keras:**
  - Signature **wajib** cocok dengan `quest.ScoreAwarder` (T2.5) — tapi `scoring`
    **tidak** mengimpor `quest` untuk memastikan itu. Kecocokan diverifikasi saat
    `server` menyuntik (compile-time).
  - `OnQuestUncompleted` **tidak memutar balik streak** (ADR-009) — memutar balik
    dengan benar butuh rekonstruksi urutan hari. Beri komentar yang menunjuk
    ADR-009 supaya tidak terlihat seperti bug.
  - `Level` selalu lewat `LevelForXP`, tak pernah `level++`.
  - `TotalPoints` tak boleh negatif — clamp di 0 dan pertimbangkan mencatatnya.
- **DoD:** 2 method; wallet, streak, dan transaksi konsisten.
- **Verifikasi:** test T3.11.

## T3.6 — `scoring/service.go`: query

- **Sentuh:** `internal/modules/scoring/service.go`
- **Isi:**
  ```go
  func (s *service) GetScore(ctx, userID) (ScoreResponse, error)
  func (s *service) GetStreak(ctx, userID) (StreakResponse, error)
  func (s *service) Leaderboard(ctx, limit int) ([]LeaderboardEntry, error)
  ```
  `Leaderboard`: `repo.Leaderboard(limit)` → kumpulkan userID → `dir.NamesByIDs`
  → gabung + isi `Rank` (1-based, urut poin turun).
- **Aturan:** kalau `NamesByIDs` tak menemukan sebuah ID (user terhapus), jangan
  error — isi nama fallback dan tetap tampilkan.
- **DoD:** 3 query; leaderboard berisi nama & rank.
- **Verifikasi:** `curl .../leaderboard`

## T3.7 — `scoring/handler.go` + `routes.go`

- **Sentuh:** `internal/modules/scoring/handler.go`,
  `internal/modules/scoring/routes.go`
- **Isi:** handler `score`, `streak`, `leaderboard`; rute:
  ```
  GET /me/score
  GET /me/streak
  GET /leaderboard        # query param ?limit= (default 20, clamp maks 100)
  ```
- **Keputusan yang harus diambil:** scaffold menaruh ketiganya di grup
  terproteksi. `/me/*` jelas butuh login. **`/leaderboard`: tetap di grup
  terproteksi untuk MVP** (papan peringkat internal, bukan halaman publik). Kalau
  nanti dibuka publik → ADR baru.
- **Aturan:** handler tipis; `limit` divalidasi & di-clamp di handler.
- **DoD:** 3 endpoint jalan.
- **Verifikasi:** `make build` + curl.

## T3.8 — `scoring/module.go`

- **Sentuh:** `internal/modules/scoring/module.go`
- **Isi:**
  ```go
  type Module struct{ svc *service; handler *handler }
  func New(db *sql.DB, dir UserDirectory) *Module
  func (m *Module) AsScoreAwarder() *service   // dikonsumsi quest lewat port-nya
  func (m *Module) RegisterRoutes(r chi.Router)
  ```
- **Aturan:** `AsScoreAwarder` mengembalikan `*service` (unexported type) —
  itu wajar dan disengaja: pemanggil hanya memakainya sebagai
  `quest.ScoreAwarder`.
- **DoD:** `AsScoreAwarder()` memenuhi `quest.ScoreAwarder` saat dikompilasi.
- **Verifikasi:** `make build` (kegagalan di sini = signature port tak cocok).

## T3.9 — `server`: ganti no-op jadi scoring sungguhan

- **Sentuh:** `internal/server/server.go`, `internal/server/router.go`
- **Isi:**
  ```go
  userMod    := user.New(db, jwt)
  scoringMod := scoring.New(db, userMod.AsUserDirectory())
  questMod   := quest.New(db, scoringMod.AsScoreAwarder())   // port!
  ```
  Hapus `noopScorer` dari T2.9 dan komentar `TODO(T3.9)`-nya. Mount
  `scoringMod.RegisterRoutes(r)` di grup terproteksi.
- **Aturan:** urutan instansiasi penting — `user` dulu (dibutuhkan scoring),
  lalu `scoring`, lalu `quest`.
- **DoD:** `noopScorer` sudah tak ada di codebase.
- **Verifikasi:** `grep -rn "noopScorer" internal/` → kosong.

## T3.10 — Atomicity: log + poin dalam satu transaksi

- **Sentuh:** `internal/modules/quest/service.go`,
  `internal/modules/quest/repository.go`, `internal/modules/scoring/*`
  (atau **hanya** `docs/DECISIONS.md` kalau ditunda)
- **Masalah:** `CompleteQuest` melakukan dua tulisan ke DB — `CreateLog` (tabel
  quest) lalu `OnQuestCompleted` (tabel scoring). Kalau yang kedua gagal, log
  sudah terlanjur ada dan poin tak pernah masuk. ARCHITECTURE menyebut ini
  ("idealnya satu transaksi DB").
- **Dua pilihan:**
  1. **Kerjakan sekarang:** alirkan `*sql.Tx` lewat context atau lewat argumen
     port (mis. `ScoreAwarder.OnQuestCompleted(ctx, tx, ...)`) — tapi ini
     membocorkan `*sql.Tx` ke dalam kontrak port, jadi port tak lagi murni.
     Alternatif lebih bersih: helper transaksi di `platform/database`
     (`WithTx(ctx, db, fn)`) + repository yang bisa menerima `sql.DB`/`sql.Tx`
     lewat interface `execer` bersama.
  2. **Tunda:** biarkan berurutan; kalau `OnQuestCompleted` gagal, service
     mengembalikan error **dan** menghapus log yang baru dibuat (kompensasi
     manual).
- **Aturan:** apa pun yang dipilih, **wajib ditulis sebagai ADR baru**
  (ADR-016) dengan konsekuensinya. Jangan diam-diam dibiarkan.
- **DoD:** perilaku saat scorer gagal terdefinisi & teruji, dan tercatat di
  `docs/DECISIONS.md`.
- **Verifikasi:** test dengan fake scorer yang selalu error — pastikan tidak ada
  log yatim yang tertinggal.

## T3.11 — Test phase 3

- **Sentuh (baru):** `internal/modules/scoring/domain_test.go`,
  `internal/modules/scoring/service_test.go`
- **Isi:**
  - `domain_test.go`:
    - `LevelForXP` di **titik batas** tiap level (xp tepat di ambang, tepat di
      bawah, 0, dan nilai besar).
    - `NextStreak` 4 kasus: pertama kali (LastActive nol) → 1; hari sama → tetap;
      hari berikutnya → +1; bolong ≥2 hari → reset ke 1. Plus: `Longest` tak
      pernah turun.
  - `service_test.go` (fake `Repository` + fake `UserDirectory`):
    - `OnQuestCompleted` → poin, XP, level, streak, transaksi (+) benar.
    - `OnQuestCompleted` lalu `OnQuestUncompleted` → poin & XP kembali ke semula,
      ada transaksi (−), **streak tidak berubah** (mengunci perilaku ADR-009).
    - `Leaderboard` mengisi nama & rank; ID yang tak dikenal tetap tampil.
- **Smoke test manual:**
  ```bash
  curl -sX POST localhost:8080/api/v1/quests/$ID/complete -H "Authorization: Bearer $TOKEN"
  curl -s localhost:8080/api/v1/me/score  -H "Authorization: Bearer $TOKEN"   # poin naik
  curl -s localhost:8080/api/v1/me/streak -H "Authorization: Bearer $TOKEN"   # current: 1
  curl -s localhost:8080/api/v1/leaderboard -H "Authorization: Bearer $TOKEN"
  curl -sX POST localhost:8080/api/v1/quests/$ID/uncomplete -H "Authorization: Bearer $TOKEN"
  curl -s localhost:8080/api/v1/me/score -H "Authorization: Bearer $TOKEN"    # poin balik
  ```
- **DoD:** `make test` hijau; smoke test sesuai.

---

## Exit criteria Phase 3

- [ ] `make fmt && make vet && make build && make test` bersih.
- [ ] Complete quest menaikkan poin, XP, level, dan streak; transaksi tercatat.
- [ ] Uncomplete mengembalikan poin & mencatat transaksi negatif; streak tetap
      (sesuai ADR-009, bukan bug).
- [ ] `/me/score`, `/me/streak`, `/leaderboard` benar.
- [ ] `scoring` tidak meng-import `user` maupun `quest`; SQL scoring tidak
      menyentuh tabel `users`.
- [ ] `noopScorer` sudah dihapus.
- [ ] Keputusan atomicity (T3.10) tercatat sebagai ADR-016.
