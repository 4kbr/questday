/**
 * Daftar IANA timezone umum (fallback browser lama) + util default browser.
 * Dipakai `RegisterForm` (F1.6) dan `TimezoneSelect` (F4.2). Ditaruh di `lib/`
 * (bukan `features/auth/lib/`) karena lintas-fitur — oxlint memblokir alias
 * fitur. Timezone menentukan batas hari untuk streak (ADR-006) — wajib IANA
 * valid.
 */
export const COMMON_TIMEZONES = [
  'Asia/Jakarta',
  'Asia/Makassar',
  'Asia/Jayapura',
  'Asia/Singapore',
  'Asia/Tokyo',
  'Europe/London',
  'America/New_York',
  'America/Los_Angeles',
  'UTC',
] as const

export function browserTimezone(): string {
  return Intl.DateTimeFormat().resolvedOptions().timeZone || 'Asia/Jakarta'
}

/**
 * Opsi untuk `Select`. Kalau zona browser tak ada di daftar umum, taruh di
 * depan supaya nilai default selalu bisa dipilih.
 */
export function timezoneOptions(): string[] {
  const list = [...COMMON_TIMEZONES] as string[]
  const tz = browserTimezone()
  return list.includes(tz) ? list : [tz, ...list]
}
