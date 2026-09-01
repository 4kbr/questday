import { Suspense } from 'react'
import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { AppShell } from '@/components/layout/AppShell'
import { PageSkeleton } from '@/components/PageSkeleton'
import { PATHS } from '@/routes/paths'
import { useAuthStore } from '@/stores/auth.store'

/**
 * Gerbang rute terproteksi (F1.7).
 *
 * Tak ada token → lempar ke `/login`, simpan lokasi asal di `state.from` supaya
 * `LoginForm` bisa mengembalikan user ke halaman yang tadi dituju.
 * Ada token → render shell + `<Outlet/>`.
 */
export function ProtectedRoute() {
  const token = useAuthStore((s) => s.token)
  const location = useLocation()

  if (!token) {
    return <Navigate to={PATHS.login} replace state={{ from: location }} />
  }

  // Satu `Suspense` di dalam `AppShell`: saat chunk halaman (lazy) diunduh,
  // sidebar/topbar tetap tampil dan hanya area konten yang menampilkan skeleton.
  return (
    <AppShell>
      <Suspense fallback={<PageSkeleton />}>
        <Outlet />
      </Suspense>
    </AppShell>
  )
}
