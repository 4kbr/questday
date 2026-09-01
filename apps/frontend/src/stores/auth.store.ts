import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '@/apis/types'

/**
 * Store client-state untuk sesi auth (ADR-018 & ADR-020).
 *
 * HANYA client state. DILARANG menyimpan quest/score/leaderboard di sini — itu
 * server state, milik TanStack Query. `user` di sini adalah *cache sesi* supaya
 * layout bisa langsung menampilkan nama tanpa menunggu `GET /me`; sumber
 * kebenaran tetap query `useMe`.
 */
type AuthState = {
  token: string | null
  user: User | null
  setSession: (token: string, user: User) => void
  setUser: (user: User) => void
  logout: () => void
  isAuthenticated: () => boolean
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      token: null,
      user: null,
      setSession: (token, user) => set({ token, user }),
      setUser: (user) => set({ user }),
      logout: () => set({ token: null, user: null }),
      isAuthenticated: () => !!get().token,
    }),
    {
      name: 'questday-auth',
      // Persist HANYA client state, bukan fungsi.
      partialize: (state) => ({ token: state.token, user: state.user }),
    },
  ),
)

// Wiring `setTokenGetter` / `setUnauthorizedHandler` dipindah ke `@/lib/session`
// (F1.9) supaya handler 401 bisa memakai `endSession()` tanpa circular import.
// `main.tsx` meng-import '@/lib/session' saat bootstrap.
