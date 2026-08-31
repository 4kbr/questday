# Phase 2 — Quest

**Tujuan:** inti aplikasi — user melihat quest hari ini, mencentangnya, dan
mengelola definisi quest. Setelah phase ini QuestDay sudah berguna (walau poin
belum tampil; itu Phase 3).

**Prasyarat:** Phase 1 selesai (butuh shell, token, route terproteksi). Backend
Phase 2 jalan — atau `VITE_USE_MOCK=true`.

**Endpoint yang dipakai:** `GET /quests`, `POST /quests`, `PATCH /quests/{id}`,
`DELETE /quests/{id}`, `GET /quests/today`, `POST /quests/{id}/complete`,
`POST /quests/{id}/uncomplete`.

**Ini phase terberat di frontend** — 12 task, dan di sinilah optimistic update
hidup.

---

## F2.1 — `apis/quest.api.ts`

- **Sentuh (baru):** `src/apis/quest.api.ts`
- **Isi:**
  ```ts
  export const questApi = {
    list:       ()                                  => api.get<Quest[]>('/quests'),
    today:      ()                                  => api.get<TodayQuests>('/quests/today'),
    create:     (body: CreateQuestRequest)          => api.post<Quest>('/quests', body),
    update:     (id: string, body: UpdateQuestRequest) => api.patch<Quest>(`/quests/${id}`, body),
    archive:    (id: string)                        => api.delete<void>(`/quests/${id}`),
    complete:   (id: string)                        => api.post<void>(`/quests/${id}/complete`),
    uncomplete: (id: string)                        => api.post<void>(`/quests/${id}/uncomplete`),
  }
  ```
- **Aturan:** murni HTTP, nol logika. Semua type dari `apis/types.ts`.
- **Catatan:** `DELETE /quests/{id}` di backend adalah **arsip (soft delete)**,
  bukan hapus permanen — namai `archive`, bukan `delete`, supaya UI tak
  menjanjikan hal yang salah.
- **DoD:** 7 fungsi terketik.
- **Verifikasi:** `npm run typecheck`

## F2.2 — `questKeys`: query key terpusat

- **Sentuh (baru):** `src/features/quest/queries/keys.ts`
- **Isi:**
  ```ts
  export const questKeys = {
    all:   ()          => ['quests'] as const,
    list:  ()          => [...questKeys.all(), 'list'] as const,
    today: ()          => [...questKeys.all(), 'today'] as const,
    detail:(id: string)=> [...questKeys.all(), 'detail', id] as const,
  }
  ```
- **Kenapa task tersendiri:** key yang diketik langsung sebagai array literal di
  banyak file adalah sumber bug invalidasi paling umum — satu salah ketik dan
  data diam-diam basi. Phase 3 (F3.6) bergantung penuh pada key ini.
- **Aturan:** **dilarang** menulis `['quests', 'today']` di luar file ini.
- **DoD:** semua hook & invalidasi memakai `questKeys`.
- **Verifikasi:** `grep -rn "'quests'" src/ | grep -v keys.ts` → kosong.

## F2.3 — `features/quest/queries`: 7 hook

- **Sentuh (baru):** `src/features/quest/queries/quest.queries.ts`
- **Isi:** `useTodayQuests`, `useQuests`, `useCreateQuest`, `useUpdateQuest`,
  `useArchiveQuest`, `useCompleteQuest`, `useUncompleteQuest`.
  Mutation `onSuccess` → `invalidateQueries({ queryKey: questKeys.all() })`.
- **Aturan:**
  - Hanya file ini yang memanggil `questApi`.
  - Invalidasi **ke scoring belum ditambahkan di sini** — itu F3.6, setelah
    `scoringKeys` ada. Beri komentar `// TODO(F3.6): invalidate scoring` supaya
    tidak terlupa.
- **DoD:** 7 hook; data segar setelah tiap mutation.
- **Verifikasi:** buat quest → daftar langsung berubah tanpa reload.

## F2.4 — `features/quest/schemas`: zod

- **Sentuh (baru):** `src/features/quest/schemas/quest.schema.ts`
- **Isi:**
  ```ts
  export const questFormSchema = z.object({
    title: z.string().min(1, 'Judul wajib diisi'),
    note: z.string().optional(),
    category: z.string().min(1),
    difficulty: z.enum(['easy', 'medium', 'hard']),
  })
  ```
