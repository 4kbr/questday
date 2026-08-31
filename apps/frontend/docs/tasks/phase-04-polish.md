# Phase 4 — Settings & polish

**Tujuan:** melengkapi yang tersisa dan merapikan sampai layak dipakai orang
lain: pengaturan profil, dark mode, penanganan error yang manusiawi, responsive,
dan build produksi. Setelah phase ini → **MVP FRONTEND DONE**.

**Prasyarat:** Phase 3 selesai. **F4.1–F4.3 butuh backend T1.11 (`PATCH /me`)** —
kalau belum ada, kerjakan F4.4 ke bawah dulu.

---

## F4.1 — `updateMe` + `useUpdateProfile`

- **Sentuh:** `src/apis/auth.api.ts`;
  **baru:** `src/features/auth/queries/profile.queries.ts`
- **Isi:**
  ```ts
  // auth.api.ts
  updateMe: (body: UpdateProfileRequest) => api.patch<AuthResponse>('/me', body)
  // queries
  export function useUpdateProfile()
  ```
- **Aturan penting:** `PATCH /me` mengembalikan **`AuthResponse`** (token + user),
  bukan hanya user — karena timezone ikut di JWT claims (ADR-013), mengubah
  timezone **wajib** menerbitkan token baru (ADR-022). Kalau backend hanya
  mengembalikan user, itu bug backend, bukan sesuatu yang diakali frontend.
- **Blocker:** butuh backend T1.11 + kontrak `PATCH /me` (backend T0.15).
- **DoD:** fungsi & hook ada, terketik dari schema.
- **Verifikasi:** `npm run typecheck`

## F4.2 — Halaman Settings + pemilih timezone

- **Sentuh (baru):** `src/pages/SettingsPage.tsx`,
  `src/features/auth/components/ProfileForm.tsx`,
  `src/components/TimezoneSelect.tsx`
- **Isi:** `Card` "Profil": email (read-only, backend tak mengizinkan ganti
  email di MVP), `display_name` (editable), timezone (editable). Tombol simpan
  nonaktif selama tak ada perubahan.
  `TimezoneSelect`: daftar IANA dari `Intl.supportedValuesOf('timeZone')` dengan
  pencarian; fallback daftar pendek untuk browser lama.
- **Aturan:**
  - Beri catatan kecil di bawah field timezone: mengubahnya mengubah **batas
    hari** untuk quest & streak (ADR-006). Ini bukan preferensi tampilan.
  - Komponen `TimezoneSelect` ditaruh di `components/` (bukan `features/auth/`)
    karena juga dipakai form register (F1.6) — kalau F1.6 sudah menulis versi
    sendiri, satukan di sini.
- **DoD:** profil bisa diubah dan tersimpan.
- **Verifikasi:** ubah nama → reload → nama baru bertahan.

## F4.3 — Simpan token baru setelah timezone berubah

- **Sentuh:** `src/features/auth/queries/profile.queries.ts`,
  `src/stores/auth.store.ts`
- **Isi:**
  ```ts
  onSuccess: ({ token, user }) => {
    useAuthStore.getState().setSession(token, user)   // token BARU
    queryClient.invalidateQueries({ queryKey: questKeys.today() })
    queryClient.invalidateQueries({ queryKey: scoringKeys.all() })
  }
  ```
- **Aturan keras:** token lama masih memuat timezone lama. Kalau tidak diganti,
  backend akan terus menghitung "hari ini" dengan timezone lama sampai token
  kedaluwarsa (hingga 24 jam) — user mengubah setelan tapi tak terjadi apa-apa,
  gejala yang sangat membingungkan untuk dilacak.
- **Kenapa invalidate quest & scoring:** batas hari berubah, jadi "quest hari ini"
  dan streak bisa berbeda.
- **DoD:** setelah ganti timezone, `/quests/today` mengembalikan tanggal menurut
  timezone baru.
- **Verifikasi:** ubah timezone ke zona yang berbeda tanggalnya, lihat dashboard.

## F4.4 — Dark mode

- **Sentuh (baru):** `src/stores/ui.store.ts`,
  `src/components/layout/ThemeToggle.tsx`; **Ubah:** `src/main.tsx`
- **Isi:** zustand persist `theme: 'light' | 'dark' | 'system'`; terapkan dengan
  menambah/menghapus class `dark` di `<html>`; mode `system` mengikuti
  `prefers-color-scheme`. Toggle di topbar.
- **Aturan:** terapkan tema **sebelum** render pertama (skrip kecil di
  `index.html` atau di awal `main.tsx`) supaya tak ada kedipan putih saat
  memuat halaman dalam mode gelap.
- **DoD:** tema bertahan setelah reload, tanpa kedipan.
- **Verifikasi:** set dark → reload → tetap gelap sejak frame pertama.

## F4.5 — Toast seragam untuk mutation

- **Sentuh:** `src/lib/`, semua file `queries`
- **Isi:** helper kecil yang memetakan hasil mutation ke `sonner` toast:
  sukses → pesan singkat; gagal → `ApiError.message` dari backend.
- **Aturan:**
  - **Jangan** menampilkan toast untuk hal yang sudah terlihat di UI (mis.
    centang quest yang sudah optimistic) — hanya untuk aksi yang efeknya tak
    langsung kasat mata (simpan profil, buat/edit/arsip quest).
  - 409 `already_completed` **tidak** memunculkan toast merah (lihat F2.6).
  - Pesan error ditampilkan dari backend, bukan dikarang frontend — kecuali
    error jaringan (backend tak terjangkau).
