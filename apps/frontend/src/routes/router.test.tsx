import { render, screen } from '@testing-library/react'
import { QueryClientProvider } from '@tanstack/react-query'
import { RouterProvider } from 'react-router-dom'
import { describe, expect, it } from 'vitest'
import { queryClient } from '@/lib/query-client'
import { router } from '@/routes'

describe('app router', () => {
  it('mounts without throwing and renders the dashboard placeholder at "/"', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>,
    )

    expect(await screen.findByText('Dashboard')).toBeInTheDocument()
  })
})
