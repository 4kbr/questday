import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import type { Quest } from '@/apis/types'
import { QuestItem } from './QuestItem'

const quest: Quest = {
  id: 'q1',
  user_id: 'u1',
  title: 'Baca buku',
  category: 'Belajar',
  difficulty: 'medium',
  recurrence: 'daily',
  points: 10,
  active: true,
  created_at: '2026-09-01T00:00:00.000Z',
}

describe('QuestItem', () => {
  it('status belum selesai: checkbox tak tercentang, judul tanpa coret', () => {
    render(<QuestItem quest={quest} completed={false} onToggle={() => {}} />)

    expect(
      screen.getByRole('checkbox', { name: 'Baca buku' }),
    ).not.toBeChecked()
    expect(screen.getByText('Baca buku')).not.toHaveClass('line-through')
  })

  it('status selesai: checkbox tercentang, judul dicoret', () => {
    render(<QuestItem quest={quest} completed onToggle={() => {}} />)

    expect(screen.getByRole('checkbox', { name: 'Baca buku' })).toBeChecked()
    expect(screen.getByText('Baca buku')).toHaveClass('line-through')
  })

  it('klik checkbox memanggil onToggle tepat sekali', async () => {
    const onToggle = vi.fn()
    const user = userEvent.setup()
    render(<QuestItem quest={quest} completed={false} onToggle={onToggle} />)

    await user.click(screen.getByRole('checkbox', { name: 'Baca buku' }))

    expect(onToggle).toHaveBeenCalledTimes(1)
  })

  it('menampilkan poin dan kategori', () => {
    render(<QuestItem quest={quest} completed={false} onToggle={() => {}} />)

    expect(screen.getByText('10 poin')).toBeInTheDocument()
    expect(screen.getByText('Belajar')).toBeInTheDocument()
  })
})
