import { createBrowserRouter } from 'react-router-dom'
import { PATHS } from '@/routes/paths'
import { GuestRoute } from '@/routes/GuestRoute'
import { ProtectedRoute } from '@/routes/ProtectedRoute'
import LoginPage from '@/pages/LoginPage'
import RegisterPage from '@/pages/RegisterPage'
import DashboardPage from '@/pages/DashboardPage'
import QuestsPage from '@/pages/QuestsPage'
import LeaderboardPage from '@/pages/LeaderboardPage'
import SettingsPage from '@/pages/SettingsPage'

export const router = createBrowserRouter([
  {
    element: <GuestRoute />,
    children: [
      { path: PATHS.login, element: <LoginPage /> },
      { path: PATHS.register, element: <RegisterPage /> },
    ],
  },
  {
    element: <ProtectedRoute />,
    children: [
      { index: true, element: <DashboardPage /> },
      { path: PATHS.quests, element: <QuestsPage /> },
      { path: PATHS.leaderboard, element: <LeaderboardPage /> },
      { path: PATHS.settings, element: <SettingsPage /> },
    ],
  },
])