- **Aturan:** `difficulty` **wajib** `z.enum` yang persis sama dengan enum
  backend (`quest/domain.go`, backend T2.1) — kalau melenceng, backend menolak
  400 setelah form terlihat valid.
- **Catatan:** kategori di MVP bebas teks; kalau nanti jadi enum tetap di
  backend, ganti jadi `Select` dan turunkan pilihannya dari kontrak.
- **DoD:** satu schema dipakai form create & edit.
- **Verifikasi:** `npm run typecheck`

## F2.5 — `QuestItem`: baris quest + checkbox

- **Sentuh (baru):** `src/features/quest/components/QuestItem.tsx`
- **Isi:** satu baris: `Checkbox` (status selesai) + judul + `Badge` kategori +
  `Badge` difficulty berwarna (easy hijau / medium kuning / hard merah) + poin.
  Judul dicoret & diredupkan saat selesai.
- **Aturan:** komponen ini **presentational** — terima `quest`, `completed`, dan
  `onToggle` sebagai props. Tak memanggil hook mutation sendiri, supaya bisa
  dipakai ulang & gampang dites.
- **DoD:** tampil benar untuk kedua status, di light & dark.
- **Verifikasi:** render manual dengan data dummy.

## F2.6 — Optimistic update complete/uncomplete + rollback

- **Sentuh:** `src/features/quest/queries/quest.queries.ts`
- **Isi:** pada `useCompleteQuest` / `useUncompleteQuest`:
  ```ts
  onMutate: async (id) => {
    await queryClient.cancelQueries({ queryKey: questKeys.today() })
    const prev = queryClient.getQueryData(questKeys.today())
    queryClient.setQueryData(questKeys.today(), draft => /* balik flag completed */)
    return { prev }
  },
  onError: (_e, _id, ctx) => queryClient.setQueryData(questKeys.today(), ctx.prev),
  onSettled: () => queryClient.invalidateQueries({ queryKey: questKeys.today() }),
  ```
- **Kenapa:** mencentang quest harus terasa instan. Tanpa optimistic update ada
  jeda terlihat, dan user menekan dua kali.
- **Aturan keras:**
  - `cancelQueries` sebelum menulis cache — kalau tidak, refetch yang sedang
    jalan bisa menimpa perubahan optimistik dan checkbox "melompat balik".
  - `onError` **wajib** mengembalikan `ctx.prev` — UI tak boleh tertinggal di
    keadaan bohong saat request gagal.
  - **Kasus khusus 409 `already_completed`:** ini bukan kegagalan sungguhan
    (mis. dua tab, atau double-click). Perlakukan sebagai sukses: jangan tampilkan
    toast merah, cukup `invalidateQueries` supaya UI ikut kebenaran server.
    Kode error dibaca dari `ApiError.code` (F0.6), bukan dari teks pesan.
- **DoD:** centang terasa instan; gagal jaringan → kembali ke keadaan semula;
  409 tak memunculkan error merah.
- **Verifikasi:** matikan backend/mock di tengah aksi, perhatikan rollback.

## F2.7 — Halaman Dashboard

- **Sentuh (baru):** `src/pages/DashboardPage.tsx`,
  `src/features/quest/components/TodayQuestList.tsx`
- **Isi:** judul + tanggal hari ini (**dari field `date` milik
  `TodayQuestsResponse`**, bukan `new Date()`), daftar `QuestItem`, ringkasan
  kecil "3 dari 5 selesai" + progress bar. Slot kosong di atas untuk kartu score
  & streak (diisi F3.5).
- **Aturan keras:** tanggal yang ditampilkan **wajib** berasal dari response
  backend. Backend menghitungnya dari timezone user (ADR-006); `new Date()` di
  browser bisa beda hari dan membuat user bingung kenapa quest-nya "reset".
- **DoD:** daftar quest hari ini tampil; centang bekerja.
- **Verifikasi:** ubah timezone OS, pastikan tanggal tetap ikut backend.

## F2.8 — `QuestFormDialog`: create & edit

- **Sentuh (baru):** `src/features/quest/components/QuestFormDialog.tsx`
- **Isi:** satu `Dialog` untuk dua mode — `mode: 'create' | 'edit'`. shadcn
  `Form` + `zodResolver(questFormSchema)`. Mode edit mengisi nilai awal dari
  quest yang dipilih.
