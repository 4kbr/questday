# AGENTS.md

Panduan untuk AI coding agent (dan kolaborator baru) yang bekerja di repo
QuestDay. Baca ini sebelum menulis kode. Tujuannya menjaga konsistensi arsitektur
dan agar keputusan tak diam-diam ditabrak.

## Orientasi dulu

Sebelum mengubah apa pun, baca:

1. `docs/ARCHITECTURE.md` — bentuk sistem & aturan dependensi.
2. `docs/DECISIONS.md` — keputusan yang sudah diambil (JANGAN dilanggar diam-diam).
3. `docs/GUIDES.md` — resep langkah-per-langkah untuk tugas umum.

Kalau tugas cocok dengan salah satu resep di GUIDES, ikuti resep itu.

## Prinsip proyek

- **Modular monolith.** Batas antar-module itu suci. Lihat aturan dependensi di
  ARCHITECTURE. Pelanggaran paling umum yang harus dihindari: satu module
  meng-import package module lain.
- **MVP dulu.** Fokus: `user`, `quest`, `scoring`. `achievement` = v2, jangan
  dikerjakan kecuali diminta eksplisit.
- **Pemilik kode ingin menulis sendiri.** Scaffold ini sengaja berisi TODO, bukan
  implementasi. JANGAN mengisi implementasi besar-besaran tanpa diminta. Kalau
  diminta, kerjakan sepotong dan jelaskan.

## Aturan keras (jangan dilanggar)

1. **Domain bersih.** `domain.go` tak mengimpor HTTP, SQL, atau package kita yang
   lain. Aturan bisnis murni diam di sini.
2. **SQL hanya di `repository_postgres.go`.** Tak ada SQL di service/handler.
3. **Handler tipis.** Handler hanya decode/validate, panggil service, tulis
   response lewat `platform/httpx`. Tak ada logika bisnis di handler.
4. **Antar-module lewat port.** Butuh module lain → definisikan interface di
   module peminta, suntik implementasi via `server`. Contoh acuan:
   `quest.ScoreAwarder` diimplementasi `scoring`.
5. **`platform/*` netral domain.** Tak boleh import `modules/*`.
6. **"Hari ini" = timezone user.** Jangan pakai `time.Now()` mentah untuk logika
   harian/streak. Selalu lewat timezone user (ADR-006).
7. **Aturan poin & level satu sumber.** Poin di `quest.Quest.Points()`, level di
   `scoring.LevelForXP`, streak di `scoring.NextStreak`. Jangan menyebar angka
   ajaib ke tempat lain.
8. **Kontrak lebih dulu.** Menambah/mengubah endpoint → ubah
   `contracts/openapi.yaml` sebelum/berbarengan dengan implementasi.
9. **Jangan bocorkan data sensitif.** `PasswordHash` dan sejenisnya tak pernah
   masuk response. Query selalu difilter `user_id`.

## Saat menambah sesuatu

Ikuti GUIDES:

- Endpoint baru → GUIDES §2. Module baru → §4. Module saling bicara → §5.
  Migrasi → §6.
- Selalu perbarui `contracts/openapi.yaml`.
- Keputusan desain baru → tambah entri ADR di `docs/DECISIONS.md` (jangan hapus
  yang lama; supersede bila berubah).

## Gaya kode

- Go idiomatik. `gofmt` wajib (`make fmt`), lolos `go vet` (`make vet`).
- Error dibungkus dengan konteks (`fmt.Errorf("...: %w", err)`); error domain
  didefinisikan di `domain.go` dan dipetakan ke HTTP di handler/httpx.
- Nama package huruf kecil, tanpa garis bawah. Interface kecil & spesifik.
- Komentar boleh Bahasa Indonesia (mengikuti gaya scaffold ini).

## Verifikasi sebelum selesai

- `make fmt && make vet && make build` bersih.
- `make test` lolos untuk kode yang kamu sentuh.
- Tidak ada import lintas-module yang melanggar aturan.
- Kontrak & DECISIONS sudah diperbarui bila relevan.

## Batasan

- Jangan menambah dependency besar tanpa alasan yang dicatat di DECISIONS.
- Jangan mengganti keputusan di DECISIONS tanpa menulis ADR baru yang
  men-supersede-nya.
- Jangan mengerjakan `achievement` atau fitur v2 kecuali diminta.
- Ragu soal desain? Berhenti dan tanya, jangan tebak lalu memaksakan.

## Agent Rules

1. Jika memungkinkan selalu gunakan skill `caveman` sebagai default response jawaban dan tanpa mengurangi kualitas jawaban.
2. Jangan kerjakan task sendiri, Selalu panggil subagents jika memungkinkan untuk membagi dan menyelesaikan task yang diberikan.
3. Untuk main agent dan subagents, Selalu gunakan `high` effort dan model untuk mode planning dan `medium` atau `low` effort dan model untuk mode implementasi. sehingga bisa menghasilkan jawaban yang lebih optimal dengan implementasi yang efisien dan efektif untuk menghemat token.

---

## Guard

- Jangan baca envrionment file `.env` dan sejenisnya selain file `.env.example` yang menjadi template tanpa meminta izin dari user
- Jangan tambah co-author dicommit tanpa diminta
