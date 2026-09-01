import { PartyPopper, Sparkles } from 'lucide-react'
import { toast } from 'sonner'
import { EmptyState } from '@/components/EmptyState'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { Skeleton } from '@/components/ui/skeleton'
import { QuestItem } from './QuestItem'
import {
  isBenignQuestToggleError,
  useCompleteQuest,
  useTodayQuests,
  useUncompleteQuest,
} from '../queries/quest.queries'

type TodayQuestListProps = {
  // Dipasang DashboardPage nanti (F2.8) untuk membuka QuestFormDialog.
  onCreate?: () => void
}

// Tanggal ditampilkan dari value backend (ADR-006). Konstruktor Date di sini
// hanya mem-parse string `YYYY-MM-DD` milik backend untuk di-format; ia TIDAK
// dipakai untuk memutuskan "hari apa" — itu sudah dihitung backend dari
// timezone user.
function formatBackendDate(date: string): string {
  return new Intl.DateTimeFormat('id-ID', { dateStyle: 'full' }).format(
    new Date(`${date}T00:00:00`),
  )
}

function LoadingRows() {
  return (
    <div className="space-y-2">
      {[0, 1, 2].map((i) => (
        <div
          key={i}
          className="flex items-center gap-3 rounded-lg border border-border px-3 py-2.5"
        >
          <Skeleton className="size-4 rounded-[4px]" />
          <Skeleton className="h-4 flex-1" />
          <Skeleton className="h-5 w-16 rounded-full" />
          <Skeleton className="h-5 w-14 rounded-full" />
          <Skeleton className="h-4 w-12" />
        </div>
      ))}
    </div>
  )
}

export function TodayQuestList({ onCreate }: TodayQuestListProps) {
  const { data, isPending, isError, refetch } = useTodayQuests()
  const complete = useCompleteQuest()
  const uncomplete = useUncompleteQuest()

  if (isPending) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-5 w-56" />
        <LoadingRows />
      </div>
    )
  }

  if (isError) {
    return (
      <div className="space-y-3 py-6 text-center">
        <p className="text-sm text-muted-foreground">
          Gagal memuat quest hari ini.
        </p>
        <Button variant="outline" size="sm" onClick={() => refetch()}>
          Coba lagi
        </Button>
      </div>
    )
  }

  const items = data.items
  const total = items.length
  const done = items.filter((it) => it.completed).length
  const pct = total === 0 ? 0 : Math.round((done / total) * 100)

  function handleToggle(id: string, currentlyCompleted: boolean) {
    const mutation = currentlyCompleted ? uncomplete : complete
    mutation.mutate(id, {
      onError: (err) => {
        if (!isBenignQuestToggleError(err)) {
          toast.error('Gagal memperbarui quest. Coba lagi.')
        }
      },
    })
  }

  if (total === 0) {
    return (
      <EmptyState
        icon={Sparkles}
        title="Belum ada quest hari ini"
        description="Mulai dengan membuat quest pertamamu — kebiasaan kecil yang ingin kamu jaga tiap hari."
        action={
          onCreate ? (
            <Button size="sm" onClick={onCreate}>
              Buat quest pertama
            </Button>
          ) : null
        }
      />
    )
  }

  return (
    <div className="space-y-4">
      <p className="text-sm text-muted-foreground capitalize">
        {formatBackendDate(data.date)}
      </p>

      <div className="space-y-1.5">
        <p className="text-sm font-medium">
          {done} dari {total} selesai
        </p>
        <Progress value={pct} />
      </div>

      {done === total ? (
        <div className="flex flex-col items-center gap-2 py-6 text-center">
          <PartyPopper className="size-8 text-primary" />
          <p className="text-sm font-medium">
            Semua quest hari ini selesai. Kerja bagus!
          </p>
        </div>
      ) : null}

      <div className="space-y-2">
        {items.map((it) => (
          <QuestItem
            key={it.quest.id}
            quest={it.quest}
            completed={it.completed}
            onToggle={() => handleToggle(it.quest.id, it.completed)}
          />
        ))}
      </div>
    </div>
  )
}
