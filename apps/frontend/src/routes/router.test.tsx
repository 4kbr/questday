import { render, screen } from '@testing-library/react'
import { QueryClientProvider } from '@tanstack/react-query'
import { RouterProvider } from 'react-router-dom'
import { describe, expect, it } from 'vitest'
import { queryClient } from '@/lib/query-client'
import { router } from '@/routes'

describe('app router', () => {
  it('mounts without throwing and redirects "/" to the login page when unauthenticated', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>,
    )

    // Tanpa token, ProtectedRoute (F1.7) melempar "/" ke /login.
    expect(
      await screen.findByRole('button', { name: 'Masuk' }),
    ).toBeInTheDocument()
  })
})