- **Aturan:**
  - Satu komponen untuk dua mode, bukan dua komponen kembar — field & validasinya
    identik.
  - Mode edit mengirim **hanya field yang berubah** (backend
    `UpdateQuestRequest` memakai pointer untuk partial update — backend T2.2).
  - Tutup dialog hanya setelah mutation sukses; saat gagal dialog tetap terbuka
    dengan pesan error.
- **DoD:** buat & edit quest jalan dari satu komponen.
- **Verifikasi:** buat quest baru, lalu edit judulnya.

## F2.9 — Halaman Quests

- **Sentuh (baru):** `src/pages/QuestsPage.tsx`,
  `src/features/quest/components/QuestTable.tsx`
- **Isi:** daftar **semua** definisi quest (bukan hanya hari ini) — `Table`
  shadcn dengan kolom judul, kategori, difficulty, poin, status aktif, dan menu
  aksi (Edit / Arsipkan). Tombol "Quest Baru" di kanan atas membuka
  `QuestFormDialog`.
- **Aturan:** Arsipkan **wajib** lewat `AlertDialog` konfirmasi, dan teksnya
  menyebut "arsipkan" — bukan "hapus permanen", karena backend melakukan soft
  delete agar riwayat & streak tetap utuh (backend T2.4).
- **DoD:** CRUD lengkap dari satu halaman.
- **Verifikasi:** buat → edit → arsip; quest terarsip tak muncul lagi di
  dashboard.

## F2.10 — Empty / loading / error state

- **Sentuh:** semua komponen daftar di phase ini
- **Isi:**
  - **Loading:** `Skeleton` shadcn berbentuk mirip konten aslinya — bukan spinner
    di tengah layar kosong.
  - **Empty (belum punya quest):** ilustrasi/ikon + kalimat ramah + tombol
    "Buat quest pertama" yang membuka dialog. Ini layar pertama yang dilihat user
    baru — jangan biarkan kosong melompong.
  - **Empty (semua quest hari ini selesai):** pesan perayaan, beda dari empty di atas.
  - **Error:** pesan + tombol "Coba lagi" yang memanggil `refetch()`.
- **Aturan:** tak ada halaman yang boleh menampilkan layar putih atau `undefined`.
- **DoD:** empat keadaan itu ada di dashboard & halaman quests.
- **Verifikasi:** matikan mock/backend, dan uji dengan akun kosong.

## F2.11 — MSW handler quest

- **Sentuh:** `src/mocks/handlers.ts`
- **Isi:** 7 endpoint quest dengan state in-memory (array quest + set log per
  tanggal) supaya complete/uncomplete benar-benar mengubah `GET /quests/today`.
  Sertakan jalur gagal: complete dua kali → 409 `already_completed`; id asing →
  404 `quest_not_found`.
- **Aturan:** mock harus **stateful** — kalau selalu mengembalikan data statis,
  optimistic update (F2.6) tak pernah benar-benar teruji.
- **DoD:** alur penuh bisa dicoba tanpa backend.
- **Verifikasi:** `VITE_USE_MOCK=true npm run dev`

## F2.12 — Test phase 2

- **Sentuh (baru):** `src/features/quest/components/QuestItem.test.tsx`,
  `src/features/quest/queries/quest.queries.test.tsx`
- **Isi:**
  - `QuestItem`: render kedua status; klik checkbox memanggil `onToggle` sekali.
  - `useCompleteQuest`: cache `today` berubah **sebelum** request selesai
    (optimistic); request gagal → cache kembali ke nilai semula; respons 409
    tidak diperlakukan sebagai error.
  - `QuestFormDialog`: judul kosong → pesan validasi, submit tak terpanggil.
- **DoD:** `npm run test` hijau.
- **Verifikasi:** `npm run test`

---

## Exit criteria Phase 2

- [ ] `npm run lint && npm run typecheck && npm run test && npm run build` bersih.
- [ ] Buat quest → muncul di halaman Quests **dan** di Dashboard hari ini.
- [ ] Centang quest terasa instan; gagal jaringan → checkbox kembali sendiri.
- [ ] Complete dua kali (mis. dua tab) tidak memunculkan error merah.
- [ ] Edit & arsip jalan; arsip pakai konfirmasi dan bahasanya "arsipkan".
- [ ] Tanggal di dashboard berasal dari response backend, bukan `new Date()`.
- [ ] Loading, empty, dan error state ada di kedua halaman.
- [ ] Tak ada query key yang diketik literal di luar `keys.ts`.
