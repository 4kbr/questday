# Hands-off — State `apps/frontend`

Dokumen orientasi cepat khusus **frontend**. Untuk state seluruh proyek lihat
`docs/HANDS-OFF.md` (root). Baca `AGENTS.md` dulu, lalu ini.

**Terakhir diperbarui:** 2026-09-01 · **Branch:** `feature/frontend/master`

---

## TL;DR

**Phase 0 (Setup & fondasi) + Phase 1 (Auth & shell) + Phase 2 (Quest) +
Phase 3 (Scoring) SELESAI & terverifikasi.** Pondasi (Vite + tema + lapisan API +
router + MSW + test) berdiri, di atasnya sesi auth penuh (register/login/logout,
rute terlindungi, shell SaaS), inti aplikasi quest (dashboard hari ini dengan
centang optimistik + rollback, CRUD definisi quest, MSW quest handler stateful),
dan kini gamifikasi: kartu score/streak di dashboard, halaman leaderboard dengan
sorotan baris sendiri, dan invalidasi lintas-fitur sehingga centang quest
langsung menggeser poin + streak tanpa reload.

**Belum di-commit.** Semua file `apps/frontend/*` masih untracked/termodifikasi.
`docs/tasks/`, `contracts/`, `AGENTS.md`, root `docs/` tak tersentuh (selain
dokumen ini + `docs/tasks/README.md` untuk papan progres).

**Phase 3 Batch 1 (F3.1–F3.5) SELESAI** — `gen:api` diregenerasi (kontrak
menghapus response 404 di score/streak), `apis/scoring.api.ts` (HTTP murni),
`features/scoring/queries/{keys.ts,scoring.queries.ts}` (`useScore` / `useStreak`
/ `useLeaderboard`), komponen presentational `ScoreCard` + `StreakCard`, dan
keduanya terpasang di `DashboardPage` (grid 2 kolom, skeleton saat loading,
fallback inline saat error). Tak ada rumus level/XP di FE — bar ScoreCard cuma
rasio tampilan `xp / (xp + points_to_next_level)`.

**Phase 3 Batch 2 (F3.6–F3.8) SELESAI** —
- **F3.6**: `useCompleteQuest` / `useUncompleteQuest` di `quest.queries.ts` kini
  meng-invalidate `scoringKeys.score()` + `scoringKeys.streak()` +
  `scoringKeys.all()` (leaderboard) di `onSettled` (bukan `onSuccess` — gagal pun
  re-sync). Import `scoringKeys` via path relatif `../../scoring/queries/keys`
  (KONSTANTA key saja — bukan hook/komponen scoring). Komentar `// TODO(F3.6)`
  dihapus. 5 hook quest lain tetap invalidate `questKeys.all()`.
- **F3.7**: `pages/LeaderboardPage.tsx` (nyata) + `features/scoring/components/
  LeaderboardTable.tsx`. Tabel kolom Rank/Nama/Poin; `rank` dari `entry.rank`
  (BUKAN indeks array); rank 1-3 dapat ikon `Medal` (emas/perak/perunggu); baris
  user sendiri `bg-muted/50 font-medium` + penanda "(kamu)"; nama kosong →
  "Pengguna dihapus". Page: `useLeaderboard()` + `useAuthStore((s)=>s.user?.id)`,
  skeleton loading, error + "Coba lagi" → `refetch()`, empty →
  "Belum ada yang mengumpulkan poin". Tanpa chrome layout (AppShell dari route).
- **F3.8**: MSW handler scoring stateful di `mocks/handlers.ts` — `GET /me/score`,
  `GET /me/streak`, `GET /leaderboard`. Terhubung ke `todayCompleted` + `quests`
  + `POINTS_BY_DIFFICULTY`: `sumCompletedPoints()`, `computeScore()` (rumus level
  MOCK-ONLY), `computeStreak()` (+ `longestStreakSeen` module-level). Leaderboard
  = seeded user (poin = skor live) + 2 user palsu tetap (Budi 45 / Sari 15),
  sort desc, rank 1-based setelah sort, slice ke `limit` (query param, default
  20). `/me/score` & `/me/streak` tak lagi 404 di mock. Centang/uncentang quest
  di mock kini menggeser score + streak + leaderboard sekaligus.

