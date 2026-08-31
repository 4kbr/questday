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

---

<!--
TEMPLATE ENTRI BARU (salin saat menambah keputusan):

## ADR-0NN — Judul singkat
**Status:** Diterima
**Konteks:** Masalah/situasi yang memicu keputusan.
**Keputusan:** Apa yang dipilih (dan alternatif yang ditolak, kalau relevan).
**Konsekuensi:** Untung-rugi & implikasi lanjutan.
-->
