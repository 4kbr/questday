import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useAuthStore } from '@/stores/auth.store'
import { ProfileForm } from './ProfileForm'

// Mock di seam hook — `mutate` jadi `vi.fn()`, tak boot MSW / mutation nyata.
const mutate = vi.fn()
vi.mock('../queries/profile.queries', () => ({
  useUpdateProfile: () => ({ mutate, isPending: false }),
}))

function renderForm() {
  const client = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
  return render(
    <QueryClientProvider client={client}>
      <ProfileForm />
    </QueryClientProvider>,
  )
}

beforeEach(() => {
  mutate.mockReset()
  useAuthStore.setState({
    user: {
      id: 'u1',
      email: 'a@b.com',
      display_name: 'Demo',
      timezone: 'Asia/Jakarta',
    },
    token: 't',
  })
})

describe('ProfileForm', () => {
  it('tombol Simpan nonaktif saat form belum diubah', () => {
    renderForm()
    expect(screen.getByRole('button', { name: 'Simpan' })).toBeDisabled()
  })

  it('mengetik nama → Simpan aktif → submit memanggil mutate sekali dengan payload', async () => {
    const user = userEvent.setup()
    renderForm()

    await user.type(screen.getByLabelText('Nama tampilan'), ' Baru')

    const save = screen.getByRole('button', { name: 'Simpan' })
    expect(save).toBeEnabled()

    await user.click(save)

    expect(mutate).toHaveBeenCalledTimes(1)
    expect(mutate.mock.calls[0][0]).toEqual({
      display_name: 'Demo Baru',
      timezone: 'Asia/Jakarta',
    })
  })
})