**Phase 3 Batch 3 (F3.9) SELESAI** — 3 file test baru, semua hijau:
- `features/scoring/components/ScoreCard.test.tsx` — render skor 0, tengah, dan
  satu kasus di mana `level` / `points_to_next_level` dari props SENGAJA tak cocok
  dengan rumus lokal naif (`floor(xp/seratus)+1`); test gagal kalau ada yang
  menyisipkan rumus level di komponen. Bar hanya dicek keberadaannya (role
  `progressbar`), bukan nilai persisnya.
- `features/quest/queries/invalidation.test.tsx` — `msw/node` `setupServer`
  (`GET /quests/today`, `POST /quests/:id/complete`), `vi.spyOn(queryClient,
  'invalidateQueries')`, `useCompleteQuest().mutate()` → assert dipanggil dengan
  `questKeys.today()`, `scoringKeys.score()`, `scoringKeys.streak()`. Import lintas
  -fitur via path relatif (`../../scoring/queries/keys`). Spy di-restore di
  `afterEach`.
- `features/scoring/components/LeaderboardTable.test.tsx` — baris user sendiri
  dapat `bg-muted/50 font-medium` + "(kamu)"; `rank` dibaca dari `entry.rank`
  (entries dengan `rank` tak urut → baris pertama tetap rank 2); `display_name`
  kosong → "Pengguna dihapus".

**Total test: 30 (11 file), semua hijau.**

### Phase 3 selesai — catatan lintas-batch

- **Invalidasi lintas-fitur** ada di `quest.queries.ts` `onSettled` milik
  `useCompleteQuest` / `useUncompleteQuest`: invalidate `scoringKeys.score()` +
  `scoringKeys.streak()` + `scoringKeys.all()` (leaderboard) + `questKeys.today()`.
  Komentar `// TODO(F3.6)` sudah dihapus. `scoringKeys` di-import via path relatif
  `../../scoring/queries/keys` — KONSTANTA key saja, bukan hook/komponen scoring.
- **MSW handler scoring stateful** di `handlers.ts` (`computeScore` /
  `computeStreak`), terhubung ke `todayCompleted` — centang quest di mock
  menggeser score + streak + leaderboard sekaligus.
- **Rumus level/XP MOCK-ONLY** hidup di `handlers.ts`. Frontend tak pernah
  menghitung level/XP; `ScoreCard` cuma rasio tampilan
  `xp / (xp + points_to_next_level)` untuk bar.
- **Sorotan baris sendiri di leaderboard** via `useAuthStore((s) => s.user?.id)`
  di `LeaderboardPage`, diteruskan sebagai `currentUserId` ke `LeaderboardTable`.

**Mulai dari sini:** `apps/frontend/docs/tasks/phase-04-polish.md` **F4.1**
(`updateMe` + `useUpdateProfile` — butuh backend T1.11).

---

## Progres per phase

| Phase | File task | Status |
|---|---|---|
| 0 — Setup & fondasi | `docs/tasks/phase-00-setup.md` | ✅ **selesai** (F0.1–F0.12) |
| 1 — Auth & shell | `docs/tasks/phase-01-auth.md` | ✅ **selesai** (F1.1–F1.11) |
| 2 — Quest | `docs/tasks/phase-02-quest.md` | ✅ **selesai** (F2.1–F2.12) |
| 3 — Scoring | `docs/tasks/phase-03-scoring.md` | ✅ **selesai** (F3.1–F3.9) |
| 4 — Settings & polish | `docs/tasks/phase-04-polish.md` | ⬜ belum mulai — mulai dari **F4.1** |

