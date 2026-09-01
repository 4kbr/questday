import { useMutation } from '@tanstack/react-query'
import { authApi } from '@/apis/auth.api'
import type { UpdateProfileRequest } from '@/apis/types'
import { queryClient } from '@/lib/query-client'
import { useAuthStore } from '@/stores/auth.store'
// Path relatif lintas-fitur: KONSTANTA key saja (oxlint memblokir alias fitur).
import { questKeys } from '../../quest/queries/keys'
import { scoringKeys } from '../../scoring/queries/keys'
import { authKeys } from './auth.queries'

/**
 * `PATCH /me` — satu-satunya pemanggil `authApi.updateMe`. Komponen memanggil
 * hook ini, tak pernah `authApi` langsung.
 *
 * ADR-022/ADR-013: response `AuthResponse` membawa token BARU karena timezone
 * ikut di JWT claims. Token lama masih memuat timezone lama → backend terus
 * menghitung "hari ini" dengan zona lama sampai token kedaluwarsa. Karena itu
 * `setSession(token, user)` wajib, lalu quest hari ini & scoring di-invalidate
 * (batas hari bisa bergeser).
 */
export function useUpdateProfile() {
  return useMutation({
    mutationFn: (body: UpdateProfileRequest) =>
      authApi.updateMe(body).then((r) => r.data),
    onSuccess: ({ token, user }) => {
      useAuthStore.getState().setSession(token, user)
      queryClient.invalidateQueries({ queryKey: authKeys.me() })
      queryClient.invalidateQueries({ queryKey: questKeys.today() })
      queryClient.invalidateQueries({ queryKey: scoringKeys.all() })
    },
  })
}
