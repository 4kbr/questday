import { api } from '@/apis/client'
import type { LeaderboardEntry, Score, Streak } from '@/apis/types'

/**
 * Lapisan HTTP murni untuk scoring. Tak tahu React, nol logika. `client.ts`
 * sudah meng-unwrap amplop {data} di interceptor sukses — jadi `res.data` di
 * sini = payload sesuai schema (ADR-025). Jangan unwrap lagi.
 *
 * `limit` default 20 = default backend, supaya perilaku identik antara dikirim
 * dan tidak.
 */
export const scoringApi = {
  score: () => api.get<Score>('/me/score'),
  streak: () => api.get<Streak>('/me/streak'),
  leaderboard: (limit = 20) =>
    api.get<LeaderboardEntry[]>('/leaderboard', { params: { limit } }),
}
