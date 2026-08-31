# Phase 2 — Quest

**Tujuan:** module `quest` lengkap — definisi quest (CRUD), quest hari ini, dan
complete/uncomplete. Poin **belum** bertambah di phase ini: `ScoreAwarder`
disuntik sebagai no-op supaya quest bisa diuji sendirian.

**Prasyarat:** Phase 1 selesai (butuh middleware yang mengisi `userID` +
`timezone` ke context, dan tabel `quests` / `quest_logs`).

**Kenapa scoring ditunda:** GUIDES §8 langkah 5-6. Mengerjakan dua module
setengah jadi sekaligus bikin sulit tahu yang mana yang salah saat gagal.

**Ini module terberat di MVP** — 8 file, 10 method repository, 7 use case.

---

## T2.1 — `quest/domain.go`

- **Sentuh:** `internal/modules/quest/domain.go`
- **Isi:**
  ```go
  type Category   string  // olahraga, belajar, tidur, ngoding, ...
  type Difficulty string  // easy | medium | hard
  type Recurrence string  // daily (MVP)
  type LogStatus  string  // completed (MVP)

  type Quest struct {
      ID, UserID, Title, Note string
      Category   Category
      Difficulty Difficulty
      Recurrence Recurrence
      Active     bool
      CreatedAt  time.Time
  }
  type QuestLog struct {
      ID, QuestID, UserID string
      Date          time.Time  // tanggal LOKAL user (ADR-006)
      Status        LogStatus
      PointsAwarded int
      CompletedAt   time.Time
  }
  func (q Quest) Points() int
  var ErrQuestNotFound, ErrNotOwner, ErrAlreadyCompleted, ErrNotCompleted error
  ```
- **Aturan keras:** `Quest.Points()` adalah **satu-satunya sumber aturan poin**
  (ADR-007). Jangan ada angka poin di service, handler, SQL, atau scoring.
  Tabel poin per difficulty didefinisikan di sini sebagai konstanta.
- **Aturan:** domain tak import HTTP/SQL/package kita yang lain. Tambahkan juga
  validasi enum di sini (mis. `func (d Difficulty) Valid() bool`) supaya service
  tak menebak.
- **DoD:** 2 entitas, 4 enum, `Points()`, 4 error.
- **Verifikasi:** `make vet`; test di T2.10.

## T2.2 — `quest/dto.go`

- **Sentuh:** `internal/modules/quest/dto.go`
- **Isi:**
  ```go
  type CreateQuestRequest struct {
      Title      string `validate:"required"`
      Note       string
      Category   string `validate:"required"`
      Difficulty string `validate:"required,oneof=easy medium hard"`
  }
  type UpdateQuestRequest struct {   // partial update -> pointer semua
      Title, Note, Category, Difficulty *string
      Active *bool
  }
  type CompleteQuestRequest struct{} // MVP: kosong; field progress/value = v2
  type QuestResponse    struct{ ... }
  type QuestLogResponse struct{ ... }
  type TodayQuestsResponse struct {
      Date  string   // "2026-08-31", tanggal lokal user
      Items []struct{ Quest QuestResponse; Completed bool }
  }
  func toQuestResponse(q Quest) QuestResponse
  ```
- **Keputusan yang diselesaikan di sini:** scaffold menyebut
  `CompleteQuestRequest` "opsional progress/value" tanpa merinci. **MVP: body
  kosong** — poin sepenuhnya dari `Quest.Points()`. Kalau nanti butuh progress,
  itu keputusan baru → ADR baru.
- **Aturan:** `UpdateQuestRequest` pakai pointer supaya "tidak dikirim" bisa
  dibedakan dari "dikirim kosong".
- **DoD:** 7 tipe + mapper.
- **Verifikasi:** `make build`

## T2.3 — Verifikasi kolom migrasi cukup

- **Sentuh:** — (hanya cek; kalau kurang, **buat migrasi baru**, jangan edit
  `000001`)
- **Isi:** pastikan `quests` & `quest_logs` dari T0.3 menampung semua field di
  T2.1: `quests.note`, `quests.active`, `quest_logs.points_awarded`,
  `quest_logs.status`.
- **Aturan:** migrasi yang sudah diterapkan **tidak boleh diedit** (GUIDES §6) —
  kalau kurang, `make migrate-create name=alter_quests_...`.
- **DoD:** tiap field entitas punya kolomnya.
- **Verifikasi:** `psql -c '\d quests' -c '\d quest_logs'`

## T2.4 — `quest/repository.go` + `repository_postgres.go`

- **Sentuh:** `internal/modules/quest/repository.go`,
  `internal/modules/quest/repository_postgres.go`
