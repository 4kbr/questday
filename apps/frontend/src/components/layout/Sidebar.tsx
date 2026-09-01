import { NavLink } from 'react-router-dom'
import { LayoutDashboard, ListChecks, Settings, Trophy } from 'lucide-react'
import { cn } from '@/lib/utils'
import { PATHS } from '@/routes/paths'

const NAV = [
  { to: PATHS.dashboard, label: 'Dashboard', icon: LayoutDashboard, end: true },
  { to: PATHS.quests, label: 'Quests', icon: ListChecks, end: false },
  { to: PATHS.leaderboard, label: 'Leaderboard', icon: Trophy, end: false },
  { to: PATHS.settings, label: 'Settings', icon: Settings, end: false },
]

export function Sidebar() {
  return (
    <aside className="flex w-60 shrink-0 flex-col border-r bg-sidebar text-sidebar-foreground">
      <div className="flex h-14 items-center px-6 text-lg font-bold">
        QuestDay
      </div>
      <nav className="flex flex-col gap-1 p-3">
        {NAV.map(({ to, label, icon: Icon, end }) => (
          <NavLink
            key={to}
            to={to}
            end={end}
            className={({ isActive }) =>
              cn(
                'flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors',
                isActive
                  ? 'bg-sidebar-accent text-sidebar-accent-foreground'
                  : 'text-muted-foreground hover:bg-sidebar-accent/50 hover:text-sidebar-accent-foreground',
              )
            }
          >
            <Icon className="size-4" />
            {label}
          </NavLink>
        ))}
      </nav>
    </aside>
  )
}
