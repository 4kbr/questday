import { useState } from 'react'
import { ClipboardList } from 'lucide-react'
import { EmptyState } from '@/components/EmptyState'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { toastApiError, toastSuccess } from '@/lib/toast'
// Path relatif: oxlint `no-restricted-imports` memblokir alias fitur.
import { QuestFormDialog } from '../features/quest/components/QuestFormDialog'
import { QuestTable } from '../features/quest/components/QuestTable'
import {
  useArchiveQuest,
  useQuests,
} from '../features/quest/queries/quest.queries'
import type { Quest } from '@/apis/types'

type DialogState = {
  open: boolean
  mode: 'create' | 'edit'
  quest?: Quest
}

function LoadingRows() {
  return (
    <div className="space-y-2 rounded-lg border border-border p-4">
      {[0, 1, 2, 3, 4].map((i) => (
        <div key={i} className="flex items-center gap-4 py-1.5">
          <Skeleton className="h-4 flex-1" />
          <Skeleton className="h-4 w-24" />
          <Skeleton className="h-5 w-16 rounded-full" />
          <Skeleton className="h-4 w-10" />
          <Skeleton className="h-7 w-7 rounded-md" />
        </div>
      ))}
    </div>
  )
}

export default function QuestsPage() {
  const { data: quests, isPending, isError, refetch } = useQuests()
  const archive = useArchiveQuest()

  const [dialog, setDialog] = useState<DialogState>({
    open: false,
    mode: 'create',
  })
  const [archiving, setArchiving] = useState<Quest | null>(null)

  function openCreate() {
    setDialog({ open: true, mode: 'create', quest: undefined })
  }

  function openEdit(quest: Quest) {
    setDialog({ open: true, mode: 'edit', quest })
  }

  function confirmArchive() {
    if (!archiving) return
    archive.mutate(archiving.id, {
      onSuccess: () => {
        toastSuccess('Quest diarsipkan')
        setArchiving(null)
      },
      onError: (err) => toastApiError(err),
    })
  }

  return (
    <div className="mx-auto max-w-4xl space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="font-heading text-xl font-semibold">Quests</h1>
        <Button onClick={openCreate}>Quest Baru</Button>
      </div>

      {isPending ? (
        <LoadingRows />
      ) : isError ? (
        <div className="space-y-3 rounded-lg border border-border py-10 text-center">
          <p className="text-sm text-muted-foreground">
            Gagal memuat daftar quest.
          </p>
          <Button variant="outline" size="sm" onClick={() => refetch()}>
            Coba lagi
          </Button>
        </div>
      ) : quests.length === 0 ? (
        <EmptyState
          className="rounded-lg border border-border"
          icon={ClipboardList}
          title="Belum ada quest"
          description="Quest adalah kebiasaan harian yang ingin kamu jaga. Buat yang pertama."
          action={
            <Button size="sm" onClick={openCreate}>
              Buat quest pertama
            </Button>
          }
        />
      ) : (
        <QuestTable
          quests={quests}
          onEdit={openEdit}
          onArchive={(quest) => setArchiving(quest)}
        />
      )}

      <QuestFormDialog
        mode={dialog.mode}
        quest={dialog.quest}
        open={dialog.open}
        onOpenChange={(open) => setDialog((prev) => ({ ...prev, open }))}
      />

      <AlertDialog
        open={archiving !== null}
        onOpenChange={(open) => {
          if (!open) setArchiving(null)
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Arsipkan quest?</AlertDialogTitle>
            <AlertDialogDescription>
              Quest{archiving ? ` "${archiving.title}"` : ''} akan{' '}
              <strong>diarsipkan</strong>, bukan dihapus permanen. Ia berhenti
              muncul di daftar hari ini, tapi riwayat penyelesaian dan streak
              kamu tetap utuh. Quest bisa dipulihkan lagi nanti.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={archive.isPending}>
              Batal
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={(e) => {
                e.preventDefault()
                confirmArchive()
              }}
              disabled={archive.isPending}
            >
              Arsipkan
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