- **DoD:** perilaku toast konsisten di seluruh app.
- **Verifikasi:** lakukan tiap mutation, perhatikan yang muncul & tidak.

## F4.6 — ErrorBoundary + halaman 404

- **Sentuh (baru):** `src/components/ErrorBoundary.tsx`,
  `src/pages/NotFoundPage.tsx`; **Ubah:** `src/routes/index.tsx`
- **Isi:** ErrorBoundary membungkus router — tampilkan halaman "Ada yang salah" +
  tombol muat ulang, bukan layar putih. Route `*` → 404 dengan tautan kembali ke
  dashboard.
- **Aturan:** detail error hanya ditampilkan di dev (`import.meta.env.DEV`).
- **DoD:** komponen yang melempar error tak menghasilkan layar putih.
- **Verifikasi:** lempar error sengaja di satu komponen.

## F4.7 — Responsive

- **Sentuh:** `src/components/layout/*`, halaman-halaman daftar
- **Isi:** di bawah `md`, sidebar jadi `Sheet` (drawer) yang dibuka lewat tombol
  hamburger di topbar. Tabel yang tak muat berubah jadi daftar kartu, atau
  dibungkus kontainer `overflow-x-auto`.
- **Aturan:** halaman **tak boleh** bisa di-scroll horizontal di lebar 375px.
- **DoD:** semua halaman terbaca di 375px.
- **Verifikasi:** DevTools responsive, cek 375 / 768 / 1440.

## F4.8 — Audit a11y ringan

- **Sentuh:** komponen form & interaktif
- **Isi:** tiap input punya `<Label htmlFor>`; tombol ikon punya `aria-label`;
  focus ring terlihat di kedua tema; kontras teks memadai; dialog mengembalikan
  fokus saat ditutup (bawaan shadcn/Radix, pastikan tak dirusak).
- **Aturan:** jangan menghapus outline fokus demi estetika — ganti dengan ring
  yang terlihat.
- **DoD:** seluruh alur utama bisa diselesaikan **hanya dengan keyboard**.
- **Verifikasi:** Tab dari login sampai mencentang quest tanpa menyentuh mouse.

## F4.9 — Lepas mock: uji melawan backend asli

- **Sentuh:** `.env` lokal
- **Isi:** `VITE_USE_MOCK=false`, backend jalan (`make dev` + Postgres). Jalankan
  alur penuh dan **catat setiap ketidakcocokan** antara mock dan backend asli
  (bentuk envelope, nama field, kode error, status).
- **Aturan:** setiap ketidakcocokan diselesaikan dengan memperbaiki **kontrak
  atau backend**, lalu `npm run gen:api` — bukan dengan menambal di frontend.
  Kalau frontend menambal, kontrak berhenti jadi sumber kebenaran (ADR-002).
- **DoD:** seluruh alur MVP jalan melawan backend asli.
- **Verifikasi:** lihat exit criteria di bawah.

## F4.10 — Build produksi (+ container, opsional)

- **Sentuh:** `apps/frontend/Dockerfile`, `nginx.conf` (kalau dikontainerkan)
- **Isi:** `npm run build` → `dist/`, cek dengan `npm run preview`. Kalau perlu
  container: build multi-stage → nginx statis dengan fallback SPA
  (`try_files $uri /index.html`).
- **Aturan:** `VITE_*` di-*inline* saat build — jadi `VITE_API_BASE_URL`
  ditentukan saat build image, bukan saat runtime. Catat ini di README supaya tak
  ada yang mengira bisa diganti lewat env container.
- **DoD:** `npm run build` bersih; preview jalan; refresh di route dalam
  (mis. `/quests`) tidak 404.
- **Verifikasi:** `npm run build && npm run preview`, lalu reload di `/quests`.

---

## Exit criteria Phase 4 — **MVP FRONTEND DONE**

- [ ] `npm run lint && npm run typecheck && npm run test && npm run build` bersih.
- [ ] Dengan `VITE_USE_MOCK=false` dan backend asli, alur penuh sukses:
      register → login → buat quest → complete → poin & streak naik →
      leaderboard → ubah timezone → tetap login dengan token baru → logout.
- [ ] Refresh browser tetap login; 401 melempar ke login dengan bersih.
- [ ] Light & dark mode benar, tanpa kedipan saat memuat.
- [ ] Terbaca di 375px; tak ada scroll horizontal.
- [ ] Alur utama bisa diselesaikan dengan keyboard saja.
- [ ] Tak ada layar putih di keadaan apa pun (loading, kosong, error, 404).
- [ ] Tak ada `axios`/`fetch` di luar `src/apis/`:
      ```bash
      grep -rn "axios\|fetch(" src/ --include=*.ts --include=*.tsx \
        | grep -v "^src/apis/" | grep -v "^src/mocks/"
      ```
- [ ] Tak ada type API yang diketik tangan di luar `src/apis/types.ts`.
- [ ] Keputusan baru selama implementasi tercatat di `docs/DECISIONS.md`.

Setelah ini: **jangan lanjut ke fitur v2 tanpa diminta** (AGENTS.md).
Backlog v2 ada di [README.md](README.md).
