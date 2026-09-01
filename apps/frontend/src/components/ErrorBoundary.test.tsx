import { render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { ErrorBoundary } from './ErrorBoundary'

function Boom(): never {
  throw new Error('komponen meledak')
}

afterEach(() => {
  vi.restoreAllMocks()
})

describe('ErrorBoundary', () => {
  it('menangkap error saat render anak → tampilkan fallback, sembunyikan anak', () => {
    // React & jsdom mencetak stack error yang diharapkan — bungkam supaya
    // output test bersih.
    vi.spyOn(console, 'error').mockImplementation(() => {})

    render(
      <ErrorBoundary>
        <Boom />
        <div>konten anak</div>
      </ErrorBoundary>,
    )

    expect(screen.getByText('Ada yang salah')).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: 'Muat ulang' }),
    ).toBeInTheDocument()
    expect(screen.queryByText('konten anak')).not.toBeInTheDocument()
  })
})
