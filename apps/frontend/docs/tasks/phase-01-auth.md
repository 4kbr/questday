# Phase 1 — Auth & shell

**Tujuan:** user bisa register, login, dan masuk ke aplikasi ber-layout SaaS.
Setelah phase ini kerangka aplikasi berdiri: ada sesi, ada rute terlindungi, ada
sidebar & topbar yang dipakai semua halaman berikutnya.

**Prasyarat:** Phase 0 selesai. Backend Phase 1 (auth) sudah jalan — atau
kerjakan dengan `VITE_USE_MOCK=true`.

**Endpoint yang dipakai:** `POST /auth/register`, `POST /auth/login`, `GET /me`.

---

## F1.1 — `apis/auth.api.ts`

- **Sentuh (baru):** `src/apis/auth.api.ts`
- **Isi:**
  ```ts
  export const authApi = {
    register: (body: RegisterRequest) => api.post<AuthResponse>('/auth/register', body),
    login:    (body: LoginRequest)    => api.post<AuthResponse>('/auth/login', body),
    me:       ()                      => api.get<User>('/me'),
  }
  ```
- **Aturan:** fungsi di sini **murni HTTP** — tak tahu React, tak menyentuh
  store, tak melempar toast. Semua type dari `apis/types.ts` (F0.5).
- **DoD:** 3 fungsi, semua terketik dari schema.
- **Verifikasi:** `npm run typecheck`

## F1.2 — `stores/auth.store.ts`

- **Sentuh (baru):** `src/stores/auth.store.ts`; **Ubah:** `src/apis/client.ts`
- **Isi:**
  ```ts
  type AuthState = {
    token: string | null
    user: User | null
    setSession: (token: string, user: User) => void
    setUser: (user: User) => void          // dipakai F4.3
    logout: () => void
    isAuthenticated: () => boolean
  }
  // zustand + persist, key: 'questday-auth', hanya {token, user}
  ```
  Lalu rapikan interceptor di `client.ts` agar membaca token dari store ini
  (`useAuthStore.getState().token`) — menggantikan akses `localStorage` sementara
  di F0.6.
- **Aturan keras (ADR-018 & ADR-020):**
  - Store ini **hanya** untuk client state. **Dilarang** menyimpan quest, score,
    atau leaderboard di zustand — itu server state, milik TanStack Query.
  - `user` disimpan di sini sebagai *cache sesi* supaya layout bisa langsung
    menampilkan nama tanpa menunggu `GET /me`; sumber kebenarannya tetap query
    `useMe`.
  - Tak ada komponen yang membaca `localStorage` langsung.
- **DoD:** refresh browser → token & user masih ada.
- **Verifikasi:** login, refresh, cek Application → Local Storage.

## F1.3 — `features/auth/schemas`: zod

- **Sentuh (baru):** `src/features/auth/schemas/auth.schema.ts`
- **Isi:**
  ```ts
  export const loginSchema = z.object({
    email: z.string().email('Email tidak valid'),
    password: z.string().min(1, 'Password wajib diisi'),
  })
  export const registerSchema = z.object({
    email: z.string().email(),
    password: z.string().min(8, 'Minimal 8 karakter'),
    display_name: z.string().min(1),
    timezone: z.string().optional(),   // backend default 'Asia/Jakarta'
  })
  ```
- **Aturan keras:** aturan validasi **wajib cocok dengan backend**
  (`min=8` untuk password, email, `display_name` required, timezone opsional —
  lihat backend T1.2). Kalau frontend lebih longgar, user dapat 400 yang
  membingungkan; kalau lebih ketat, ada input sah yang ditolak diam-diam.
- **Catatan:** validasi client adalah **kenyamanan**, bukan pengaman. Error dari
  backend tetap harus ditampilkan (F1.5/F1.6).
- **DoD:** 2 schema; `z.infer` cocok dengan `LoginRequest`/`RegisterRequest`.
- **Verifikasi:** `npm run typecheck`