Papan progres per-task ada di `docs/tasks/README.md`.

---

## Stack terpasang (⚠️ lebih baru dari asumsi task docs)

`create-vite` v9 men-scaffold stack yang lebih baru dari yang ditulis task docs.
Semua bekerja; batch berikutnya harus sadar ini:

| Bagian | Versi / pilihan | Catatan |
|---|---|---|
| Core | **Vite 8.2**, **React 19.2**, **TypeScript 6.0** | `@vitejs/plugin-react` 6 |
| tsconfig | split (`app`/`node`), `verbatimModuleSyntax`, **`erasableSyntaxOnly`** | TS "parameter properties" di constructor **dilarang** — pakai field eksplisit |
| Styling | **Tailwind v4** (`@tailwindcss/vite`) + **shadcn/ui** v4 | style `radix-nova`, base color `neutral`, font Inter |
| Server state | **TanStack Query 5** | `src/lib/query-client.ts` |
| Client state | **Zustand** | belum dipasang — lahir di F1.2 |
| HTTP | **axios 1.20** | hanya di `src/apis/client.ts` |
| Types | **openapi-typescript 7** | butuh `--legacy-peer-deps` (peer minta TS 5.x) |
| Routing | **react-router-dom 7** | `createBrowserRouter` |
| Lint | **oxlint** (BUKAN ESLint) | keputusan pemilik kode; `.oxlintrc.json` |
| Format | **Prettier** | format saja, bukan linting |
| Mock | **MSW 2** | di balik flag `VITE_USE_MOCK` |
| Test | **Vitest 4** + Testing Library + jsdom | |

---

## Keputusan yang mengikat (dibuat selama Phase 0 — jangan ditabrak diam-diam)

1. **Response sukses BER-AMPLOP `{data}` (ADR-025).** Setiap response sukses
   berbadan JSON dibungkus `{"data": <payload>}`; `src/apis/client.ts`
   meng-unwrap sekali di interceptor sehingga `res.data` = payload langsung
   sesuai schema. `register` mengembalikan **200** (bukan 201). `/healthz` polos
   (di luar amplop) — dibiarkan apa adanya. **Error** tetap ber-amplop:
   `{"error":{code,message}}` → diubah jadi `ApiError(status, code, message)` di
   `src/apis/client.ts`. Fitur membedakan perlakuan lewat **`ApiError.code`**
   (`already_completed`, `email_taken`, dst — lihat deskripsi `ErrorResponse` di
   `contracts/openapi.yaml`), **bukan** mencocokkan string pesan.
   (Membalik keputusan awal Phase 0 yang mengasumsikan tak ada `{data}`.)
2. **oxlint, bukan ESLint** (menyimpang dari F0.10, disetujui pemilik kode).
   Prettier tetap dipakai untuk format.
3. **`no-restricted-imports` mem-blokir SELURUH `@/features/*`.** oxlint tak bisa
   membatasi "fitur lain" relatif terhadap file saat ini. Konsekuensi:
   **import di dalam satu fitur WAJIB relatif (`./`)**, bukan `@/features/...`.
   Kode bersama tetap lewat `@/components/*` atau `@/lib/*`. Ada `TODO(oxlint)`
   di `.oxlintrc.json`.
4. **`.npmrc` → `legacy-peer-deps=true`** supaya `npm install` polos jalan untuk
   dev baru (konflik `openapi-typescript@7` ↔ `typescript@6`).
5. **`typecheck` = `tsc -b --noEmit`** (root tsconfig solution-style; `tsc
   --noEmit` polos lolos tanpa mengecek apa-apa).
6. **shadcn CLI di `devDependencies`.** `next-themes` + `tw-animate-css` sengaja
   tetap di `dependencies` (runtime).

Semua keputusan ini juga terangkum di `apps/frontend/README.md`.

