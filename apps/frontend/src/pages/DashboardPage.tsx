import { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
// Pakai path relatif: oxlint `no-restricted-imports` memblokir alias fitur.
import { QuestFormDialog } from '../features/quest/components/QuestFormDialog'
import { TodayQuestList } from '../features/quest/components/TodayQuestList'

export default function DashboardPage() {
  // Halaman dirender di dalam AppShell via ProtectedRoute — jangan gambar
  // sidebar/topbar di sini.
  const [createOpen, setCreateOpen] = useState(false)

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <h1 className="font-heading text-xl font-semibold">Dashboard</h1>

      {/* F3.5: slot kartu skor & streak di sini. */}

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
