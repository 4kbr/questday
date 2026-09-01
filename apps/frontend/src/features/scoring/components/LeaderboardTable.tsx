import { Medal } from 'lucide-react'
import type { LeaderboardEntry } from '@/apis/types'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { cn } from '@/lib/utils'

/**
 * Tabel leaderboard. **Presentational murni** — data & loading/error di halaman.
 *
 * ATURAN: `rank` datang dari `entry.rank` (backend), JANGAN dihitung dari indeks
 * array — kalau backend menambah paging nanti, indeks lokal akan salah.
 */
type LeaderboardTableProps = {
  entries: LeaderboardEntry[]
  currentUserId?: string
}

// Warna medali untuk rank 1-3; rank >= 4 tak dapat.
const MEDAL_CLASS: Record<number, string> = {
  1: 'text-yellow-500',
  2: 'text-zinc-400',
  3: 'text-amber-700',
}

function RankCell({ rank }: { rank: number }) {
  const medal = MEDAL_CLASS[rank]
  if (medal) {
    return (
      <span className="inline-flex items-center gap-1 font-medium tabular-nums">
        <Medal className={cn('size-4', medal)} aria-hidden />
        {rank}
      </span>
    )
  }
  return <span className="tabular-nums">{rank}</span>
}

export function LeaderboardTable({
  entries,
  currentUserId,
}: LeaderboardTableProps) {
  return (
    <div className="w-full overflow-x-auto">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-16">Rank</TableHead>
            <TableHead>Nama</TableHead>
            <TableHead className="text-right">Poin</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {entries.map((entry) => {
            const isMe = entry.user_id === currentUserId
            // Backend menjaga display_name (T3.6), tapi FE jangan anggap mustahil.
            const name = entry.display_name?.trim()
              ? entry.display_name
              : 'Pengguna dihapus'
            return (
              <TableRow
                key={entry.user_id}
                className={cn(isMe && 'bg-muted/50 font-medium')}
              >
                <TableCell>
                  <RankCell rank={entry.rank} />
                </TableCell>
                <TableCell>
                  {name}
                  {isMe && (
                    <span className="ml-1 text-muted-foreground">(kamu)</span>
                  )}
                </TableCell>
                <TableCell className="text-right tabular-nums">
                  {entry.points}
                </TableCell>
              </TableRow>
            )
          })}
        </TableBody>
      </Table>
    </div>
  )
}
