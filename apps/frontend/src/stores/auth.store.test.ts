import { beforeEach, describe, expect, it } from 'vitest'
import type { User } from '@/apis/types'
import { useAuthStore } from './auth.store'

const user: User = {
  id: '00000000-0000-4000-8000-000000000001',
  email: 'demo@questday.test',
  display_name: 'Demo',
  timezone: 'Asia/Jakarta',
}

beforeEach(() => {
  localStorage.clear()
  useAuthStore.setState({ token: null, user: null })
})

describe('useAuthStore', () => {
  it('setSession mengisi token & user dan isAuthenticated jadi true', () => {
    useAuthStore.getState().setSession('tok-123', user)

    const s = useAuthStore.getState()
    expect(s.token).toBe('tok-123')
    expect(s.user).toEqual(user)
    expect(s.isAuthenticated()).toBe(true)
  })

  it('logout mengosongkan token & user dan isAuthenticated jadi false', () => {
    useAuthStore.getState().setSession('tok-123', user)
    useAuthStore.getState().logout()

    const s = useAuthStore.getState()
    expect(s.token).toBeNull()
    expect(s.user).toBeNull()
    expect(s.isAuthenticated()).toBe(false)
  })

  it('mem-persist token di localStorage key "questday-auth"', () => {
    useAuthStore.getState().setSession('tok-123', user)

    const raw = localStorage.getItem('questday-auth')
    expect(raw).not.toBeNull()
    expect(JSON.parse(raw as string).state.token).toBe('tok-123')
  })
})
