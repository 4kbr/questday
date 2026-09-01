import { useQuery } from '@tanstack/react-query'
import { scoringApi } from '@/apis/scoring.api'
import { scoringKeys } from './keys'

/**
 * Hanya file ini yang memanggil `scoringApi`. Komponen memanggil hook `use*`,
 * tak pernah `scoringApi` langsung.
 */

export function useScore() {
  return useQuery({
    queryKey: scoringKeys.score(),
    queryFn: () => scoringApi.score().then((r) => r.data),
  })
}

export function useStreak() {
  return useQuery({
    queryKey: scoringKeys.streak(),
    queryFn: () => scoringApi.streak().then((r) => r.data),
  })
}

export function useLeaderboard(limit = 20) {
  return useQuery({
    queryKey: scoringKeys.leaderboard(limit),
    queryFn: () => scoringApi.leaderboard(limit).then((r) => r.data),
  })
}
