# Phase 0 — Setup & fondasi

**Tujuan:** menyiapkan project Vite yang jalan, tema shadcn ala SaaS, dan
**lapisan API** yang dipakai semua fitur. Setelah phase ini belum ada halaman
berisi, tapi semua bahan untuk membangun fitur sudah ada.

**Prasyarat:** Node 20+, npm. **Dan `contracts/openapi.yaml` sudah terisi
schema** — itu backend **T0.15** (Phase 0 backend). Kalau kontrak belum siap,
F0.4 memblokir; kerjakan F0.1–F0.3 dulu.

**Kenapa duluan:** `apis/client.ts` dan `schema.gen.ts` dipakai setiap fitur.
Salah bentuk di sini menular ke semua phase.

---

## F0.1 — Scaffold Vite + React + TypeScript

- **Sentuh (baru):** seluruh `apps/frontend/`
- **Isi:**
  ```bash
  cd apps/frontend
  npm create vite@latest . -- --template react-ts
  npm install
  ```
  Bersihkan boilerplate: hapus `src/App.css`, logo, dan isi demo `App.tsx`.
- **Aturan:** semua tinggal di `apps/frontend/` — jangan mengotori root monorepo
  dengan `node_modules` atau `package.json` (ADR-001: tiap app berdiri sendiri).
- **DoD:** `npm run dev` menyalakan halaman kosong tanpa error.
- **Verifikasi:** `npm run dev` lalu buka `http://localhost:5173`

## F0.2 — Tailwind CSS + path alias

- **Sentuh:** `tailwind.config.ts`, `src/index.css`, `vite.config.ts`,
  `tsconfig.json`
- **Isi:** pasang Tailwind (v4 lewat `@tailwindcss/vite`, atau v3 + PostCSS —
  ikuti versi yang didukung shadcn saat itu). Lalu alias `@/` → `src/`:
  ```ts
  // vite.config.ts
  resolve: { alias: { '@': path.resolve(__dirname, './src') } }
  // tsconfig.json
  "paths": { "@/*": ["./src/*"] }
  ```
- **Aturan:** alias `@/` wajib — semua import di task berikutnya memakainya.
- **DoD:** class Tailwind bekerja; `import x from '@/lib/utils'` resolve.
- **Verifikasi:** `npm run build`

## F0.3 — shadcn/ui init + tema SaaS

- **Sentuh:** `components.json`, `src/index.css`, `src/components/ui/*`,
  `src/lib/utils.ts`
- **Isi:**
  ```bash
  npx shadcn@latest init
  npx shadcn@latest add button input label card form dialog table \
    dropdown-menu avatar badge skeleton sonner separator select switch checkbox
  ```
  Tetapkan token tema sekali di sini: warna primary, `--radius`, dan **variabel
  dark mode** (`.dark`). Font: satu sans (mis. Inter) untuk seluruh app.
- **Arah desain (dipakai semua phase):** SaaS — **sidebar kiri** (nav: Dashboard,
  Quests, Leaderboard, Settings) + **topbar** (judul halaman kiri, avatar kanan),
  konten di `Card`, spacing lega, border halus, aksen warna hanya untuk aksi
  utama & badge difficulty.
- **Aturan:** `src/components/ui/*` adalah hasil generate shadcn — boleh diubah,
  tapi jangan taruh logika bisnis di situ. Komponen milik kita di
  `features/*/components/` atau `components/layout/`.
- **DoD:** `<Button>` shadcn tampil benar di light & dark.
- **Verifikasi:** render satu Button + Card di `App.tsx`, toggle class `dark`
  di `<html>`.

## F0.4 — `gen:api`: generate type dari kontrak

- **Sentuh:** `package.json` (script), `src/apis/schema.gen.ts` (hasil generate)
- **Isi:**
  ```bash
  npm i -D openapi-typescript
  ```
  ```json
  "scripts": {
    "gen:api": "openapi-typescript ../../contracts/openapi.yaml -o src/apis/schema.gen.ts"
  }
  ```
- **Aturan keras:**
  - `schema.gen.ts` **tak pernah diedit tangan** — beri header komentar yang
    menyatakan itu, dan masukkan ke `.eslintignore` / `prettierignore`.
  - Kalau type yang dibutuhkan tak ada, **perbaiki kontraknya**, bukan hasil
    generate-nya (ADR-002 & ADR-019).
- **Blocker:** butuh backend T0.15 selesai. Kalau `components.schemas` masih
  `{}`, hasil generate kosong dan F0.5 tak bisa jalan.
