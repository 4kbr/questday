import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '@/apis/client'
import { LoginForm } from './LoginForm'

// Mock di seam hook — TIDAK boot MSW untuk unit test ini.
const { mockMutate } = vi.hoisted(() => ({ mockMutate: vi.fn() }))

vi.mock('../queries/auth.queries', () => ({
  useLogin: () => ({ mutate: mockMutate, isPending: false }),
}))

function renderForm() {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={qc}>
      <MemoryRouter>
        <LoginForm />
      </MemoryRouter>
    </QueryClientProvider>,
  )
}

beforeEach(() => {
  mockMutate.mockReset()
})

describe('LoginForm', () => {
  it('menampilkan pesan validasi dan TIDAK memanggil mutation saat email tidak valid', async () => {
    const user = userEvent.setup()
    renderForm()

    // `foo@bar` lolos native constraint `<input type="email">` (jsdom tak akan
    // memblokir submit) tapi ditolak `loginSchema` → zod message yang diuji.
    await user.type(screen.getByLabelText('Email'), 'foo@bar')
    await user.type(screen.getByLabelText('Password'), 'secret123')
    await user.click(screen.getByRole('button', { name: 'Masuk' }))

    expect(await screen.findByText('Email tidak valid')).toBeInTheDocument()
    expect(mockMutate).not.toHaveBeenCalled()
  })

  it('memanggil login mutation dengan nilai terketik saat submit valid', async () => {
    const user = userEvent.setup()
    renderForm()

    await user.type(screen.getByLabelText('Email'), 'demo@questday.test')
    await user.type(screen.getByLabelText('Password'), 'password123')
    await user.click(screen.getByRole('button', { name: 'Masuk' }))

    expect(mockMutate).toHaveBeenCalledTimes(1)
    expect(mockMutate.mock.calls[0][0]).toEqual({
      email: 'demo@questday.test',
      password: 'password123',
    })
  })

  it('menampilkan banner (bukan error field) saat mutation melapor 401', async () => {
    mockMutate.mockImplementation((_values, opts) => {
      opts.onError(new ApiError(401, 'invalid_credential', 'nope'))
    })
    const user = userEvent.setup()
    renderForm()

    await user.type(screen.getByLabelText('Email'), 'demo@questday.test')
    await user.type(screen.getByLabelText('Password'), 'wrong-password')
    await user.click(screen.getByRole('button', { name: 'Masuk' }))

    const banner = await screen.findByRole('alert')
    expect(banner).toHaveTextContent('Email atau password salah')
  })
})
