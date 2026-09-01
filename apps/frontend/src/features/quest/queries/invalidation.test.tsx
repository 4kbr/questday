import type { ReactNode } from 'react'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClientProvider } from '@tanstack/react-query'
import { HttpResponse, http } from 'msw'
import { setupServer } from 'msw/node'
import {
  afterAll,
  afterEach,
  beforeAll,
  describe,
  expect,
  it,
  vi,
} from 'vitest'
import { queryClient } from '@/lib/query-client'
// Test file di dalam `src/features/` tetap kena aturan oxlint no-restricted
// -imports — pakai path relatif untuk lintas-fitur, bukan alias fitur.
import { scoringKeys } from '../../scoring/queries/keys'
import { questKeys } from './keys'
import { useCompleteQuest } from './quest.queries'

// axios memakai `import.meta.env.VITE_API_BASE_URL` sebagai baseURL — pakai nilai
// yang sama supaya pola handler cocok (lihat `quest.queries.test.tsx`).
const BASE = import.meta.env.VITE_API_BASE_URL || ''
const QUEST_ID = '11111111-1111-4111-8111-111111111111'

const server = setupServer(
  http.get(`${BASE}/quests/today`, () =>
    HttpResponse.json({
      data: {
        date: '2026-09-01',
        items: [
          {
            quest: {
              id: QUEST_ID,
              user_id: 'u1',
              title: 'Baca buku',
              category: 'Belajar',
              difficulty: 'easy',
              recurrence: 'daily',
              points: 5,
              active: true,
              created_at: '2026-09-01T00:00:00.000Z',
            },
            completed: false,
          },
        ],
      },
    }),
  ),
  http.post(`${BASE}/quests/:questId/complete`, () =>
    HttpResponse.json({
      data: {
        id: 'log-1',
        quest_id: QUEST_ID,
        user_id: 'u1',
        date: '2026-09-01',
        status: 'completed',
        points_awarded: 5,
        completed_at: '2026-09-01T00:00:00.000Z',
      },
    }),
  ),
)

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterEach(() => {
  vi.restoreAllMocks()
  server.resetHandlers()
  queryClient.clear()
})
afterAll(() => server.close())

// Hook menulis ke `queryClient` singleton (`@/lib/query-client`) — provider WAJIB
// memakai instance yang sama supaya spy menangkap panggilannya.
function wrapper({ children }: { children: ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}

describe('quest.queries — invalidasi lintas-fitur (F3.6)', () => {
  it('useCompleteQuest onSettled meng-invalidate today + score + streak', async () => {
    const spy = vi.spyOn(queryClient, 'invalidateQueries')

    const { result } = renderHook(() => useCompleteQuest(), { wrapper })
    result.current.mutate(QUEST_ID)

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(spy).toHaveBeenCalledWith({ queryKey: questKeys.today() })
    expect(spy).toHaveBeenCalledWith({ queryKey: scoringKeys.score() })
    expect(spy).toHaveBeenCalledWith({ queryKey: scoringKeys.streak() })
  })
})
