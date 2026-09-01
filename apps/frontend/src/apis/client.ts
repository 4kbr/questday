import axios from 'axios'
import type { AxiosError, InternalAxiosRequestConfig } from 'axios'

/**
 * Error terpadu untuk semua kegagalan HTTP.
 *
 * Catatan: field ditulis eksplisit (bukan TS "parameter properties") karena
 * tsconfig mengaktifkan `erasableSyntaxOnly` — konstruksi TS-only tidak boleh
 * menghasilkan kode runtime.
 */
export class ApiError extends Error {
  readonly status: number
  readonly code: string

  constructor(status: number, code: string, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
  }
}

// Satu-satunya instance HTTP di seluruh app. Jangan bikin `axios.create` lain.
export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
})

/**
 * Akses token bisa disuntik supaya `stores/auth.store` (F1.2) bisa mem-wire
 * tanpa circular import. Default: baca langsung dari localStorage.
 *
 * PENTING untuk F1.2: zustand `persist` WAJIB memakai `name: 'questday.auth'`
 * agar key di bawah ini cocok, dan token disimpan di `state.token`.
 */
let tokenGetter: () => string | null = () => {
  try {
    return (
      JSON.parse(localStorage.getItem('questday.auth') ?? 'null')?.state
        ?.token ?? null
    )
  } catch {
    return null
  }
}

export function setTokenGetter(fn: () => string | null) {
  tokenGetter = fn
}

// Dipanggil saat 401. Default no-op; F1.2 menyuntik logout + redirect ke /login.
let onUnauthorized: () => void = () => {}

export function setUnauthorizedHandler(fn: () => void) {
  onUnauthorized = fn
}

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = tokenGetter()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  // Sukses: response mentah, TANPA unwrap. Kontrak = sumber kebenaran, dan
  // sukses adalah objek RAW (tidak ada envelope `{data:...}`).
  (response) => response,
  (error: AxiosError<{ error?: { code?: string; message?: string } }>) => {
    if (!error.response) {
      // Network error / tak ada response.
      return Promise.reject(new ApiError(0, 'internal_error', error.message))
    }

    const status = error.response.status
    const body = error.response.data
    const apiErr =
      body?.error?.code && body?.error?.message
        ? new ApiError(status, body.error.code, body.error.message)
        : new ApiError(status, 'internal_error', error.message)

    if (status === 401) {
      onUnauthorized()
    }

    return Promise.reject(apiErr)
  },
)
