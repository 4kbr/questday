import { Suspense } from 'react'
import { Navigate, Outlet } from 'react-router-dom'
import { PageSkeleton } from '@/components/PageSkeleton'
import { PATHS } from '@/routes/paths'
import { useAuthStore } from '@/stores/auth.store'

/**
 * Kebalikan `ProtectedRoute` (F1.7): user yang sudah login tak perlu melihat
 * halaman login/register lagi — lempar ke dashboard.
 */
export function GuestRoute() {
  const token = useAuthStore((s) => s.token)

  if (token) {
    return <Navigate to={PATHS.dashboard} replace />
  }

  return (
    <Suspense fallback={<PageSkeleton />}>
      <Outlet />
    </Suspense>
  )
}
