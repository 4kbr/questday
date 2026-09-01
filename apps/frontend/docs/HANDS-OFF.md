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

**Total test: 30 (11 file) → 39 (15 file) setelah Phase 4 Batch 4, semua hijau.**

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

**Phase 4 Batch 1 (F4.1–F4.3) SELESAI** —
- **F4.1**: `authApi.updateMe(body) → api.patch<AuthResponse>('/me', body)`.
  `features/auth/queries/profile.queries.ts` baru: `useUpdateProfile()` —
  `onSuccess({token,user})` → `useAuthStore.getState().setSession(token,user)`
  (token BARU, ADR-022/013) lalu invalidate `authKeys.me()` + `questKeys.today()`
  + `scoringKeys.all()`. `scoringKeys`/`questKeys` di-import via path relatif
  (`../../quest/queries/keys`, `../../scoring/queries/keys`) — konstanta key saja.
- **F4.2**: `timezones.ts` DIPINDAH `features/auth/lib/` → `src/lib/timezones.ts`
  (lintas-fitur; `RegisterForm` sekarang `import { browserTimezone } from
  '@/lib/timezones'`). `src/components/TimezoneSelect.tsx` baru — combobox
  `Popover`+`Command` (shadcn `command`/`popover`/`sheet` ditambah via CLI;
  `textarea`+`input-group` ikut sebagai dep, `button`/`input`/`dialog`
  di-refresh). Daftar zona dari `Intl.supportedValuesOf('timeZone')`, fallback
  `COMMON_TIMEZONES`. `features/auth/components/ProfileForm.tsx` baru (RHF +
  `profileSchema`), `pages/SettingsPage.tsx` nyata (Card "Profil"). `RegisterForm`
  field timezone kini pakai `<TimezoneSelect>` (import `Select*` + `timezoneOptions`
  dihapus).
- **F4.2 schema**: `auth.schema.ts` tambah `profileSchema`
  (`display_name` min 1, `timezone` min 1), `ProfileValues`, `_ProfileSyncCheck`
  (resolve → `true` vs `UpdateProfileRequest`).
- **F4.3**: tercakup di `useUpdateProfile().onSuccess` (setSession token baru +
  invalidate). Mock `PATCH ${BASE}/me` di `handlers.ts` (dekat `GET /me`) —
  mutasi `seededUser.display_name`/`.timezone`, balas `{token: MOCK_TOKEN, user:
  seededUser}`.
- Gate hijau: typecheck / lint (0 error) / test (30, tak nambah) / build /
  format:check. `grep @/features/` & `grep 'quests'|'scoring'` (di luar keys.ts)
  kosong.

**Phase 4 Batch 2 (F4.4–F4.6) SELESAI** —
- **F4.4 dark mode**: `src/stores/ui.store.ts` baru (zustand persist key
  `questday-ui`, `partialize` → `{theme}`, `theme: 'light'|'dark'|'system'`).
  Ekspor `resolveTheme` / `applyTheme` (toggle class `dark` di `<html>`).
  `onRehydrateStorage` panggil `applyTheme` saat import; listener `matchMedia`
  module-level ikut perubahan OS saat mode `system`.
  `src/components/layout/ThemeToggle.tsx` baru — `DropdownMenu` (trigger icon
  button `aria-label="Ubah tema"`, ikon `Sun`/`Moon`/`Monitor` per `theme`), 3
  item Terang/Gelap/Sistem, aktif ditandai `font-medium` + `Check`. Dipasang di
  `Topbar` dalam `div.flex.items-center.gap-2` kiri dari avatar dropdown.
  `main.tsx` `import '@/stores/ui.store'` (side-effect). Anti-flash: IIFE inline
  di `<head>` `index.html` (sebelum font `<link>`) baca `localStorage`
  `questday-ui` → set class `dark` sebelum render pertama. Konfirmasi bertahan
  di `dist/index.html` (`grep -c questday-ui` = 1).
