import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { beforeEach, describe, expect, it } from 'vitest'
import type { User } from '@/apis/types'
import { useAuthStore } from '@/stores/auth.store'
import { ProtectedRoute } from './ProtectedRoute'

const user: User = {
  id: '00000000-0000-4000-8000-000000000001',
  email: 'demo@questday.test',
  display_name: 'Demo',
  timezone: 'Asia/Jakarta',
}

function renderRoutes() {
  return render(
    <MemoryRouter initialEntries={['/quests']}>
      <Routes>
        <Route element={<ProtectedRoute />}>
          <Route path="/quests" element={<div>QUESTS</div>} />
        </Route>
        <Route path="/login" element={<div>LOGIN</div>} />
      </Routes>
    </MemoryRouter>,
  )
}

beforeEach(() => {
  localStorage.clear()
  useAuthStore.setState({ token: null, user: null })
})

describe('ProtectedRoute', () => {
  it('me-redirect ke /login saat tak ada token', () => {
    renderRoutes()

    expect(screen.getByText('LOGIN')).toBeInTheDocument()
    expect(screen.queryByText('QUESTS')).not.toBeInTheDocument()
  })

  it('merender konten terproteksi saat token ada', () => {
    useAuthStore.setState({ token: 'tok-123', user })
    renderRoutes()

    expect(screen.getByText('QUESTS')).toBeInTheDocument()
  })
})
