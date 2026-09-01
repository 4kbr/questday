import { QueryClient } from '@tanstack/react-query'
import { ApiError } from '@/apis/client'

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      refetchOnWindowFocus: true,
      // Jangan retry 4xx — 401/404/409 tak membaik dengan diulang.
      retry: (count, err) =>
        !(err instanceof ApiError && err.status >= 400 && err.status < 500) &&
        count < 2,
    },
  },
})
