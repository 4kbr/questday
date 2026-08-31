# Phase 3 — Scoring

**Tujuan:** gamifikasi terlihat — poin, XP, level, streak, dan leaderboard.
Setelah phase ini mencentang quest langsung terasa "berbuah": angka di dashboard
ikut naik tanpa reload.

**Prasyarat:** Phase 2 selesai (butuh `questKeys` — F3.6 bergantung penuh
padanya). Backend Phase 3 jalan — atau `VITE_USE_MOCK=true`.

**Endpoint yang dipakai:** `GET /me/score`, `GET /me/streak`, `GET /leaderboard`.

---

## F3.1 — `apis/scoring.api.ts`

- **Sentuh (baru):** `src/apis/scoring.api.ts`
- **Isi:**
  ```ts
  export const scoringApi = {
    score:       ()             => api.get<Score>('/me/score'),
    streak:      ()             => api.get<Streak>('/me/streak'),
    leaderboard: (limit = 20)   => api.get<LeaderboardEntry[]>('/leaderboard', { params: { limit } }),
  }
  ```
- **Aturan:** murni HTTP. `limit` punya default yang sama dengan default backend
  (20) supaya perilakunya tak berbeda antara dikirim & tidak.
- **DoD:** 3 fungsi terketik.
- **Verifikasi:** `npm run typecheck`

## F3.2 — `scoringKeys` + queries

- **Sentuh (baru):** `src/features/scoring/queries/keys.ts`,
  `src/features/scoring/queries/scoring.queries.ts`
- **Isi:**
  ```ts
  export const scoringKeys = {
    all:         ()               => ['scoring'] as const,
    score:       ()               => [...scoringKeys.all(), 'score'] as const,
    streak:      ()               => [...scoringKeys.all(), 'streak'] as const,
    leaderboard: (limit: number)  => [...scoringKeys.all(), 'leaderboard', limit] as const,
  }
  export function useScore()
  export function useStreak()
  export function useLeaderboard(limit?: number)
  ```
- **Aturan:** sama seperti `questKeys` — **dilarang** menulis `['scoring', ...]`
  literal di luar file ini. F3.6 memanggil key-key ini dari fitur lain.
- **DoD:** 3 hook + key terpusat.
- **Verifikasi:** `grep -rn "'scoring'" src/ | grep -v keys.ts` → kosong.

## F3.3 — `ScoreCard`

- **Sentuh (baru):** `src/features/scoring/components/ScoreCard.tsx`
- **Isi:** `Card` berisi: total poin (angka besar), **Level N** sebagai badge,
  dan `Progress` bar menuju level berikutnya dengan keterangan
  "120 XP lagi menuju Level 4".
- **Aturan keras:** **jangan menghitung ulang level atau ambang XP di frontend.**
  Level datang dari backend (`ScoreResponse.level`), sisa XP dari
  `points_to_next_level`. Aturan level adalah satu sumber di
  `scoring.LevelForXP` (ADR-007) — menduplikasi rumusnya di sini artinya dua
  tempat harus diubah bersamaan, dan cepat atau lambat melenceng.
- **Catatan:** kalau `points_to_next_level` tidak ada di kontrak, **tambahkan ke
  kontrak & backend**, jangan hitung sendiri.
- **DoD:** kartu tampil benar termasuk saat skor 0 (level 1, bar kosong).
- **Verifikasi:** render dengan data 0, tengah, dan tepat di ambang level.

## F3.4 — `StreakCard`

- **Sentuh (baru):** `src/features/scoring/components/StreakCard.tsx`
- **Isi:** `Card` berisi streak berjalan (angka besar + ikon api) dan streak
  terpanjang sebagai teks pendukung. Saat `current === 0`, tampilkan ajakan
  ("Selesaikan satu quest untuk memulai streak") alih-alih angka nol yang dingin.
- **Aturan:** `last_active` ditampilkan apa adanya dari backend (sudah tanggal
  lokal user) — jangan diproses ulang dengan timezone browser.
- **DoD:** tampil benar untuk streak 0, 1, dan besar.
- **Verifikasi:** render dengan tiga data itu.

## F3.5 — Pasang kartu di Dashboard

- **Sentuh:** `src/pages/DashboardPage.tsx`
- **Isi:** isi slot yang disiapkan di F2.7 — grid 2 kolom (desktop) berisi
  `ScoreCard` + `StreakCard` di atas daftar quest hari ini. Skeleton saat loading.
- **Aturan:** Dashboard tetap tipis — hanya merakit komponen fitur, tak
  mengandung logika sendiri.
- **DoD:** dashboard menampilkan skor, streak, dan quest hari ini sekaligus.
- **Verifikasi:** buka `/`.

