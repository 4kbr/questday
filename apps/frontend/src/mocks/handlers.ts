import { http, HttpResponse } from 'msw'
import type { ApiErrorBody, HealthResponse } from '@/apis/types'

/**
 * Base URL untuk semua handler. Sama dengan `baseURL` axios di `apis/client.ts`
 * sehingga request yang keluar cocok dengan path lengkap yang di-intercept.
 * Contoh: `http://localhost:8080/api/v1`.
 */
const BASE = import.meta.env.VITE_API_BASE_URL

/**
 * Amplop error konsisten dengan kontrak: `{"error":{code,message}}`.
 * Dipakai ulang oleh handler fitur di phase berikutnya (F1.10 / F2.11 / F3.8).
 */
export function errorBody(code: string, message: string): ApiErrorBody {
  return { error: { code, message } }
}

/**
 * Bangun response error ber-amplop dengan status HTTP tertentu.
 * Sukses tetap RAW (tanpa `{data}`), hanya error yang di-amplop — sesuai kontrak.
 */
export function errorResponse(status: number, code: string, message: string) {
  return HttpResponse.json(errorBody(code, message), { status })
}

export const handlers = [
  http.get(`${BASE}/healthz`, () =>
    HttpResponse.json({ status: 'ok' } satisfies HealthResponse),
  ),
]
