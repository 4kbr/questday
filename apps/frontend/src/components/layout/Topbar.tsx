import { useLocation, useNavigate } from 'react-router-dom'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { endSession } from '@/lib/session'
import { PATHS } from '@/routes/paths'
import { useAuthStore } from '@/stores/auth.store'

const TITLES: Record<string, string> = {
  [PATHS.dashboard]: 'Dashboard',
  [PATHS.quests]: 'Quests',
  [PATHS.leaderboard]: 'Leaderboard',
  [PATHS.settings]: 'Settings',
}

function initials(name: string): string {
  const parts = name.trim().split(/\s+/).filter(Boolean)
  const first = parts[0] ?? ''
  const last = parts.length > 1 ? (parts[parts.length - 1] ?? '') : ''
  const out = ((first[0] ?? '') + (last[0] ?? first[1] ?? '')).toUpperCase()
  return out || '?'
}

export function Topbar() {
  const location = useLocation()
  const navigate = useNavigate()
  const user = useAuthStore((s) => s.user)
  const title = TITLES[location.pathname] ?? 'QuestDay'
  const name = user?.display_name ?? 'User'

  return (
    <header className="flex h-14 shrink-0 items-center justify-between border-b px-6">
      <h1 className="text-sm font-semibold">{title}</h1>
      <DropdownMenu>
        <DropdownMenuTrigger className="rounded-full outline-none focus-visible:ring-2 focus-visible:ring-ring">
          <Avatar className="size-8">
            <AvatarFallback>{initials(name)}</AvatarFallback>
          </Avatar>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-56">
          <DropdownMenuLabel className="flex flex-col gap-0.5">
            <span className="font-medium">{name}</span>
            <span className="text-xs font-normal text-muted-foreground">
              {user?.email}
            </span>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem onSelect={() => navigate(PATHS.settings)}>
            Settings
          </DropdownMenuItem>
          <DropdownMenuItem onSelect={() => endSession()}>
            Logout
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </header>
  )
}
