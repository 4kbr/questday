import { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
// Pakai path relatif: oxlint `no-restricted-imports` memblokir alias fitur.
import { QuestFormDialog } from '../features/quest/components/QuestFormDialog'
import { TodayQuestList } from '../features/quest/components/TodayQuestList'
import { ScoreCard } from '../features/scoring/components/ScoreCard'
import { StreakCard } from '../features/scoring/components/StreakCard'
import {
  useScore,
  useStreak,
} from '../features/scoring/queries/scoring.queries'

// Skeleton berbentuk kartu — dipakai saat kedua query masih loading.
function CardSkeleton() {
  return (
    <Card>
      <CardHeader>
        <Skeleton className="h-4 w-16" />
      </CardHeader>
      <CardContent className="space-y-3">
        <Skeleton className="h-8 w-24" />
        <Skeleton className="h-1 w-full" />
        <Skeleton className="h-3 w-40" />
      </CardContent>
    </Card>
  )
}

export default function DashboardPage() {
  // Halaman dirender di dalam AppShell via ProtectedRoute — jangan gambar
  // sidebar/topbar di sini.
  const [createOpen, setCreateOpen] = useState(false)

  // Halaman yang mengorkestrasi data; kartu tetap presentational.
  const score = useScore()
  const streak = useStreak()

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <h1 className="font-heading text-xl font-semibold">Dashboard</h1>

      <div className="grid gap-4 sm:grid-cols-2">
        {score.isLoading || streak.isLoading ? (
          <>
            <CardSkeleton />
            <CardSkeleton />
          </>
        ) : (
          <>
            {score.isError || !score.data ? (
              <Card>
                <CardContent className="py-6 text-sm text-muted-foreground">
                  Gagal memuat skor
                </CardContent>
              </Card>
            ) : (
              <ScoreCard score={score.data} />
            )}

            {streak.isError || !streak.data ? (
              <Card>
                <CardContent className="py-6 text-sm text-muted-foreground">
                  Gagal memuat streak
                </CardContent>
              </Card>
            ) : (
              <StreakCard streak={streak.data} />
            )}
          </>
        )}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Quest hari ini</CardTitle>
        </CardHeader>
        <CardContent>
          <TodayQuestList onCreate={() => setCreateOpen(true)} />
        </CardContent>
      </Card>

      <QuestFormDialog
        mode="create"
        open={createOpen}
        onOpenChange={setCreateOpen}
      />
    </div>
  )
}