### Rekonsiliasi pasca-merge (2026-09-01)

`main` di-merge membawa **ADR-025**. Penyesuaian frontend: `schema.gen.ts`
di-regenerate (response 2xx kini `{data: <schema>}`), `client.ts` sekarang
meng-unwrap `{data}` di interceptor sukses, key localStorage default diselaraskan
ke `questday-auth` (hyphen, cocok dgn `persist` F1.2), dan `VITE_PORT` didukung
(`.env.example` + `vite.config.ts` + `vite-env.d.ts`). `register` kini 200.

---

## Yang sudah ada di `src/`

```
src/
├── apis/
│   ├── client.ts        # axios instance TUNGGAL + ApiError + setTokenGetter/setUnauthorizedHandler
│   ├── client.test.ts   # map error→ApiError, success raw, header Bearer
│   ├── auth.api.ts       # register / login / me — murni HTTP (F1.1)
│   ├── schema.gen.ts     # GENERATE dari contracts/openapi.yaml — jangan diedit tangan
│   └── types.ts          # alias ramah (User, AuthResponse, Quest, Score, Streak, ...)
├── lib/
│   ├── query-client.ts  # QueryClient: staleTime 30s, no-retry 4xx, refetchOnWindowFocus
│   ├── session.ts        # endSession() — SATU jalur teardown (logout + queryClient.clear() + redirect); wiring setTokenGetter/setUnauthorizedHandler (F1.9)
│   └── utils.ts          # cn() dari shadcn
├── routes/
│   ├── paths.ts         # PATHS — konstanta, jangan tulis path string literal di komponen
│   ├── index.tsx         # createBrowserRouter — GuestRoute (login/register) + ProtectedRoute (nested) → 6 route
│   ├── ProtectedRoute.tsx / ProtectedRoute.test.tsx  # tanpa token → /login (state.from); ada token → <AppShell><Outlet/>
│   ├── GuestRoute.tsx    # sudah login → lempar ke dashboard
│   └── router.test.tsx
├── pages/                # 6 halaman tipis; Login/Register + Dashboard/Quests/Leaderboard nyata
│                         #   DashboardPage kini juga menampilkan ScoreCard + StreakCard (F3.5);
│                         #   LeaderboardPage nyata (F3.7); Settings masih placeholder
├── apis/
│   ├── quest.api.ts      # list/today/create/update/archive/complete/uncomplete — murni HTTP (F2.1)
│   └── scoring.api.ts    # score / streak / leaderboard — murni HTTP (F3.1)
├── mocks/
│   ├── handlers.ts      # /healthz + auth + QUEST stateful (7 handler, store in-memory, jalur 404/409) MSW; helper errorBody/errorResponse + dataResponse ({data} envelope, ADR-025)
│   └── browser.ts        # setupWorker(...handlers)
├── components/
│   ├── ui/               # 18+ komponen shadcn (generated) — termasuk progress + alert-dialog (F2)
│   ├── EmptyState.tsx    # ikon + judul + deskripsi + aksi — dipakai dashboard/quests kosong (F2.10)
│   └── layout/           # AppShell (+ <Toaster/> sonner) + Sidebar + Topbar — shell SaaS semua route terproteksi (F1.8)
├── features/
│   ├── auth/
│   │   ├── components/   # LoginForm (+ test), RegisterForm
│   │   ├── queries/      # authKeys, useMe, useLogin, useRegister (satu-satunya pemanggil authApi)
│   │   ├── schemas/      # loginSchema / registerSchema (zod, cocok validasi backend)
│   │   └── lib/          # timezones.ts (daftar IANA + default browser)
│   ├── quest/
│   │   ├── components/   # QuestItem (+ test), TodayQuestList, QuestFormDialog (+ test), QuestTable
│   │   ├── queries/      # questKeys (keys.ts) + 7 hook (+ test); complete/uncomplete optimistik di questKeys.today(); isBenignQuestToggleError
│   │   ├── schemas/      # questFormSchema / QuestFormValues (zod)
│   │   └── lib/          # difficulty.ts (warna + label badge easy/medium/hard)
│   └── scoring/
│       ├── components/   # ScoreCard (+ test), StreakCard, LeaderboardTable (+ test)
│       └── queries/      # scoringKeys (keys.ts) + useScore / useStreak / useLeaderboard
├── stores/
│   ├── auth.store.ts     # zustand + persist key 'questday-auth' — token & user (client state saja)
│   └── auth.store.test.ts
├── test/setup.ts         # @testing-library/jest-dom/vitest
├── main.tsx              # enableMocking() → import '@/lib/session' → QueryClientProvider → RouterProvider → Devtools(DEV)
├── index.css             # Tailwind v4 + token tema shadcn (:root + .dark) + Inter
└── vite-env.d.ts         # tipe VITE_API_BASE_URL, VITE_USE_MOCK, VITE_PORT?
```