- **F4.5 toast seragam**: `src/lib/toast.ts` baru — `toastSuccess(msg)` +
  `toastApiError(err, fallback)` (baca `ApiError.message`; `status === 0` →
  "Tidak dapat terhubung ke server"; non-ApiError → fallback; tak mengarang
  detail). Sukses di-wire: `QuestFormDialog` create → "Quest dibuat", edit →
  "Quest diperbarui" (di `onSuccess`, sebelum tutup dialog); `QuestsPage`
  archive → "Quest diarsipkan"; `ProfileForm` → "Profil disimpan".
  Error toast HANYA di `QuestsPage` archive (`onError` → `toastApiError`) karena
  tak ada UI inline di sana. `QuestFormDialog` + `ProfileForm` TETAP pakai alert
  inline saja untuk error (dialog/form terlihat; hindari pesan ganda) — tak
  ditambah error toast. Toggle complete/uncomplete TIDAK dapat toast (F2.6;
  sudah kasat mata; `already_completed` senyap). `TodayQuestList` `toast.error`
  lama dibiarkan (refactor opsional dilewati). Hook tetap bebas toast.
- **F4.6 ErrorBoundary + 404**: `src/components/ErrorBoundary.tsx` baru (class
  component — diizinkan meski `erasableSyntaxOnly`; field `state` eksplisit,
  bukan param-property), `getDerivedStateFromError` + `componentDidCatch`
  (`console.error` hanya DEV), fallback "Ada yang salah" + tombol Muat ulang,
  `<pre>` detail hanya DEV. `src/pages/NotFoundPage.tsx` baru (default export,
  "404 — Halaman tak ditemukan" + `Link` ke `PATHS.dashboard`, di luar
  AppShell → centering minimal). `src/routes/RouteError.tsx` baru
  (`useRouteError()` → pesan + link dashboard, detail hanya DEV).
  `routes/index.tsx`: `errorElement: <RouteError />` di KEDUA grup route +
  route catch-all `{ path: '*', element: <NotFoundPage /> }` top-level.
  `main.tsx`: `<RouterProvider>` dibungkus `<ErrorBoundary>` di dalam
  `QueryClientProvider`, `<StrictMode>` tetap terluar.
- Gate hijau: typecheck / lint (0 error; 3 warning lama di `ui/`) / test (30,
  tak nambah) / build / format:check. `grep @/features/` di `src/` kosong.
  Format via `npm run format` (prettier ubah gaya IIFE `index.html` ke
  semicolon-prefix — normal).

**Phase 4 Batch 3 (F4.7–F4.8) SELESAI** —
- **F4.7 responsive**: `src/components/layout/SidebarNav.tsx` baru — array `NAV`
  + `<nav>` + `<NavLink>` (styling aktif `cn(...)` tetap), prop
  `{ onNavigate?: () => void }` dipanggil `onClick` tiap link (drawer mobile
  menutup sheet). Brand "QuestDay" TIDAK di `SidebarNav` — tiap kontainer
  menggambarnya sendiri. `Sidebar.tsx` jadi desktop-only:
  `aside` `hidden md:flex` (dulu `flex`), merender `<SidebarNav />`. `Topbar.tsx`
  dapat hamburger `Button variant=ghost size=icon` `aria-label="Buka menu"`
  `className="md:hidden"` di kiri jauh, membuka `Sheet side=left w-64 p-0`
  (`SheetHeader h-14` + `SheetTitle` QuestDay, lalu `<SidebarNav onNavigate={()
  => setOpen(false)} />`); `const [open,setOpen]=useState(false)`. Grup kiri
  Topbar `div.flex.min-w-0.items-center.gap-2`, `<h1>` kini `truncate`.
  `AppShell` `<main>` → `p-4 md:p-6`. `QuestTable` wrapper dapat `w-full`
  (sudah `overflow-x-auto rounded-lg border`); `LeaderboardTable` dibungkus
  `<div className="w-full overflow-x-auto">` baru.
- **F4.8 a11y**: `TimezoneSelect` sudah forward `id` ke trigger + terima
  `aria-label` (tak berubah). `RegisterForm` field timezone kini
  `FormLabel htmlFor="register-timezone"` + `TimezoneSelect id="register-timezone"
  aria-label="Timezone"`; `ProfileForm` `TimezoneSelect` dapat `aria-label="Timezone"`
  (id/htmlFor sudah ada). Avatar `DropdownMenuTrigger` di `Topbar` dapat
  `aria-label="Menu akun"` (ring focus-visible sudah ada). Audit `outline-none`:
  semua hit di `src/components/ui/*` (shadcn) berpasangan `focus-visible:ring-*`
  atau kontainer non-interaktif; satu hit non-ui di `Topbar` avatar trigger
  berpasangan `focus-visible:ring-2 focus-visible:ring-ring`. Audit `onClick=`:
  semua di elemen button-like (`Button` / `DropdownMenuItem` / `AlertDialogAction`)
  — tak ada `div onClick`. Tak ada `onCloseAutoFocus` preventDefault di mana pun.
