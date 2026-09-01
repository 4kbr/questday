import { Flame } from 'lucide-react'
import type { Streak } from '@/apis/types'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

/**
 * Kartu streak. **Presentational murni** — angka dari props.
 *
 * `last_active` ditampilkan apa adanya dari backend (sudah tanggal lokal user,
 * ADR-006) — jangan diproses ulang lewat timezone browser.
 */
type StreakCardProps = {
  streak: Streak
}

export function StreakCard({ streak }: StreakCardProps) {
  const cold = streak.current === 0

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-sm font-medium text-muted-foreground">
          Streak
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-2">
        {cold ? (
          <p className="flex items-center gap-2 text-sm font-medium">
            <Flame className="size-5 shrink-0 text-muted-foreground" />
            Selesaikan satu quest untuk memulai streak
          </p>
        ) : (
          <div className="flex items-center gap-2">
            <Flame className="size-6 shrink-0 text-orange-500" />
            <span className="text-3xl font-semibold tabular-nums">
              {streak.current}
            </span>
          </div>
        )}

        <p className="text-xs text-muted-foreground">
          Terpanjang: {streak.longest}
        </p>

        {streak.last_active != null && (
          <p className="text-xs text-muted-foreground">
            Terakhir aktif: {streak.last_active}
          </p>
        )}
      </CardContent>
    </Card>
  )
}
