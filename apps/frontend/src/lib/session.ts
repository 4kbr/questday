import { setTokenGetter, setUnauthorizedHandler } from '@/apis/client'
import { queryClient } from '@/lib/query-client'
import { PATHS } from '@/routes/paths'
import { useAuthStore } from '@/stores/auth.store'

/**
 * Satu-satunya jalur teardown sesi (F1.9).
 *
 * Dipakai dua tempat:
 *  - interceptor 401 di `apis/client.ts` (di luar React Router)
 *  - tombol Logout di `Topbar`
 *
 * WAJIB `queryClient.clear()` — tanpa itu user berikutnya di browser yang sama
 * bisa melihat data user sebelumnya dari cache (kebocoran data antar-akun).
 */
export function endSession() {
  useAuthStore.getState().logout() // token & user -> null, persist bersih
  queryClient.clear() // WAJIB — cegah kebocoran data antar-akun
  // hard redirect: dipakai juga dari interceptor 401 di luar React Router
  if (window.location.pathname !== PATHS.login) {
    window.location.assign(PATHS.login)
  }
}

// --- Wiring ke lapisan HTTP (side-effect, jalan sekali saat modul di-import) ---
// Dipindah dari `stores/auth.store.ts` ke sini untuk memutus circular import:
// `session.ts` butuh `auth.store`, tapi `auth.store` tak lagi butuh `session.ts`.
// `main.tsx` meng-import '@/lib/session' saat bootstrap.
setTokenGetter(() => useAuthStore.getState().token)
setUnauthorizedHandler(() => endSession())
