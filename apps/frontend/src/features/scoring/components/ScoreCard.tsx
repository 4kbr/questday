import type { Score } from '@/apis/types'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'

/**
 * Kartu skor. **Presentational murni** — semua angka datang langsung dari
 * props; halaman yang mengambil data & menangani loading/error.
 *
 * ATURAN KERAS: jangan menghitung ulang level atau ambang XP di frontend.
 * Aturan level satu sumber di backend (`scoring.LevelForXP`, ADR-007).
 */
type ScoreCardProps = {
  score: Score
}

export function ScoreCard({ score }: ScoreCardProps) {
  // Rasio TAMPILAN untuk bar — bukan rumus level. Hanya memakai angka yang
  // backend beri (`xp` + `points_to_next_level`); tak butuh tabel ambang.
  // Denominator 0 (mis. skor awal) → 0 supaya tak NaN.
  const denom = score.xp + score.points_to_next_level
  const progress = denom === 0 ? 0 : (score.xp / denom) * 100

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-sm font-medium text-muted-foreground">
          Poin
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="flex items-baseline justify-between gap-2">
          <span className="text-3xl font-semibold tabular-nums">
            {score.total_points}
          </span>
          <Badge variant="secondary">Level {score.level}</Badge>
        </div>

        <Progress value={progress} />

        <p className="text-xs text-muted-foreground">
          {score.points_to_next_level} XP lagi menuju Level {score.level + 1}
        </p>
      </CardContent>
    </Card>
  )
}
