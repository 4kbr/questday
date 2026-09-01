import { z } from 'zod'
import type {
  LoginRequest,
  RegisterRequest,
  UpdateProfileRequest,
} from '@/apis/types'

/**
 * Validasi client = kenyamanan, bukan pengaman. Aturan WAJIB cocok dengan
 * backend (password min 8, email, display_name required, timezone opsional).
 */
export const loginSchema = z.object({
  email: z.string().email('Email tidak valid'),
  password: z.string().min(1, 'Password wajib diisi'),
})

export const registerSchema = z.object({
  email: z.string().email('Email tidak valid'),
  password: z.string().min(8, 'Minimal 8 karakter'),
  display_name: z.string().min(1, 'Nama wajib diisi'),
  timezone: z.string().optional(),
})

/**
 * Form profil (F4.2). `timezone` selalu ada di form (default dari user), jadi
 * `min(1)` — bukan `optional()` seperti di kontrak `UpdateProfileRequest`.
 */
export const profileSchema = z.object({
  display_name: z.string().min(1, 'Nama wajib diisi'),
  timezone: z.string().min(1),
})

export type LoginValues = z.infer<typeof loginSchema>
export type RegisterValues = z.infer<typeof registerSchema>
export type ProfileValues = z.infer<typeof profileSchema>

// type-check: gagal `tsc` bila schema menyimpang dari kontrak backend.
export type _LoginSyncCheck = LoginValues extends LoginRequest ? true : never
export type _RegisterSyncCheck = RegisterValues extends RegisterRequest
  ? true
  : never
export type _ProfileSyncCheck = ProfileValues extends UpdateProfileRequest
  ? true
  : never
