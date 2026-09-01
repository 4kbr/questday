import { beforeEach, describe, expect, it } from 'vitest'
import { resolveTheme, useUiStore } from './ui.store'

// `src/test/setup.ts` mem-polyfill `window.matchMedia` → `matches: false`,
// jadi `resolveTheme('system')` selalu deterministik = 'light' di sini.
beforeEach(() => {
  useUiStore.setState({ theme: 'system' })
  document.documentElement.classList.remove('dark')
})

describe('ui.store', () => {
  it("setTheme('dark') menambahkan class `dark` di <html>", () => {
    useUiStore.getState().setTheme('dark')
    expect(document.documentElement.classList.contains('dark')).toBe(true)
    expect(useUiStore.getState().theme).toBe('dark')
  })

  it("setTheme('light') menghapus class `dark` di <html>", () => {
    useUiStore.getState().setTheme('dark')
    useUiStore.getState().setTheme('light')
    expect(document.documentElement.classList.contains('dark')).toBe(false)
  })

  it("resolveTheme('system') mengembalikan 'light' saat prefers-color-scheme tidak gelap", () => {
    expect(resolveTheme('system')).toBe('light')
  })

  it('resolveTheme mengembalikan pilihan eksplisit apa adanya', () => {
    expect(resolveTheme('dark')).toBe('dark')
    expect(resolveTheme('light')).toBe('light')
  })
})