## F3.6 — Invalidasi lintas-fitur setelah complete/uncomplete

- **Sentuh:** `src/features/quest/queries/quest.queries.ts`
- **Isi:** pada `onSettled` milik `useCompleteQuest` & `useUncompleteQuest`,
  tambahkan:
  ```ts
  queryClient.invalidateQueries({ queryKey: questKeys.today() })
  queryClient.invalidateQueries({ queryKey: scoringKeys.score() })
  queryClient.invalidateQueries({ queryKey: scoringKeys.streak() })
  // leaderboard: cukup invalidate scoringKeys.all() kalau ingin ikut segar
  ```
  Hapus komentar `// TODO(F3.6)` dari F2.3.
- **Kenapa task tersendiri:** **inilah alasan memilih TanStack Query.** Satu aksi
  (centang quest) mengubah tiga sumber data di server. Dengan zustand murni,
  ketiganya harus di-refetch manual di tiap tempat yang memanggil complete — dan
  satu tempat yang lupa berarti angka di layar bohong.
- **Aturan:**
  - `quest.queries` boleh mengimpor `scoringKeys` (hanya konstanta key), tapi
    **jangan** mengimpor komponen atau hook milik `features/scoring`. Batas fitur
    dijaga: yang dibagi hanyalah key.
  - Invalidasi di `onSettled`, bukan `onSuccess` — supaya setelah gagal pun UI
    kembali sinkron dengan server.
- **DoD:** centang quest → poin & streak di dashboard ikut berubah tanpa reload.
- **Verifikasi:** centang quest sambil melihat kedua kartu; lalu uncomplete dan
  pastikan poin turun kembali.

## F3.7 — Halaman Leaderboard

- **Sentuh (baru):** `src/pages/LeaderboardPage.tsx`,
  `src/features/scoring/components/LeaderboardTable.tsx`
- **Isi:** `Table` dengan kolom rank, nama, poin. Rank 1-3 diberi badge/warna
  khusus. **Baris user sendiri disorot** (bandingkan `entry.user_id` dengan
  `auth.store.user.id`).
- **Aturan:**
  - Rank datang dari backend, **jangan dihitung ulang dari indeks array** —
    kalau backend nanti menambah paging, indeks lokal akan salah.
  - Nama yang kosong (user terhapus) ditampilkan sebagai fallback yang sopan,
    bukan `undefined` — backend sudah menjaga ini (backend T3.6), frontend jangan
    menganggapnya mustahil.
  - Empty state: "Belum ada yang mengumpulkan poin".
- **DoD:** leaderboard tampil dengan nama, rank, dan sorotan diri sendiri.
- **Verifikasi:** buka `/leaderboard` dengan ≥2 user di mock.

## F3.8 — MSW handler scoring

- **Sentuh:** `src/mocks/handlers.ts`
- **Isi:** 3 endpoint scoring, **terhubung dengan state quest** di F2.11:
  complete quest di mock → poin & streak ikut naik.
- **Aturan:** kalau mock scoring statis sementara mock quest stateful, F3.6 akan
  terlihat "berhasil" padahal tak membuktikan apa-apa. Sambungkan keduanya.
- **DoD:** alur centang → angka naik bisa didemokan tanpa backend.
- **Verifikasi:** `VITE_USE_MOCK=true npm run dev`

## F3.9 — Test phase 3

- **Sentuh (baru):** `src/features/scoring/components/ScoreCard.test.tsx`,
  `src/features/quest/queries/invalidation.test.tsx`
- **Isi:**
  - `ScoreCard`: render skor 0 dan skor tengah; **tak ada perhitungan level di
    komponen** (nilai yang tampil sama persis dengan props).
  - Invalidasi: setelah `useCompleteQuest` selesai, `questKeys.today()`,
    `scoringKeys.score()`, dan `scoringKeys.streak()` ketiganya ditandai stale
    (spy pada `queryClient.invalidateQueries`).
  - `LeaderboardTable`: baris milik user sendiri mendapat penanda.
- **DoD:** `npm run test` hijau.
- **Verifikasi:** `npm run test`

---

## Exit criteria Phase 3

- [ ] `npm run lint && npm run typecheck && npm run test && npm run build` bersih.
- [ ] Dashboard menampilkan poin, XP, level, dan streak.
- [ ] Centang quest → ketiga angka ikut segar tanpa reload; uncomplete →
      poin turun kembali.
- [ ] Leaderboard menampilkan rank & nama, baris sendiri tersorot.
- [ ] Tak ada rumus level/XP yang ditulis ulang di frontend.
- [ ] Tak ada query key literal di luar kedua `keys.ts`.
- [ ] `features/quest` hanya mengimpor `scoringKeys` dari `features/scoring` —
      bukan komponen atau hook-nya.