- Gate: typecheck bersih; lint 0 error (3 warning lama `ui/`); build sukses
  (`dist/` dihapus); format:check bersih; `grep @/features/` di `src/` kosong.
  Test: **29/30 hijau**; 1 gagal (`QuestFormDialog.test.tsx` "submit valid …")
  HANYA saat run penuh karena timeout worker 5000ms di mesin yang sangat
  terbebani (environment 164s) — lulus 2/2 saat file itu dijalankan sendiri
  (environment 8s). Bukan regresi: Batch 3 tak menyentuh `QuestFormDialog`.
  Tak ada file test yang diubah (tak ada test layout yang meng-assert
  visibilitas sidebar).

**Phase 4 Batch 4 (F4.10 + route code-splitting + test Phase 4 + docs) SELESAI** —
- **Route code-splitting**: 6 halaman di `routes/index.tsx` kini
  `const XPage = lazy(() => import('@/pages/XPage'))` (`NotFoundPage` tetap eager
  — catch-all top-level tanpa induk `Suspense`, ukurannya sepele). Satu
  `<Suspense fallback={<PageSkeleton/>}>` membungkus `<Outlet/>` di
  **`GuestRoute`** dan di **`ProtectedRoute` DI DALAM `AppShell`** (chrome
  sidebar/topbar tetap tampil saat chunk halaman diunduh) — bukan per-route.
  `src/components/PageSkeleton.tsx` baru (3 `Skeleton` + `p-6`, `aria-busy`).
  `routes/index.tsx` dapat `/* oxlint-disable react/only-export-components */`
  (file router, bukan file komponen). Bundle: **single chunk ~693 KB → entry
  365 KB** + ~15 chunk halaman/vendor terpisah (LoginPage 2.2 KB, DashboardPage
  14 KB, QuestsPage 9 KB, dst).
- **F4.10 build produksi**: `npm run build` bersih; `npm run preview -- --port
  4173` + `curl /quests` → **HTTP 200**, HTML shell berisi `id="root"` (deep
  link OK; router `*` catch-all + preview SPA fallback). `dist/index.html` tetap
  membawa anti-flash `questday-ui` (`grep -c` = 1). `dist/` dihapus setelah cek.
