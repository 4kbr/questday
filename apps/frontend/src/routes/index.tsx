// File router: definisi `lazy()` per-route hidup berdampingan dengan export
// `router`. Aturan fast-refresh tak relevan di sini (bukan file komponen).
/* oxlint-disable react/only-export-components */
import { lazy } from 'react'
import { createBrowserRouter } from 'react-router-dom'
import { PATHS } from '@/routes/paths'
import { GuestRoute } from '@/routes/GuestRoute'
import { ProtectedRoute } from '@/routes/ProtectedRoute'
import { RouteError } from '@/routes/RouteError'
import NotFoundPage from '@/pages/NotFoundPage'

// Route code-splitting (F4.10): tiap halaman jadi chunk terpisah supaya bundle
// awal jauh lebih kecil. Fallback `Suspense` (`<PageSkeleton/>`) hidup di
// `GuestRoute` / `ProtectedRoute` — untuk rute terproteksi ia berada DI DALAM
// `AppShell`, jadi chrome (sidebar/topbar) tetap terlihat saat chunk halaman
// diunduh. `NotFoundPage` sengaja eager: catch-all top-level tanpa induk
// `Suspense`, dan ukurannya sepele.
const LoginPage = lazy(() => import('@/pages/LoginPage'))
const RegisterPage = lazy(() => import('@/pages/RegisterPage'))
const DashboardPage = lazy(() => import('@/pages/DashboardPage'))
const QuestsPage = lazy(() => import('@/pages/QuestsPage'))
const LeaderboardPage = lazy(() => import('@/pages/LeaderboardPage'))
const SettingsPage = lazy(() => import('@/pages/SettingsPage'))

export const router = createBrowserRouter([
  {
    element: <GuestRoute />,
    errorElement: <RouteError />,
    children: [
      { path: PATHS.login, element: <LoginPage /> },
      { path: PATHS.register, element: <RegisterPage /> },
    ],
  },
  {
    element: <ProtectedRoute />,
    errorElement: <RouteError />,
    children: [
      { index: true, element: <DashboardPage /> },
      { path: PATHS.quests, element: <QuestsPage /> },
      { path: PATHS.leaderboard, element: <LeaderboardPage /> },
      { path: PATHS.settings, element: <SettingsPage /> },
    ],
  },
  { path: '*', element: <NotFoundPage /> },
])
