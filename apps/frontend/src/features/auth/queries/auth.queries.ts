import { useMutation, useQuery } from '@tanstack/react-query'
import { authApi } from '@/apis/auth.api'
import { useAuthStore } from '@/stores/auth.store'

export const authKeys = { me: () => ['auth', 'me'] as const }

/**
 * Hanya hook di folder ini yang memanggil `authApi`. Komponen memanggil hook,
 * tak pernah `authApi` langsung.
 */
export function useMe() {
  const token = useAuthStore((s) => s.token)
  return useQuery({
    queryKey: authKeys.me(),
    queryFn: () => authApi.me().then((r) => r.data),
    enabled: !!token,
  })
}

export function useLogin() {
  return useMutation({
    mutationFn: authApi.login,
    onSuccess: (r) => {
      useAuthStore.getState().setSession(r.data.token, r.data.user)
    },
  })
}

export function useRegister() {
  return useMutation({
    mutationFn: authApi.register,
    onSuccess: (r) => {
      useAuthStore.getState().setSession(r.data.token, r.data.user)
    },
  })
}
