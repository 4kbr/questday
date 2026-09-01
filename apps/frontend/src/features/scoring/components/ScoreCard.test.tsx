import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import type { Score } from '@/apis/types'
import { ScoreCard } from './ScoreCard'

describe('ScoreCard', () => {
  it('skor 0: total 0, Level 1, teks "100 XP lagi menuju Level 2", bar ada', () => {
    const score: Score = {
      total_points: 0,
      xp: 0,
      level: 1,
      points_to_next_level: 100,
    }
    render(<ScoreCard score={score} />)

    expect(screen.getByText('0')).toBeInTheDocument()
    expect(screen.getByText('Level 1')).toBeInTheDocument()
    expect(screen.getByText('100 XP lagi menuju Level 2')).toBeInTheDocument()
    expect(screen.getByRole('progressbar')).toBeInTheDocument()
  })

  it('skor tengah: angka tampil PERSIS seperti props', () => {
    const score: Score = {
      total_points: 250,
      xp: 250,
      level: 3,
      points_to_next_level: 50,
    }
    render(<ScoreCard score={score} />)

    expect(screen.getByText('250')).toBeInTheDocument()
    expect(screen.getByText('Level 3')).toBeInTheDocument()
    expect(screen.getByText('50 XP lagi menuju Level 4')).toBeInTheDocument()
  })

  it('level & points_to_next_level dipakai VERBATIM dari props, bukan dihitung', () => {
    // Sengaja: level (9) & ambang (7) TIDAK cocok dengan rumus lokal naif
    // (floor dari xp dibagi seratus, plus satu) yang akan bilang Level 3.
    // Kalau ada yang menyisipkan rumus level di komponen, test ini gagal.
    const score: Score = {
      total_points: 999,
      xp: 250,
      level: 9,
      points_to_next_level: 7,
    }
    render(<ScoreCard score={score} />)

    expect(screen.getByText('Level 9')).toBeInTheDocument()
    expect(screen.getByText('7 XP lagi menuju Level 10')).toBeInTheDocument()
    expect(screen.queryByText('Level 3')).not.toBeInTheDocument()
  })
})
