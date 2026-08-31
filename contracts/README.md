# Contracts

Kontrak API QuestDay. **Sumber kebenaran bersama** antara `apps/backend`
(yang meng-implement) dan `apps/frontend` (yang meng-consume).

Ditaruh di root, bukan di dalam salah satu app, supaya tidak ada app yang
"memiliki" kontrak — keduanya sama-sama tunduk ke sini.

## Isi

- `openapi.yaml` — spesifikasi REST API (OpenAPI 3.1).

## Alur pemakaian (disarankan)

1. Ubah kontrak di sini **lebih dulu** saat menambah/mengubah endpoint.
2. Backend menyesuaikan handler + DTO agar cocok dengan kontrak.
3. Frontend generate client/type dari file ini (nanti).

Lihat `docs/GUIDES.md` → "Menambah endpoint API".