- **Test Phase 4** (4 file baru, semua hijau):
  `stores/ui.store.test.ts` (`setTheme` toggle class `dark` di `<html>`;
  `resolveTheme('system')` → `'light'` via matchMedia polyfill),
  `components/ErrorBoundary.test.tsx` (`<Boom/>` throw → fallback "Ada yang
  salah" + "Muat ulang", anak disembunyikan; `console.error` di-spy),
  `features/auth/components/ProfileForm.test.tsx` (QueryClientProvider + set
  `useAuthStore` user; Simpan disabled saat pristine → enabled setelah ketik →
  `mutate` dipanggil 1× dengan `{display_name, timezone}`; `vi.mock`
  `../queries/profile.queries`),
  `components/TimezoneSelect.test.tsx` (trigger tampil value; buka → ketik
  "Tokyo" di `CommandInput` → klik opsi `Asia/Tokyo` → `onChange('Asia/Tokyo')`;
  polyfill `scrollIntoView` + `ResizeObserver` untuk cmdk).
- **Total test: 39 (15 file), semua hijau** (dari 30/11).
- **README** `apps/frontend/README.md` dapat section "## Build produksi"
  (VITE_* inline saat build, code-splitting, deep link + SPA fallback rewrite).
- Gate hijau: typecheck / lint (0 error; 3 warning lama `ui/`) / test (39) /
  build / format:check. `grep -rn "@/features/" src/` kosong.

**Phase 4 SELESAI — MVP FRONTEND DONE.**

**F4.9 (uji melawan backend asli) — SELESAI 2026-09-01.** Smoke curl penuh lawan
backend Go asli (Postgres :5433, `go build ./cmd/api`, migrasi v1): register →
login → `GET /me` → `POST /quests` → `GET /quests/today` → complete (+ 409
`already_completed` pada complete kedua) → `GET /me/score` → `GET /me/streak` →
`GET /leaderboard` → `PATCH /me` (token BARU terbit) → login salah (401
`invalid_credential`) → PATCH quest asing (404 `quest_not_found`) → no-auth (401).
Semua **cocok** dengan MSW mock & kontrak, **kecuali satu**:

- **Mismatch:** middleware auth backend memakai `error.code: "unauthorized"`
  untuk semua 401 (token hilang/invalid/kedaluwarsa). Kontrak tak mendaftarkannya
  dan MSW `GET /me` no-token dulu memakai `invalid_credential`.
  **Perbaikan (bukan tambal FE):** (1) `contracts/openapi.yaml` `ErrorResponse`
  deskripsi kini mendaftarkan `unauthorized` + catatan bedanya dgn
  `invalid_credential`; (2) MSW `GET /me` no-token → `errorResponse(401,
  'unauthorized', ...)`. `npm run gen:api` dijalankan (nol perubahan tipe —
  `error.code` cuma string bebas). **Nol dampak fungsional**: `client.ts`
  menangani 401 by-status (`endSession()`), tak pernah cocokkan `error.code`
  untuk auth.

Backend dimatikan setelah smoke (server tak ditinggalkan jalan); container
Postgres dibiarkan (dikelola user).

---

## Progres per phase

| Phase | File task | Status |
|---|---|---|
| 0 — Setup & fondasi | `docs/tasks/phase-00-setup.md` | ✅ **selesai** (F0.1–F0.12) |
| 1 — Auth & shell | `docs/tasks/phase-01-auth.md` | ✅ **selesai** (F1.1–F1.11) |
| 2 — Quest | `docs/tasks/phase-02-quest.md` | ✅ **selesai** (F2.1–F2.12) |
| 3 — Scoring | `docs/tasks/phase-03-scoring.md` | ✅ **selesai** (F3.1–F3.9) |
| 4 — Settings & polish | `docs/tasks/phase-04-polish.md` | ✅ **selesai** (F4.1–F4.10) — **MVP FRONTEND DONE** |

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
│   ├── toast.ts          # toastSuccess / toastApiError — pemetaan mutation → sonner seragam (F4.5)
│   ├── timezones.ts      # COMMON_TIMEZONES + browserTimezone() — lintas-fitur (dipindah dari features/auth/lib, F4.2)
│   └── utils.ts          # cn() dari shadcn
├── routes/
│   ├── paths.ts         # PATHS — konstanta, jangan tulis path string literal di komponen
│   ├── index.tsx         # createBrowserRouter — 6 route LAZY (React.lazy + import()); NotFoundPage eager; catch-all '*'
│   ├── ProtectedRoute.tsx / ProtectedRoute.test.tsx  # tanpa token → /login (state.from); ada token → <AppShell><Suspense><Outlet/>
│   ├── GuestRoute.tsx    # sudah login → dashboard; else <Suspense fallback={<PageSkeleton/>}><Outlet/>
│   ├── RouteError.tsx    # errorElement kedua grup route — "Ada yang salah" + link dashboard, detail hanya DEV (F4.6)
│   └── router.test.tsx
├── pages/                # 7 halaman tipis (semua default export); Login/Register + Dashboard/Quests/Leaderboard/Settings nyata
│                         #   DashboardPage: ScoreCard + StreakCard (F3.5); LeaderboardPage nyata (F3.7);
│                         #   SettingsPage: Card "Profil" + ProfileForm (F4.2); NotFoundPage: route '*' 404 (F4.6)
├── apis/
│   ├── quest.api.ts      # list/today/create/update/archive/complete/uncomplete — murni HTTP (F2.1)
│   └── scoring.api.ts    # score / streak / leaderboard — murni HTTP (F3.1)
├── mocks/
│   ├── handlers.ts      # /healthz + auth + QUEST stateful (7 handler, store in-memory, jalur 404/409) MSW; helper errorBody/errorResponse + dataResponse ({data} envelope, ADR-025)
│   └── browser.ts        # setupWorker(...handlers)
├── components/
│   ├── ui/               # 20+ komponen shadcn (generated) — + command / popover / sheet / textarea / input-group (F4.2)
│   ├── EmptyState.tsx    # ikon + judul + deskripsi + aksi — dipakai dashboard/quests kosong (F2.10)
│   ├── ErrorBoundary.tsx (+ test)  # class component — bungkus RouterProvider; fallback "Ada yang salah" + Muat ulang (F4.6)
│   ├── PageSkeleton.tsx  # fallback <Suspense> untuk halaman lazy (route code-splitting, Batch 4)
│   ├── TimezoneSelect.tsx (+ test)  # combobox IANA (Popover+Command); Intl.supportedValuesOf, fallback COMMON_TIMEZONES (F4.2)
│   └── layout/           # AppShell (+ <Toaster/>) + Sidebar (desktop-only) + Topbar (hamburger+Sheet) + SidebarNav + ThemeToggle (F4.4/F4.7)
├── features/
│   ├── auth/
│   │   ├── components/   # LoginForm (+ test), RegisterForm, ProfileForm (+ test) (F4.2)
│   │   ├── queries/      # authKeys, useMe, useLogin, useRegister, useUpdateProfile (profile.queries.ts, F4.1/F4.3)
│   │   └── schemas/      # loginSchema / registerSchema / profileSchema (zod, cocok validasi backend)
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
│   ├── auth.store.test.ts
│   ├── ui.store.ts       # zustand + persist key 'questday-ui' — theme light/dark/system; resolveTheme/applyTheme (F4.4)
│   └── ui.store.test.ts  # (Batch 4)
├── test/setup.ts         # @testing-library/jest-dom/vitest + polyfill window.matchMedia
├── main.tsx              # enableMocking() → import '@/lib/session' + '@/stores/ui.store' → QueryClientProvider → <ErrorBoundary><RouterProvider> → Devtools(DEV)
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

### Phase 4 — SELESAI (F4.1–F4.10) → MVP FRONTEND DONE

- **Dark mode** (F4.4): `stores/ui.store.ts` (zustand persist `questday-ui`,
  `theme: light|dark|system`), `resolveTheme` / `applyTheme` toggle class `dark`
  di `<html>`; listener `matchMedia` module-level untuk mode `system`. Toggle
  `components/layout/ThemeToggle.tsx` (DropdownMenu) di Topbar. **Anti-flash**:
  IIFE inline di `<head>` `index.html` baca `localStorage` `questday-ui` sebelum
  render pertama (bertahan di `dist/index.html`).
- **Toast seragam** (F4.5): `lib/toast.ts` — `toastSuccess(msg)` +
  `toastApiError(err, fallback)` (`ApiError.message`; `status===0` → pesan
  jaringan). Sukses di-wire pada aksi yang efeknya tak kasat mata (create/edit/
  archive quest, simpan profil). Toggle complete/uncomplete TANPA toast (F2.6).
- **ErrorBoundary + 404 + RouteError** (F4.6): `components/ErrorBoundary.tsx`
  (class) membungkus `<RouterProvider>` di `main.tsx`; `pages/NotFoundPage.tsx`
  route `*`; `routes/RouteError.tsx` `errorElement` kedua grup. Detail error
  hanya `import.meta.env.DEV`.
- **Responsive** (F4.7): `Sidebar` `hidden md:flex` (desktop-only), `Topbar`
  hamburger → `Sheet` drawer (`SidebarNav` dipakai bersama). Tabel dibungkus
  `overflow-x-auto`. Tak ada scroll horizontal di 375px.
- **A11y ringan** (F4.8): tiap input `<Label htmlFor>`, tombol ikon `aria-label`,
  focus ring tak dihapus (semua `outline-none` berpasangan `focus-visible:ring`).
- **Settings + PATCH /me** (F4.1–F4.3): `authApi.updateMe → api.patch<AuthResponse>`,
  `useUpdateProfile()` `onSuccess({token,user})` → `setSession` (token BARU,
  ADR-022/013) + invalidate `authKeys.me()` / `questKeys.today()` /
  `scoringKeys.all()`. `ProfileForm` (RHF + `profileSchema`), `SettingsPage` Card
  "Profil" (email read-only). Mock `PATCH /me` di `handlers.ts`.
- **Route code-splitting** (Batch 4): 6 halaman `React.lazy` +
  `<Suspense fallback={<PageSkeleton/>}>` sekali di `GuestRoute` dan di
  `ProtectedRoute` (di dalam `AppShell`). Bundle: entry ~693 KB → 365 KB + ~15
  chunk terpisah.
- **F4.9 = SELESAI** (smoke curl lawan backend Go asli, 2026-09-01). Semua alur
  MVP cocok mock/kontrak kecuali `error.code: "unauthorized"` (401 middleware) —
  diperbaiki di `contracts/openapi.yaml` (deskripsi `ErrorResponse`) + MSW
  `GET /me`. Nol dampak fungsional (`client.ts` handle 401 by-status). Backend
  dimatikan setelah smoke. **MVP FRONTEND DONE.**

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
