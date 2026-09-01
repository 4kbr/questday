import type { Quest } from '@/apis/types'
import { Badge } from '@/components/ui/badge'
import { Checkbox } from '@/components/ui/checkbox'
import { cn } from '@/lib/utils'
import { difficultyClass, difficultyLabel } from '../lib/difficulty'

/**
 * Satu baris quest. **Presentational murni** — tak memanggil hook mutation
 * sendiri. Status & aksi datang lewat props supaya gampang dipakai ulang &
 * dites (F2.5).
 */
type QuestItemProps = {
  quest: Quest
  completed: boolean
  onToggle: () => void
}

export function QuestItem({ quest, completed, onToggle }: QuestItemProps) {
  return (
    <label className="flex cursor-pointer items-center gap-3 rounded-lg border border-border px-3 py-2.5">
      <Checkbox
        checked={completed}
        onCheckedChange={onToggle}
        aria-label={quest.title}
      />

      <span
        className={cn(
          'min-w-0 flex-1 truncate text-sm font-medium',
          completed && 'text-muted-foreground line-through',
        )}
      >
        {quest.title}
      </span>

      <Badge variant="outline" className="shrink-0">
        {quest.category}
      </Badge>

      <Badge
        variant="outline"
        className={cn(
          'shrink-0 border-transparent',
          difficultyClass[quest.difficulty],
        )}
      >
        {difficultyLabel[quest.difficulty]}
      </Badge>

      <span className="shrink-0 text-xs text-muted-foreground tabular-nums">
        {quest.points} poin
      </span>
    </label>
  )
}
