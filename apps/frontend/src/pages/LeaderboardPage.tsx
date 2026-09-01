import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Skeleton } from '@/components/ui/skeleton'
import { useAuthStore } from '@/stores/auth.store'
// Path relatif: oxlint `no-restricted-imports` memblokir alias fitur.
import { LeaderboardTable } from '../features/scoring/components/LeaderboardTable'
import { useLeaderboard } from '../features/scoring/queries/scoring.queries'

// Skeleton baris tabel saat loading.
function TableSkeleton() {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-16">Rank</TableHead>
          <TableHead>Nama</TableHead>
          <TableHead className="text-right">Poin</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {Array.from({ length: 5 }).map((_, i) => (
          <TableRow key={i}>
            <TableCell>
              <Skeleton className="h-4 w-6" />
            </TableCell>
            <TableCell>
              <Skeleton className="h-4 w-32" />
            </TableCell>
            <TableCell className="text-right">
              <Skeleton className="ml-auto h-4 w-10" />
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}

export default function LeaderboardPage() {
  // Dirender di dalam AppShell via ProtectedRoute — jangan gambar chrome layout.
  const leaderboard = useLeaderboard()
  const currentUserId = useAuthStore((s) => s.user?.id)

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <h1 className="font-heading text-xl font-semibold">Leaderboard</h1>

      {leaderboard.isLoading ? (
        <TableSkeleton />
      ) : leaderboard.isError ? (
        <div className="space-y-3 py-6 text-sm text-muted-foreground">
          <p>Gagal memuat leaderboard.</p>
          <Button variant="outline" onClick={() => leaderboard.refetch()}>
            Coba lagi
          </Button>
        </div>
      ) : !leaderboard.data || leaderboard.data.length === 0 ? (
        <p className="py-6 text-sm text-muted-foreground">
          Belum ada yang mengumpulkan poin
        </p>
      ) : (
        <LeaderboardTable
          entries={leaderboard.data}
          currentUserId={currentUserId}
        />
      )}
    </div>
  )
}
