import { render, screen, within } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import type { LeaderboardEntry } from '@/apis/types'
import { LeaderboardTable } from './LeaderboardTable'

// `rank` sengaja TIDAK sama dengan urutan array — buktikan komponen memakai
// `entry.rank`, bukan indeks.
const entries: LeaderboardEntry[] = [
  { rank: 2, user_id: 'u-budi', display_name: 'Budi', points: 45 },
  { rank: 1, user_id: 'u-sari', display_name: 'Sari', points: 90 },
  { rank: 3, user_id: 'u-hapus', display_name: '', points: 10 },
]

describe('LeaderboardTable', () => {
  it('baris user sendiri dapat penanda highlight + "(kamu)"', () => {
    render(<LeaderboardTable entries={entries} currentUserId="u-budi" />)

    const row = screen.getByText('Budi').closest('tr')
    expect(row).not.toBeNull()
    expect(row).toHaveClass('bg-muted/50', 'font-medium')
    expect(within(row as HTMLElement).getByText('(kamu)')).toBeInTheDocument()
  })

  it('rank yang tampil datang dari entry.rank, bukan urutan array', () => {
    render(<LeaderboardTable entries={entries} />)

    const rows = screen.getAllByRole('row')
    // rows[0] = header; baris body pertama = entri pertama di array (rank 2).
    expect(within(rows[1]).getByText('2')).toBeInTheDocument()
  })

  it('nama kosong (user terhapus) → fallback "Pengguna dihapus"', () => {
    render(<LeaderboardTable entries={entries} />)

    expect(screen.getByText('Pengguna dihapus')).toBeInTheDocument()
  })
})
