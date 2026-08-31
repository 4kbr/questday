# Frontend

Web app QuestDay — mengonsumsi API dari `apps/backend` sesuai kontrak di
`contracts/openapi.yaml`.

**Status:** belum dikerjakan. Rencana kerjanya sudah dipecah per phase di
[`docs/tasks/`](docs/tasks/README.md) — mulai dari sana.

## Stack

React + Vite + TypeScript · TanStack Query (server state) · Zustand (client
state) · shadcn/ui + Tailwind CSS · React Router v7 · react-hook-form + zod ·
MSW (mock dev).

Alasan tiap pilihan dicatat sebagai ADR-017..022 di `docs/DECISIONS.md`.

## Mulai cepat (setelah Phase 0 dikerjakan)

```bash
cd apps/frontend
npm install
cp .env.example .env
npm run gen:api      # generate type dari contracts/openapi.yaml
npm run dev
```

Selama backend belum jalan, set `VITE_USE_MOCK=true` di `.env` — MSW akan
menjawab request sesuai kontrak.

## Aturan yang tak boleh dilanggar

1. **HTTP hanya di `src/apis/`.** Halaman & komponen tak pernah memanggil
   `axios`/`fetch` langsung — selalu lewat hook di `features/*/queries/`.
2. **`src/apis/schema.gen.ts` hasil generate** — jangan diedit tangan. Kontrak
   berubah → `npm run gen:api`.
3. **Server state ≠ client state.** Data API → TanStack Query. Token & UI →
   Zustand.
4. **"Hari ini" ikut timezone user** (ADR-006) — pakai tanggal dari response
   backend, bukan `new Date()`.

Selengkapnya di [`docs/tasks/README.md`](docs/tasks/README.md).
