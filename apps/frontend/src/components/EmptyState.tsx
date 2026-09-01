import type { LucideIcon } from 'lucide-react'
import type { ReactNode } from 'react'
import { cn } from '@/lib/utils'

/**
 * UI bersama, presentational murni: ikon + judul + deskripsi + aksi opsional.
 * Dipakai untuk empty-state "belum ada quest" maupun perayaan "semua selesai".
 * Tak tahu apa-apa soal quest — makanya tinggal di `components/`, bukan
 * `features/`.
 */
type EmptyStateProps = {
  icon: LucideIcon
  title: string
  description?: string
  action?: ReactNode
  className?: string
}

export function EmptyState({
  icon: Icon,
  title,
  description,
  action,
  className,
}: EmptyStateProps) {
  return (
    <div
      className={cn(
        'flex flex-col items-center gap-3 py-10 text-center',
        className,
      )}
    >
      <Icon className="size-8 text-muted-foreground" />
      <div className="space-y-1">
        <p className="text-sm font-medium">{title}</p>
        {description ? (
          <p className="text-sm text-muted-foreground">{description}</p>
        ) : null}
      </div>
      {action}
    </div>
  )
}
