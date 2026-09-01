import { z } from 'zod'
import type { LoginRequest, RegisterRequest } from '@/apis/types'

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

export type LoginValues = z.infer<typeof loginSchema>
export type RegisterValues = z.infer<typeof registerSchema>

// type-check: gagal `tsc` bila schema menyimpang dari kontrak backend.
export type _LoginSyncCheck = LoginValues extends LoginRequest ? true : never
export type _RegisterSyncCheck = RegisterValues extends RegisterRequest
  ? true
  : never
