import { api } from '@/apis/client'
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  UpdateProfileRequest,
  User,
} from '@/apis/types'

/**
 * Lapisan HTTP murni untuk auth. Tak tahu React, tak menyentuh store, tak
 * melempar toast. `client.ts` sudah meng-unwrap amplop {data} di interceptor —
 * jadi `res.data` di sini = payload sesuai schema (ADR-025). Jangan unwrap lagi.
 */
export const authApi = {
  register: (body: RegisterRequest) =>
    api.post<AuthResponse>('/auth/register', body),
  login: (body: LoginRequest) => api.post<AuthResponse>('/auth/login', body),
  me: () => api.get<User>('/me'),
  updateMe: (body: UpdateProfileRequest) =>
    api.patch<AuthResponse>('/me', body),
}
