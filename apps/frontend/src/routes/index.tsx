import { createBrowserRouter } from 'react-router-dom'
import { PATHS } from '@/routes/paths'
import LoginPage from '@/pages/LoginPage'
import RegisterPage from '@/pages/RegisterPage'
import DashboardPage from '@/pages/DashboardPage'
import QuestsPage from '@/pages/QuestsPage'
import LeaderboardPage from '@/pages/LeaderboardPage'
import SettingsPage from '@/pages/SettingsPage'

export const router = createBrowserRouter([
  { path: PATHS.login, element: <LoginPage /> },
  { path: PATHS.register, element: <RegisterPage /> },
  { path: PATHS.dashboard, element: <DashboardPage /> },
  { path: PATHS.quests, element: <QuestsPage /> },
  { path: PATHS.leaderboard, element: <LeaderboardPage /> },
  { path: PATHS.settings, element: <SettingsPage /> },
])