- **Isi:** interface `Repository` **10 method**:

  | Kelompok | Method |
  |---|---|
  | Quest | `CreateQuest`, `GetQuest`, `ListQuestsByUser`, `UpdateQuest`, `ArchiveQuest` (soft delete: `active=false`) |
  | Log | `CreateLog`, `GetLog(questID, date)`, `DeleteLog`, `ListLogsByUserAndDate`, `ListActiveDates(userID, from, to)` |

  Impl: `type postgresRepository struct{ db *sql.DB }` + `newPostgresRepository`.
- **Aturan keras:**
  - SQL **hanya** di file ini.
  - **Setiap** query difilter `user_id` — termasuk `GetQuest` dan `DeleteLog`.
    Ini pertahanan kedua setelah cek kepemilikan di service.
  - `sql.ErrNoRows` → `ErrQuestNotFound`.
  - `CreateLog` mendeteksi pelanggaran `UNIQUE(quest_id, date)` → `ErrAlreadyCompleted`.
  - `ArchiveQuest` = soft delete, bukan `DELETE` — quest_logs lama harus tetap ada
    untuk perhitungan streak.
- **Catatan:** `ListActiveDates` dipakai untuk rekonstruksi streak nanti; walau
  belum terpakai di MVP, tetap implementasikan karena murah dan sudah dirancang.
- **DoD:** 10 method jalan, semua terikat `user_id`.
- **Verifikasi:** `make vet && make build`; cek manual
  `grep -c "user_id" repository_postgres.go` ≥ jumlah query.

## T2.5 — `quest/service.go` — port `ScoreAwarder` + use case

- **Sentuh:** `internal/modules/quest/service.go`
- **Isi:**
  ```go
  // Port keluar milik quest. Diimplementasi scoring, disuntik server. (ADR-005)
  type ScoreAwarder interface {
      OnQuestCompleted(ctx context.Context, userID, questID string, points int, date time.Time) error
      OnQuestUncompleted(ctx context.Context, userID, questID string, points int, date time.Time) error
  }

  type service struct{ repo Repository; scorer ScoreAwarder }
  func newService(repo Repository, scorer ScoreAwarder) *service
  ```
  7 use case: `CreateQuest`, `ListQuests`, `UpdateQuest`, `ArchiveQuest`,
  `GetToday(ctx, userID, localDate)`, `CompleteQuest(ctx, userID, questID, localDate)`,
  `UncompleteQuest(ctx, userID, questID, localDate)`.

  `CompleteQuest`: ambil quest (cek `q.UserID == userID` → `ErrNotOwner`; cek
  `q.Active`) → cek belum ada log tanggal itu → hitung `q.Points()` → `CreateLog`
  → `scorer.OnQuestCompleted`.
  `UncompleteQuest`: `GetLog` (tak ada → `ErrNotCompleted`) → `DeleteLog` →
  `scorer.OnQuestUncompleted`.
- **Aturan keras:**
  - **Interface ini milik `quest`**, bukan milik scoring. Ini yang menjaga arah
    dependensi (GUIDES §5). `quest` tak boleh `import .../modules/scoring`.
  - Service menerima `localDate` sebagai **argumen** — bukan menghitungnya
    sendiri, dan **dilarang memanggil `time.Now()`**. Yang menghitung adalah
    handler (T2.6).
  - Poin diambil dari `q.Points()`, tak pernah dihitung ulang di sini.
- **Catatan atomicity:** idealnya `CreateLog` + `OnQuestCompleted` satu transaksi
  DB. Di phase ini biarkan berurutan; keputusannya diambil di **T3.10**.
- **DoD:** port + 7 use case; nol import module lain.
- **Verifikasi:** `grep -n "modules/" internal/modules/quest/*.go` → hanya
  `modules/quest` sendiri.

## T2.6 — `quest/handler.go` — "hari ini" dari timezone user

- **Sentuh:** `internal/modules/quest/handler.go`
- **Isi:** `type handler struct{ svc *service }` + 7 handler yang cocok dengan
  routes T2.7. Pola tanggal lokal:
  ```go
  tz, _ := middleware.TimezoneFrom(r.Context())   // diisi Authenticator (ADR-013)
  loc, err := time.LoadLocation(tz)               // fallback "Asia/Jakarta" bila gagal
  today := time.Now().In(loc)
  localDate := time.Date(today.Year(), today.Month(), today.Day(), 0,0,0,0, time.UTC)
  ```
- **Aturan keras:** **dilarang** memakai `time.Now()` mentah untuk logika harian
  (ADR-006). Selalu lewat timezone user. User GMT+7 yang menyelesaikan quest jam
  23:00 harus tercatat di tanggal lokalnya, bukan tanggal UTC.
- **Aturan:** handler tipis; petakan error domain: `ErrQuestNotFound` → 404,
  `ErrNotOwner` → 404 (**bukan 403** — jangan bocorkan bahwa quest itu ada milik
  orang lain), `ErrAlreadyCompleted` → 409, `ErrNotCompleted` → 409.
- **Catatan:** pertimbangkan helper kecil `localDateFrom(r)` di file ini supaya
  7 handler tidak menyalin blok yang sama.