### Phase 1 selesai — yang berdiri sekarang

- **Route guard bersarang**: `GuestRoute` membungkus `/login` + `/register`
  (redirect ke dashboard bila sudah punya token); `ProtectedRoute` membungkus
  sisanya, merender `<AppShell>` + `<Outlet/>` dan menyimpan `state.from` supaya
  user kembali ke halaman yang tadi dituju setelah login.
- **`AppShell`** (Sidebar + Topbar) — shell untuk SEMUA halaman terproteksi;
  halaman berikutnya (Phase 2+) tak menggambar sidebar/topbar sendiri.
- **`endSession()` (`src/lib/session.ts`)** — satu-satunya jalur teardown sesi,
  dipakai tombol Logout DAN interceptor 401 di `client.ts`: `logout()` +
  `queryClient.clear()` (wajib, cegah kebocoran data antar-akun) + redirect ke
  `PATHS.login`. Modul ini juga mem-wire `setTokenGetter` / `setUnauthorizedHandler`
  saat di-import `main.tsx`.
- **MSW auth handlers** (`src/mocks/handlers.ts`) — `POST /auth/login`,
  `POST /auth/register` (200, bukan 201), `GET /me`, plus jalur gagal
  (401 `invalid_credential`, 409 `email_taken`). Sukses dibungkus amplop
  `{ data: <payload> }` lewat helper `dataResponse()` (ADR-025); error tetap
  `{ error: { code, message } }`.
- **Persist**: zustand `persist` key **`questday-auth`**, hanya `{token, user}`;
  refresh browser tetap login.

### Phase 2 selesai — yang berdiri sekarang

- **MSW quest handler stateful** (`src/mocks/handlers.ts`) — 7 endpoint quest di
  atas store in-memory (`quests[]` + `todayCompleted` Set), sukses ber-amplop
  `{ data }` (`dataResponse`), error `{ error:{code,message} }` (`errorResponse`),
  204 = `new HttpResponse(null,{status:204})`. `/quests/today` didaftar SEBELUM
  `/quests/:questId`. Jalur gagal: `quest_not_found` (404), `already_completed` /
  `not_completed` (409), `validation_failed` (400). Poin turun dari difficulty
  (easy 5 / medium 10 / hard 20). Complete/uncomplete benar-benar mengubah
  `GET /quests/today` → optimistic update teruji sungguhan.
- **Optimistic complete/uncomplete + rollback** di `features/quest/queries`:
  `onMutate` menulis `questKeys.today()` setelah `cancelQueries`, `onError`
  mengembalikan `ctx.prev`, `onSettled` `invalidateQueries`. Kode 409
  (`already_completed` / `not_completed`) TIDAK dianggap error merah — caller
  menyaring lewat `isBenignQuestToggleError(err)` (baca `ApiError.code`).
- **`questKeys` (`features/quest/queries/keys.ts`)** satu-satunya sumber query
  key quest — dilarang menulis `['quests', ...]` di file lain.