## F1.4 — `features/auth/queries`

- **Sentuh (baru):** `src/features/auth/queries/auth.queries.ts`
- **Isi:**
  ```ts
  export const authKeys = { me: () => ['auth', 'me'] as const }

  export function useMe()       // useQuery, enabled: !!token
  export function useLogin()    // useMutation -> onSuccess: setSession(token, user)
  export function useRegister() // useMutation -> onSuccess: setSession(...)
  ```
- **Aturan:** **hanya hook di folder ini** yang memanggil `authApi`. Komponen
  memanggil hook, tak pernah `authApi` langsung.
- **DoD:** 3 hook; login sukses mengisi store.
- **Verifikasi:** test F1.11.

## F1.5 — Halaman Login

- **Sentuh (baru):** `src/pages/LoginPage.tsx`,
  `src/features/auth/components/LoginForm.tsx`
- **Isi:** layout terpusat (card di tengah, logo/nama app di atas, tautan ke
  register di bawah). Form: shadcn `Form` + `react-hook-form` +
  `zodResolver(loginSchema)`.
- **Aturan:**
  - Tombol submit `disabled` + spinner saat `isPending` — cegah double submit.
  - Error backend ditampilkan: `ApiError.code === 'invalid_credential'` →
    pesan di atas form ("Email atau password salah"), **bukan** di field
    tertentu — backend sengaja tak memberi tahu yang mana yang salah
    (backend T1.4).
  - Sukses → redirect ke `PATHS.dashboard` (atau ke halaman yang tadi dituju,
    lihat F1.7).
- **DoD:** login salah menampilkan pesan; login benar masuk dashboard.
- **Verifikasi:** coba kedua kasus dengan MSW.

## F1.6 — Halaman Register

- **Sentuh (baru):** `src/pages/RegisterPage.tsx`,
  `src/features/auth/components/RegisterForm.tsx`
- **Isi:** sama polanya dengan Login, tambah `display_name` dan **timezone**
  (default dari browser: `Intl.DateTimeFormat().resolvedOptions().timeZone`,
  tetap bisa diubah).
- **Aturan:**
  - `ApiError.code === 'email_taken'` (409) → pesan di **field email**
    ("Email sudah terdaftar"), bukan di banner umum — di sini backend memang
    memberi tahu field-nya.
  - Timezone yang dikirim harus IANA valid (mis. `Asia/Jakarta`); pakai
    komponen pemilih yang sama dengan F4.2 kalau sudah ada, kalau belum cukup
    `Select` berisi daftar umum + default browser.
- **Kenapa timezone penting:** ini yang menentukan batas hari untuk streak
  (ADR-006). Salah di sini, quest jam 23:00 bisa tercatat di hari yang salah.
- **DoD:** register sukses langsung login (backend mengembalikan token).
- **Verifikasi:** register → langsung berada di dashboard.

## F1.7 — `routes/`: ProtectedRoute & GuestRoute

- **Sentuh:** `src/routes/index.tsx`; **baru:**
  `src/routes/ProtectedRoute.tsx`, `src/routes/GuestRoute.tsx`
- **Isi:**
  - `ProtectedRoute`: tak ada token → `<Navigate to={PATHS.login} replace
    state={{ from: location }} />`. Ada token → render `<AppShell><Outlet/></AppShell>`.
  - `GuestRoute`: sudah punya token → lempar ke dashboard (supaya user login tak
    melihat halaman login lagi).
- **Aturan:** simpan lokasi asal di `state.from` supaya setelah login user
  kembali ke halaman yang tadi dituju, bukan selalu ke dashboard.
- **DoD:** `/quests` tanpa token → `/login`; setelah login kembali ke `/quests`.
- **Verifikasi:** buka `/quests` di tab baru tanpa token.

## F1.8 — `components/layout/AppShell`

- **Sentuh (baru):** `src/components/layout/AppShell.tsx`, `Sidebar.tsx`,
  `Topbar.tsx`
