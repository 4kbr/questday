import { http, HttpResponse } from 'msw'
import type {
  ApiErrorBody,
  AuthResponse,
  HealthResponse,
  User,
} from '@/apis/types'

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
 * Error di-amplop `{error:{code,message}}`; sukses di-amplop `{data}` lewat
 * `dataResponse` (ADR-025).
 */
export function errorResponse(status: number, code: string, message: string) {
  return HttpResponse.json(errorBody(code, message), { status })
}

/**
 * Amplop sukses konsisten dengan ADR-025: `{"data": <payload>}`.
 * `client.ts` meng-unwrap sekali di interceptor sukses, jadi handler sukses
 * WAJIB membungkus payload di sini.
 */
export function dataResponse<T>(payload: T, status = 200) {
  return HttpResponse.json({ data: payload }, { status })
}

// --- Fake in-memory store supaya alur auth koheren tanpa backend --------------

/** User ter-seed; password-nya `'password123'`. */
const seededUser: User = {
  id: '00000000-0000-4000-8000-000000000001',
  email: 'demo@questday.test',
  display_name: 'Demo',
  timezone: 'Asia/Jakarta',
}
const seededPassword = 'password123'

/** Email yang sudah "terdaftar" lewat register — supaya register ulang → 409. */
const registeredEmails = new Set<string>([seededUser.email])

const MOCK_TOKEN = 'mock-jwt-token'

type LoginBody = { email?: string; password?: string }
type RegisterBody = {
  email?: string
  password?: string
  display_name?: string
  timezone?: string
}

export const handlers = [
  http.get(`${BASE}/healthz`, () =>
    HttpResponse.json({ status: 'ok' } satisfies HealthResponse),
  ),

  http.post(`${BASE}/auth/login`, async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as LoginBody
    if (body.email === seededUser.email && body.password === seededPassword) {
      return dataResponse<AuthResponse>(
        { token: MOCK_TOKEN, user: seededUser },
        200,
      )
    }
    return errorResponse(401, 'invalid_credential', 'Email atau password salah')
  }),

  http.post(`${BASE}/auth/register`, async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as RegisterBody
    const email = body.email ?? ''
    if (registeredEmails.has(email)) {
      return errorResponse(409, 'email_taken', 'Email sudah terdaftar')
    }
    const user: User = {
      id: crypto.randomUUID(),
      email,
      display_name: body.display_name ?? '',
      timezone: body.timezone ?? 'Asia/Jakarta',
    }
    registeredEmails.add(email)
    return dataResponse<AuthResponse>({ token: MOCK_TOKEN, user }, 200)
  }),

  http.get(`${BASE}/me`, ({ request }) => {
    const auth = request.headers.get('Authorization') ?? ''
    if (auth.startsWith('Bearer ')) {
      return dataResponse<User>(seededUser, 200)
    }
    return errorResponse(401, 'invalid_credential', 'Token tidak valid')
  }),
]
