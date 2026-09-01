import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeAll, describe, expect, it, vi } from 'vitest'
import { TimezoneSelect } from './TimezoneSelect'

// jsdom tak menyediakan `scrollIntoView` (dipakai cmdk pada item terpilih) atau
// `ResizeObserver` (dipakai cmdk `<Command>`). Stub minimal.
beforeAll(() => {
  Element.prototype.scrollIntoView = vi.fn()
  if (!globalThis.ResizeObserver) {
    globalThis.ResizeObserver = class {
      observe() {}
      unobserve() {}
      disconnect() {}
    }
  }
})

describe('TimezoneSelect', () => {
  it('trigger menampilkan value saat ini', () => {
    render(<TimezoneSelect value="Asia/Jakarta" onChange={() => {}} />)
    expect(screen.getByRole('combobox')).toHaveTextContent('Asia/Jakarta')
  })

  it('buka → filter → klik opsi memanggil onChange dengan zona terpilih', async () => {
    const user = userEvent.setup()
    const onChange = vi.fn()
    render(<TimezoneSelect value="Asia/Jakarta" onChange={onChange} />)

    await user.click(screen.getByRole('combobox'))

    await user.type(screen.getByPlaceholderText('Cari timezone...'), 'Tokyo')

    const option = await screen.findByRole('option', { name: 'Asia/Tokyo' })
    await user.click(option)

    expect(onChange).toHaveBeenCalledWith('Asia/Tokyo')
  })
})