- **Isi:**
  - **Sidebar** kiri: nama app + nav (Dashboard, Quests, Leaderboard, Settings)
    dengan ikon `lucide-react`, item aktif disorot berdasarkan route sekarang.
  - **Topbar**: judul halaman di kiri, kanan berisi `Avatar` + `DropdownMenu`
    (nama & email user, Settings, Logout).
  - Konten: `<Outlet />` di area utama dengan padding konsisten.
- **Aturan:** ini shell untuk **semua** halaman terproteksi — halaman berikutnya
  tak boleh menggambar sidebar/topbar sendiri. Responsive dikerjakan di F4.7;
  untuk sekarang cukup benar di desktop.
- **DoD:** shell tampil di semua route terproteksi; nav aktif benar.
- **Verifikasi:** klik tiap menu, perhatikan sorotan item aktif.

## F1.9 — Logout: bersihkan token **dan** cache Query

- **Sentuh:** `src/stores/auth.store.ts`, `src/components/layout/Topbar.tsx`
- **Isi:**
  ```ts
  // saat logout:
  useAuthStore.getState().logout()   // token & user -> null, localStorage bersih
  queryClient.clear()                // WAJIB
  navigate(PATHS.login, { replace: true })
  ```
- **Aturan keras:** `queryClient.clear()` **tidak boleh dilupakan**. Tanpa itu,
  user berikutnya yang login di browser sama akan melihat quest & skor milik user
  sebelumnya dari cache sebelum refetch selesai — kebocoran data antar-akun.
- **Catatan:** interceptor 401 di `client.ts` harus melakukan hal yang sama —
  ekstrak jadi satu fungsi `endSession()` supaya tak ada jalur yang lupa.
- **DoD:** logout → login sebagai user lain → nol data user lama muncul.
- **Verifikasi:** login A, buka dashboard, logout, login B — cek isi dashboard.

## F1.10 — MSW handler auth

- **Sentuh:** `src/mocks/handlers.ts`
- **Isi:** handler untuk `/auth/register`, `/auth/login`, `/me`. Sertakan jalur
  gagal: login salah → 401 `invalid_credential`; register email terpakai → 409
  `email_taken`.
- **Aturan:** bentuk response mengikuti envelope backend & type dari
  `apis/types.ts` (lihat F0.9). Mock jalur gagal sama pentingnya dengan jalur
  sukses — tanpa itu penanganan error di F1.5/F1.6 tak pernah teruji.
- **DoD:** semua state form bisa dicoba tanpa backend jalan.
- **Verifikasi:** `VITE_USE_MOCK=true npm run dev`

## F1.11 — Test phase 1

- **Sentuh (baru):** `src/features/auth/components/LoginForm.test.tsx`,
  `src/stores/auth.store.test.ts`
- **Isi:**
  - `LoginForm`: email tak valid → pesan validasi, submit tidak terpanggil;
    submit sukses → `setSession` terpanggil; 401 → banner error tampil.
  - `auth.store`: `setSession` mengisi, `logout` mengosongkan, persist bekerja.
  - `ProtectedRoute`: tanpa token me-redirect.
- **DoD:** `npm run test` hijau.
- **Verifikasi:** `npm run test`

---

## Exit criteria Phase 1

- [ ] `npm run lint && npm run typecheck && npm run test && npm run build` bersih.
- [ ] Register → otomatis login → dashboard.
- [ ] Login salah menampilkan pesan yang benar (401 vs 409 dibedakan).
- [ ] Refresh browser tetap login; buka route terproteksi tanpa token → `/login`,
      dan setelah login kembali ke route yang dituju.
- [ ] Logout membersihkan token **dan** `queryClient` — tak ada data user lama.
- [ ] Shell SaaS (sidebar + topbar) tampil di semua halaman terproteksi.
- [ ] Tak ada komponen yang memanggil `authApi` langsung — semua lewat hook.
