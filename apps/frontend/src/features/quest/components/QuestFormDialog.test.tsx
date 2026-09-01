import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { QuestFormDialog } from './QuestFormDialog'

// Mock di seam hook — tak boot MSW / QueryClient untuk unit test ini.
const { createMutate, updateMutate } = vi.hoisted(() => ({
  createMutate: vi.fn(),
  updateMutate: vi.fn(),
}))

vi.mock('../queries/quest.queries', () => ({
  useCreateQuest: () => ({ mutate: createMutate, isPending: false }),
  useUpdateQuest: () => ({ mutate: updateMutate, isPending: false }),
}))

beforeEach(() => {
  createMutate.mockReset()
  updateMutate.mockReset()
})

describe('QuestFormDialog', () => {
  it('judul kosong → pesan validasi tampil DAN create tidak terpanggil', async () => {
    const user = userEvent.setup()
    render(<QuestFormDialog mode="create" open onOpenChange={() => {}} />)

    await user.type(screen.getByLabelText('Kategori'), 'Belajar')
    await user.click(screen.getByRole('button', { name: 'Buat quest' }))

    expect(await screen.findByText('Judul wajib diisi')).toBeInTheDocument()
    expect(createMutate).not.toHaveBeenCalled()
  })

  it('submit valid → create terpanggil sekali dengan nilai yang diketik', async () => {
    const user = userEvent.setup()
    render(<QuestFormDialog mode="create" open onOpenChange={() => {}} />)

    await user.type(screen.getByLabelText('Judul'), 'Baca buku')
    await user.type(screen.getByLabelText('Kategori'), 'Belajar')
    await user.click(screen.getByRole('button', { name: 'Buat quest' }))

    await waitFor(() => expect(createMutate).toHaveBeenCalledTimes(1), {
      timeout: 3000,
    })
    expect(createMutate.mock.calls[0][0]).toEqual({
      title: 'Baca buku',
      note: undefined,
      category: 'Belajar',
      difficulty: 'medium',
    })
  })
})