- **DoD:** 7 handler, nol `time.Now()` di luar helper tersebut.
- **Verifikasi:** `grep -n "time.Now()" internal/modules/quest/` → hanya di helper.

## T2.7 — `quest/routes.go`

- **Sentuh:** `internal/modules/quest/routes.go`
- **Isi:**
  ```
  GET    /quests
  POST   /quests
  GET    /quests/today
  PATCH  /quests/{questId}
  DELETE /quests/{questId}
  POST   /quests/{questId}/complete
  POST   /quests/{questId}/uncomplete
  ```
- **Perangkap chi:** `/quests/today` **harus** didaftarkan sebelum
  `/quests/{questId}`, kalau tidak `today` akan tertangkap sebagai `questId`.
  Urutan di scaffold sudah benar — pertahankan.
- **Aturan:** semua rute quest ada di grup **terproteksi** (butuh Authenticator).
- **DoD:** 7 rute; `GET /quests/today` tidak nyasar ke handler detail.
- **Verifikasi:** `curl localhost:8080/api/v1/quests/today -H "Authorization: Bearer $TOKEN"`

## T2.8 — `quest/module.go`

- **Sentuh:** `internal/modules/quest/module.go`
- **Isi:**
  ```go
  type Module struct{ handler *handler }
  func New(db *sql.DB, scorer ScoreAwarder) *Module
  func (m *Module) RegisterRoutes(r chi.Router)
  ```
- **Aturan:** `scorer` datang dari luar — module tak pernah membuat sendiri.
- **DoD:** `quest.New(db, scorer)` merakit repo→service→handler.
- **Verifikasi:** `make build`

## T2.9 — `server`: mount quest + ScoreAwarder no-op

- **Sentuh:** `internal/server/server.go`, `internal/server/router.go`
- **Isi:** tambahkan di composition root:
  ```go
  type noopScorer struct{}
  func (noopScorer) OnQuestCompleted(...) error   { return nil }
  func (noopScorer) OnQuestUncompleted(...) error { return nil }

  questMod := quest.New(db, noopScorer{})
  // di grup terproteksi: questMod.RegisterRoutes(r)
  ```
- **Aturan:** no-op ini **sementara**, diganti di T3.9. Taruh di `server` (bukan
  di `quest`) dan beri komentar `// TODO(T3.9): ganti dengan scoringMod.AsScoreAwarder()`.
- **DoD:** endpoint quest jalan tanpa scoring ada.
- **Verifikasi:** smoke test T2.10.

## T2.10 — Test phase 2

- **Sentuh (baru):** `internal/modules/quest/domain_test.go`,
  `internal/modules/quest/service_test.go`
- **Isi:**
  - `domain_test.go`: `Quest.Points()` untuk easy/medium/hard; enum `Valid()`
    menolak nilai ngawur.
  - `service_test.go` (fake `Repository` + fake `ScoreAwarder` yang mencatat
    panggilan):
    - `CompleteQuest` happy path → log dibuat, scorer dipanggil sekali dengan
      poin yang sama dengan `q.Points()`.
    - quest milik user lain → `ErrNotOwner`, scorer **tidak** dipanggil.
    - complete dua kali di tanggal sama → `ErrAlreadyCompleted`.
    - `UncompleteQuest` tanpa log → `ErrNotCompleted`.
  - **Uji batas hari:** panggil `CompleteQuest` dengan `localDate` dua timezone
    berbeda pada instant yang sama (mis. UTC 2026-08-31T17:00 → Asia/Jakarta
    sudah 09-01, UTC masih 08-31) dan pastikan log jatuh di tanggal yang benar.
- **Smoke test manual:**
  ```bash
  curl -sX POST localhost:8080/api/v1/quests -H "Authorization: Bearer $TOKEN" \
    -d '{"title":"Lari 5k","category":"olahraga","difficulty":"medium"}'
  curl -s localhost:8080/api/v1/quests/today -H "Authorization: Bearer $TOKEN"
  curl -sX POST localhost:8080/api/v1/quests/$ID/complete -H "Authorization: Bearer $TOKEN"
  curl -s localhost:8080/api/v1/quests/today -H "Authorization: Bearer $TOKEN"  # completed: true
  ```
- **DoD:** `make test` hijau; smoke test sesuai.

---

## Exit criteria Phase 2

- [ ] `make fmt && make vet && make build && make test` bersih.
- [ ] CRUD quest + `/quests/today` + complete/uncomplete jalan lewat HTTP.
- [ ] Complete dua kali di hari sama → 409; quest orang lain → 404.
- [ ] Poin **belum** bertambah (no-op) tapi `quest_logs` tercatat dengan `date`
      yang benar menurut timezone user.
- [ ] `internal/modules/quest` tidak meng-import `internal/modules/scoring`.
- [ ] Tak ada angka poin di luar `quest/domain.go`.