- **DoD:** `npm run gen:api` menghasilkan file berisi 14+ schema.
- **Verifikasi:** `grep -c "QuestResponse\|ScoreResponse" src/apis/schema.gen.ts`

## F0.5 — `apis/types.ts`: alias type ramah

- **Sentuh (baru):** `src/apis/types.ts`
- **Isi:**
  ```ts
  import type { components } from './schema.gen'
  type S = components['schemas']

  export type User            = S['UserResponse']
  export type AuthResponse    = S['AuthResponse']
  export type RegisterRequest = S['RegisterRequest']
  export type LoginRequest    = S['LoginRequest']
  export type Quest           = S['QuestResponse']
  export type CreateQuestRequest = S['CreateQuestRequest']
  export type UpdateQuestRequest = S['UpdateQuestRequest']
  export type TodayQuests     = S['TodayQuestsResponse']
  export type Score           = S['ScoreResponse']
  export type Streak          = S['StreakResponse']
  export type LeaderboardEntry = S['LeaderboardEntry']
  export type ApiErrorBody    = S['ErrorResponse']
  ```
- **Kenapa:** `components['schemas']['...']` berisik dipakai di komponen. File ini
  satu-satunya jembatan — kalau nama schema di kontrak berubah, hanya file ini
  yang menyesuaikan.
- **Aturan:** **dilarang** mendefinisikan `interface Quest { ... }` sendiri di
  mana pun. Semua type API turun dari sini.
- **DoD:** semua type MVP punya alias.
- **Verifikasi:** `npx tsc --noEmit`

## F0.6 — `apis/client.ts`: axios + interceptor

- **Sentuh (baru):** `src/apis/client.ts`
- **Isi:**
  ```ts
  export const api = axios.create({ baseURL: import.meta.env.VITE_API_BASE_URL })

  // request: pasang Authorization dari auth.store (dibuat di F1.2)
  // response sukses: unwrap envelope backend  {"data": ...}  -> kembalikan .data
  // response gagal: ubah {"error":{code,message}} jadi ApiError
  export class ApiError extends Error {
    constructor(public status: number, public code: string, message: string)
  }
  // 401 -> auth.store.logout() + redirect ke /login
  ```
- **Aturan keras:**
  - **Ini satu-satunya file yang membuat instance HTTP.** Tak ada `axios.create`
    di tempat lain, tak ada `fetch()` di komponen.
  - Unwrap envelope dilakukan **di sini**, sekali — supaya `*.api.ts` dan
    komponen tak perlu menulis `res.data.data`. Bentuk envelope ditetapkan
    backend T0.7/T0.8.
  - `ApiError.code` (mis. `already_completed`, `email_taken`) dipakai fitur untuk
    membedakan perlakuan — jangan mencocokkan string pesan.
- **Catatan:** interceptor butuh `auth.store` yang baru lahir di F1.2. Untuk
  sekarang baca token lewat fungsi yang bisa disuntik, atau langsung dari
  `localStorage` dengan key yang sama — rapikan di F1.2.
- **DoD:** satu tempat untuk baseURL, token, unwrap, dan error.
- **Verifikasi:** test di F0.12.

## F0.7 — `lib/query-client.ts` + provider

- **Sentuh (baru):** `src/lib/query-client.ts`, `src/main.tsx`
- **Isi:**
  ```ts
  new QueryClient({ defaultOptions: { queries: {
    staleTime: 30_000,
    retry: (count, err) => !(err instanceof ApiError && err.status < 500) && count < 2,
    refetchOnWindowFocus: true,
  }}})
  ```
  Bungkus app dengan `<QueryClientProvider>`. Pasang React Query Devtools hanya
  di dev.
- **Aturan:** **jangan retry error 4xx** — 401/404/409 tak akan membaik dengan
  diulang, dan retry pada 409 `already_completed` bikin UI berkedip.
- **DoD:** QueryClient tunggal, dipakai seluruh app.
- **Verifikasi:** Devtools muncul di `npm run dev`.

## F0.8 — Struktur folder + kerangka router

- **Sentuh (baru):** `src/routes/paths.ts`, `src/routes/index.tsx`,
  folder-folder kosong sesuai struktur di README
- **Isi:**
  ```ts
  // paths.ts — konstanta, bukan string tersebar
  export const PATHS = {
    login: '/login', register: '/register',
    dashboard: '/', quests: '/quests',
    leaderboard: '/leaderboard', settings: '/settings',
  } as const
  ```
  Router dasar dengan placeholder page tiap route (isi menyusul per phase).