- ~~`// TODO(F3.6): invalidate scoring`~~ SELESAI di F3.6 (batch 2) —
  `useCompleteQuest` / `useUncompleteQuest` kini invalidate scoring keys.
- **Tanggal dashboard** diambil dari `TodayQuestsResponse.date` (backend hitung
  dari timezone user, ADR-006), bukan `new Date()` browser.
- **Test**: `QuestItem.test.tsx`, `QuestFormDialog.test.tsx` (seam-mock hook),
  `quest.queries.test.tsx` (`msw/node` `setupServer` + `queryClient` singleton).
  Setup test menambah polyfill `window.matchMedia` (jsdom) untuk komponen yang
  menyentuh shell / sonner.

Config root `apps/frontend/`: `.env.example`, `.gitignore`, `.npmrc`,
`.oxlintrc.json`, `.prettierrc`, `.prettierignore`, `components.json`,
`vite.config.ts`, `vitest.config.ts`, `tsconfig*.json`.

Script: `dev`, `build`, `preview`, `lint`, `format`, `format:check`, `test`,
`test:watch`, `gen:api`, `typecheck`.

---

## Cara jalan

```bash
cd apps/frontend
npm install                 # .npmrc sudah set legacy-peer-deps
cp .env.example .env
npm run gen:api             # regenerate src/apis/schema.gen.ts dari kontrak
npm run dev                 # http://localhost:5173  (VITE_USE_MOCK=true → MSW aktif)
```

Gate verifikasi (semua hijau per 2026-09-01):
`npm run typecheck && npm run lint && npm run test && npm run format:check && npm run build`

---

## Titik sambung ke backend (belum aktif)

Frontend jalan penuh dgn `VITE_USE_MOCK=true` sampai endpoint backend siap.

| Butuh | Dari backend | Dipakai FE saat |
|---|---|---|
| `contracts/openapi.yaml` terisi schema | ✅ sudah (16 schema) | F0.4 — **tidak lagi memblokir** |
| Endpoint auth jalan | Phase 1 backend | FE Phase 1 lepas mock |
| Endpoint quest jalan | Phase 2 backend | FE Phase 2 lepas mock |
| Endpoint scoring jalan | Phase 3 backend | FE Phase 3 lepas mock |
| `PATCH /me` | backend T1.11 | FE Phase 4 |

> Catatan: root `docs/HANDS-OFF.md` masih bilang `components.schemas: {}` — itu
> **stale**. Kontrak sudah terisi; `npm run gen:api` menghasilkan 16 schema.

---

## Pesan untuk agent berikutnya (Phase 1)

- **F1.2 `auth.store`**: zustand `persist` **WAJIB** `name: 'questday-auth'` dan
  simpan token di `state.token` — default `tokenGetter` di `client.ts` membaca
  key itu. Setelah store jadi, panggil `setTokenGetter(() => useAuth.getState().token)`
  dan `setUnauthorizedHandler(() => { logout(); redirect ke PATHS.login })` sekali
  saat bootstrap (mis. di `main.tsx` atau modul store).
- **F1.9 Logout**: wajib `queryClient.clear()` selain hapus token (jebakan #8 di
  root HANDS-OFF).
- HTTP baru → tambah fungsi di `src/apis/*.api.ts` (buat file baru), bungkus
  dengan hook di `src/features/*/queries/`. Komponen **tak pernah** sentuh `api`
  langsung.
- Route baru → tambah di `src/routes/paths.ts` + `src/routes/index.tsx`. Jangan
  string literal.
- Type API kurang → perbaiki `contracts/openapi.yaml`, jalankan `npm run gen:api`,
  tambah alias di `src/apis/types.ts`. Jangan bikin `interface` tangan.
- Keputusan desain baru → tambah ADR di root `docs/DECISIONS.md`, dan perbarui
  dokumen ini.
