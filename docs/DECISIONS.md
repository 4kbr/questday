# Catatan Keputusan (Decision Log)

Setiap keputusan desain yang tak-sepele dicatat di sini, gaya ADR ringan.
Tujuannya: "kenapa dulu kita memilih ini" tidak hilang. Tambah entri baru di
paling bawah; jangan hapus yang lama — kalau berubah, tulis entri baru yang
men-*supersede*.

Format tiap entri: **Konteks → Keputusan → Konsekuensi**. Status: Diterima /
Diganti (oleh #N) / Ditinjau ulang.

---

## ADR-001 — Monorepo dengan `apps/`
**Status:** Diterima
**Konteks:** Ada backend (sekarang) dan frontend (nanti). Ingin satu repo.
**Keputusan:** Struktur monorepo. Tiap app di `apps/<nama>` dan punya
`.gitignore` sendiri. Root `.gitignore` minimal (hanya hal lintas-app).
**Konsekuensi:** Docker compose & contract API berada di root (bersama), bukan
di dalam satu app.

## ADR-002 — Contract API di root (`contracts/`)
**Status:** Diterima
**Konteks:** Kontrak API dipakai backend (implement) dan frontend (consume).
**Keputusan:** Taruh di `contracts/` (root), OpenAPI 3.1, sebagai sumber
kebenaran bersama. Bukan di `apps/backend`.
**Konsekuensi:** Tak ada app yang "memiliki" kontrak; frontend tak jadi seolah
bergantung ke backend. Alur: ubah kontrak dulu, baru implement.

## ADR-003 — Backend: modular monolith (Go + chi)
**Status:** Diterima
**Konteks:** MVP solo, ingin sederhana tapi rapi & bisa berkembang.
**Keputusan:** Satu binary, module berbatas tegas di `internal/modules/*`,
kode teknis di `internal/platform/*`, perakitan di `internal/server`. Router
pakai chi (ringan, idiomatik). Go native untuk sisanya (net/http, database/sql).
**Konsekuensi:** Satu deploy, transaksi DB mudah. Batas module yang rapi menjaga
opsi memecah ke service terpisah nanti.

## ADR-004 — Pisahkan Quest (definisi) dari QuestLog (instance harian)
**Status:** Diterima
**Konteks:** "Quest per hari" — satu quest berulang tiap hari; perlu tahu apa
yang selesai di tanggal tertentu, dan menghitung streak.
**Keputusan:** Dua entitas: `Quest` (template: judul, poin, kategori, recurring)
dan `QuestLog` (penyelesaian pada tanggal tertentu). Semua perhitungan poin &
streak berbasis QuestLog. `UNIQUE(quest_id, date)` mencegah dobel.
**Konsekuensi:** Model lebih bersih; query "hari ini selesai apa" dan streak jadi
natural. Sedikit lebih banyak tabel/join — worth it.

## ADR-005 — Keterhubungan gamifikasi lewat port, bukan import langsung
**Status:** Diterima
**Konteks:** Menyelesaikan quest harus menambah poin. Tapi `quest` tak boleh
kawin-paksa dengan `scoring`.
**Keputusan:** `quest` mendefinisikan interface `ScoreAwarder` (port keluar).
`scoring.service` mengimplementasikannya. `server` menyuntik saat perakitan.
Untuk MVP, orkestrasi langsung (bukan event bus).
**Konsekuensi:** Module tetap decoupled dengan kompleksitas minimal. Titik
ekstensi ke event bus tetap terbuka (lihat ADR-009).

## ADR-006 — Batas hari mengikuti timezone user
**Status:** Diterima
**Konteks:** User GMT+7 menyelesaikan quest jam 23:00 tak boleh tercatat di hari
yang salah. Server bisa beda timezone.
**Keputusan:** Simpan `users.timezone` (IANA, default `Asia/Jakarta`). Handler
mengubah "sekarang" ke tanggal lokal user sebelum memanggil service. Kolom
tanggal di `quest_logs` bertipe DATE (tanggal lokal).
**Konsekuensi:** Perhitungan hari & streak benar lintas timezone. Semua kode yang
menyentuh "hari ini" WAJIB lewat timezone user, tak boleh pakai `time.Now()`
mentah.

## ADR-007 — Poin berbasis tingkat kesulitan
**Status:** Diterima (bisa ditinjau)
**Konteks:** Perlu aturan poin sederhana tapi terasa adil.
**Keputusan:** Poin ditentukan `Difficulty` (easy/medium/hard) lewat
`Quest.Points()`. Angka poin hanya didefinisikan di satu tempat (domain quest).
XP = poin untuk MVP; level = fungsi XP (`scoring.LevelForXP`).
**Konsekuensi:** Mudah dipahami. Kalau nanti mau poin dinamis (mis. bonus
streak), ubah di satu titik.

## ADR-008 — Kebijakan streak: reset saat bolong (belum ada freeze)
**Status:** Diterima (MVP)
**Konteks:** Streak bikin habit app terasa "hidup", tapi bisa menyebalkan kalau
kejam.
**Keputusan:** MVP: hari berurutan +1, hari sama tak berubah, ada hari bolong →
reset ke 1. "Freeze"/grace day = fitur nanti. Logika terpusat di
`scoring.NextStreak`.
**Konsekuensi:** Sederhana untuk MVP. `Longest` tetap disimpan agar pencapaian
tak hilang saat reset.

## ADR-009 — Uncomplete: rollback poin, streak dibiarkan (MVP)
**Status:** Diterima (MVP)
**Konteks:** Membatalkan penyelesaian harus mengembalikan poin. Tapi memutar
balik streak dengan benar itu rumit (butuh rekonstruksi urutan hari).
**Keputusan:** MVP: `OnQuestUncompleted` mengurangi poin & mencatat transaksi
negatif, tetapi TIDAK memutar balik streak. Ditinjau ulang bila terasa aneh.
**Konsekuensi:** Ada kemungkinan streak sedikit "murah hati" pada kasus batal.
Dapat diterima untuk MVP; catat sebagai utang yang disadari.

## ADR-010 — Achievement ditunda ke v2
**Status:** Diterima
**Konteks:** Fokus MVP = user + quest + scoring.
**Keputusan:** Module `achievement` hanya di-scaffold, tak di-mount. Digarap
setelah MVP. Saat digarap, pertimbangkan pindah keterhubungan ke event bus agar
scoring & achievement tak menumpuk di orkestrasi `quest`.
**Konsekuensi:** Struktur lengkap tanpa membebani MVP.

## ADR-011 — Module path Go: `questday`
**Status:** Diterima
**Konteks:** Scaffold memakai placeholder `github.com/yourorg/questday`. Semua
import internal mengikuti module path, jadi ini harus diputuskan sebelum satu
baris kode pun ditulis.
**Keputusan:** `module questday`. Import jadi `questday/internal/modules/quest`,
dst. Path pendek tanpa domain, karena modul ini tidak dipublish untuk di-import
repo lain.
**Konsekuensi:** Kalau suatu saat backend perlu di-import dari luar (mis. dipecah
jadi service terpisah yang berbagi package), path harus diubah dan seluruh import
ikut berubah. Murah selama masih monolith.

## ADR-012 — Akses database: `database/sql` + driver `pgx/v5/stdlib`
**Status:** Diterima
**Konteks:** Perlu memilih cara bicara ke Postgres sebelum menulis repository.
Pilihannya `database/sql` (dengan pgx stdlib atau lib/pq) versus `pgxpool` native.
**Keputusan:** `database/sql` dengan `_ "github.com/jackc/pgx/v5/stdlib"`. Semua
repository menerima `*sql.DB`, sesuai signature yang sudah dicontohkan scaffold
(`New(db *sql.DB, ...)`). `pgxpool` ditolak; `lib/pq` ditolak (maintenance mode).
**Konsekuensi:** Transaksi lintas-module gampang (`db.BeginTx`) — penting untuk
"buat log + tambah poin" yang atomik. Interface repository gampang di-fake untuk
test. Biayanya: fitur khusus Postgres di pgx (COPY, listen/notify, tipe kaya)
tidak langsung tersedia; kalau nanti dibutuhkan, bisa dibuka lewat `stdlib`
escape hatch atau ADR baru.

## ADR-013 — Timezone user dibawa lewat JWT claims
**Status:** Diterima
**Konteks:** ADR-006 mewajibkan "hari ini" dihitung dari timezone user. Tapi
scaffold tidak menyediakan jalan untuk membawa timezone itu ke handler `quest` —
middleware hanya menaruh `userID` ke context. Tiga opsi: masukkan ke JWT claims,
port `TimezoneProvider` dari quest ke user, atau middleware query DB.
**Keputusan:** Timezone ikut di JWT claims (`Claims{UserID, Timezone}`).
`middleware.Authenticator` mengisi keduanya ke request context; handler quest
membacanya lewat `middleware.TimezoneFrom`. Opsi middleware-query-DB ditolak
karena melanggar aturan `platform/*` tak boleh kenal `modules/*`.
**Konsekuensi:** Nol query tambahan per request. Tapi user yang mengubah
timezone-nya baru merasakan efeknya setelah token baru terbit — bisa diatasi
dengan menerbitkan ulang token saat profil diubah, atau TTL yang tidak terlalu
panjang. Kalau kelak timezone harus selalu fresh, pindah ke port
`TimezoneProvider` dan tulis ADR yang men-supersede ini.

## ADR-014 — Nama di leaderboard lewat port `UserDirectory`, bukan JOIN
**Status:** Diterima
**Konteks:** `LeaderboardEntry` memuat `DisplayName`, yang tinggal di tabel
`users`. Scaffold menaruh `Leaderboard()` di repository `scoring` — kalau
diimplementasi apa adanya, SQL scoring harus JOIN ke `users` dan jadi tahu skema
module lain.
**Keputusan:** `scoring.Repository.Leaderboard(limit)` hanya mengembalikan
`userID` + poin. `scoring` mendefinisikan port `UserDirectory` dengan
`NamesByIDs(ctx, ids) (map[string]string, error)`; `user` mengimplementasinya dan
mengeksposnya lewat `AsUserDirectory()`; `server` menjahit keduanya. Sesuai pola
ADR-005 (interface milik peminta).
**Konsekuensi:** Batas module tetap utuh — `scoring` tidak menyentuh tabel
`users`. Biayanya satu query tambahan per permintaan leaderboard, dan urutan
instansiasi di `server` jadi mengikat (`user` sebelum `scoring`). Kalau
leaderboard tumbuh besar dan dua query jadi mahal, pertimbangkan denormalisasi
nama atau cache — bukan JOIN.

## ADR-015 — Primary key: UUID v7 digenerate di aplikasi
**Status:** Diterima
**Konteks:** Tipe id belum diputuskan di scaffold. Pilihan: BIGSERIAL, UUID via
`gen_random_uuid()` (pgcrypto), atau UUID digenerate di Go.
**Keputusan:** Kolom `uuid`, nilainya dibuat di aplikasi dengan UUID v7
(`github.com/google/uuid`). Tidak ada `DEFAULT` di DDL.
**Konsekuensi:** Tak perlu ekstensi Postgres. ID sudah ada sebelum `INSERT`, jadi
tak perlu `RETURNING id` dan enak dipakai di dalam transaksi yang menulis ke
beberapa tabel sekaligus. UUID v7 terurut waktu, jadi indeks tidak terfragmentasi
seperti v4. ID tidak bisa ditebak — aman untuk dipakai di URL. Biaya: 16 byte per
id, lebih besar dari BIGSERIAL.

---

<!--
TEMPLATE ENTRI BARU (salin saat menambah keputusan):

## ADR-0NN — Judul singkat
**Status:** Diterima
**Konteks:** Masalah/situasi yang memicu keputusan.
**Keputusan:** Apa yang dipilih (dan alternatif yang ditolak, kalau relevan).
**Konsekuensi:** Untung-rugi & implikasi lanjutan.
-->
