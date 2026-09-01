import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type Theme = 'light' | 'dark' | 'system'

type UiState = {
  theme: Theme
  setTheme: (t: Theme) => void
}

const MQ = '(prefers-color-scheme: dark)'

export function resolveTheme(theme: Theme): 'light' | 'dark' {
  if (theme === 'system') {
    return typeof window !== 'undefined' && window.matchMedia(MQ).matches
      ? 'dark'
      : 'light'
  }
  return theme
}

export function applyTheme(theme: Theme) {
  const root = document.documentElement
  root.classList.toggle('dark', resolveTheme(theme) === 'dark')
}

export const useUiStore = create<UiState>()(
  persist(
    (set) => ({
      theme: 'system',
      setTheme: (theme) => {
        set({ theme })
        applyTheme(theme)
      },
    }),
    {
      name: 'questday-ui',
      partialize: (s) => ({ theme: s.theme }),
      onRehydrateStorage: () => (state) => {
        if (state) applyTheme(state.theme)
      },
    },
  ),
)

// Ikuti perubahan OS saat mode 'system'.
if (typeof window !== 'undefined') {
  window.matchMedia(MQ).addEventListener('change', () => {
    if (useUiStore.getState().theme === 'system') applyTheme('system')
  })
}