- **Aturan:** **dilarang** menulis path sebagai string literal di komponen
  (`<Link to="/quests">`) — selalu `PATHS.quests`. Ini yang membuat rename route
  aman.
- **DoD:** 6 route bisa dikunjungi (walau isinya placeholder).
- **Verifikasi:** klik-klik semua route di browser.

## F0.9 — MSW: mock backend untuk dev

- **Sentuh (baru):** `src/mocks/handlers.ts`, `src/mocks/browser.ts`,
  `src/main.tsx`
- **Isi:**
  ```ts
  // main.tsx
  if (import.meta.env.VITE_USE_MOCK === 'true') {
    const { worker } = await import('./mocks/browser')
    await worker.start({ onUnhandledRequest: 'warn' })
  }
  ```
  Handler dasar dulu: `GET /healthz`. Handler per fitur ditambah di F1.10, F2.11,
  F3.8.
- **Aturan keras:**
  - Bentuk response mock **wajib mengikuti `schema.gen.ts`** — ketik handler-nya
    dengan type dari `apis/types.ts` supaya mock yang melenceng ketahuan saat
    compile. Mock yang bohong lebih berbahaya daripada tak ada mock.
  - Mock juga harus meniru **envelope** `{"data": ...}` dan bentuk error
    `{"error":{code,message}}`, bukan mengembalikan objek telanjang.
  - MSW **tak pernah aktif di production build** — hanya di balik flag.
- **DoD:** `VITE_USE_MOCK=true npm run dev` → request dijawab mock;
  `false` → request menembus ke backend.
- **Verifikasi:** cek tab Network — request ter-intercept.

## F0.10 — Config project: env, gitignore, lint

- **Sentuh (baru):** `.env.example`, `.gitignore`, `eslint.config.js`,
  `.prettierrc`
- **Isi:**
  ```
  # .env.example
  VITE_API_BASE_URL=http://localhost:8080/api/v1
  VITE_USE_MOCK=true
  ```
  `.gitignore`: `node_modules`, `dist`, `.env`, `.env.*` (kecuali
  `.env.example`) — root `.gitignore` sengaja minimal (ADR-001).
  ESLint + Prettier: aktifkan `@typescript-eslint`, `react-hooks`, dan
  **larang import silang antar fitur** kalau bisa (mis. `no-restricted-imports`:
  `features/*/` tak boleh saling import — lewat `components/` atau `lib/` saja).
- **Aturan:** `.env` asli tak pernah masuk git; `.env.example` hanya placeholder.
- **DoD:** `npm run lint` bersih.
- **Verifikasi:** `npm run lint`

## F0.11 — Script `package.json` + README frontend

- **Sentuh:** `package.json`, `apps/frontend/README.md`
- **Isi:** script: `dev`, `build`, `preview`, `lint`, `format`, `test`,
  `gen:api`, `typecheck` (`tsc --noEmit`).
  README: cara mulai (install → `cp .env.example .env` → `npm run gen:api` →
  `npm run dev`), penjelasan flag mock, dan tautan ke `docs/tasks/`.
- **DoD:** orang baru bisa menjalankan project hanya dari README.
- **Verifikasi:** ikuti README dari nol di folder bersih.

## F0.12 — Test phase 0

- **Sentuh (baru):** `vitest.config.ts`, `src/test/setup.ts`,
  `src/apis/client.test.ts`
- **Isi:** pasang Vitest + Testing Library + jsdom. Test:
  - `client.ts` meng-unwrap envelope `{"data": ...}` dengan benar.
  - Response error `{"error":{code,message}}` jadi `ApiError` dengan `status` &
    `code` yang benar.
  - Render `App` tidak crash.
- **DoD:** `npm run test` hijau.
- **Verifikasi:** `npm run test`

---

## Exit criteria Phase 0

- [ ] `npm run dev`, `npm run build`, `npm run lint`, `npm run test`,
      `npm run typecheck` semuanya bersih.
- [ ] `npm run gen:api` menghasilkan `schema.gen.ts` berisi schema MVP.
- [ ] 6 route bisa dikunjungi (placeholder), layout SaaS dasar tampil.
- [ ] `VITE_USE_MOCK=true` → MSW menjawab; `false` → menembus ke backend.
- [ ] Light & dark mode keduanya terbaca.
- [ ] Tak ada `axios.create` atau `fetch(` di luar `src/apis/` dan `src/mocks/`.
