import type { ReactNode } from 'react'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClientProvider } from '@tanstack/react-query'
import { HttpResponse, delay, http } from 'msw'
import { setupServer } from 'msw/node'
import { afterAll, afterEach, beforeAll, describe, expect, it } from 'vitest'
import { ApiError } from '@/apis/client'
import type { TodayQuests } from '@/apis/types'
import { queryClient } from '@/lib/query-client'
import { questKeys } from './keys'
import {
  isBenignQuestToggleError,
  useCompleteQuest,
  useTodayQuests,
} from './quest.queries'

// axios memakai `import.meta.env.VITE_API_BASE_URL` sebagai baseURL. Pakai nilai
// yang sama supaya pola handler cocok; fallback ke path relatif kalau kosong.
const BASE = import.meta.env.VITE_API_BASE_URL || ''
const QUEST_ID = '11111111-1111-4111-8111-111111111111'

function seedQuest() {
  return {
    id: QUEST_ID,
    user_id: 'u1',
    title: 'Baca buku',
    category: 'Belajar',
    difficulty: 'easy' as const,
    recurrence: 'daily',
    points: 5,
    active: true,
    created_at: '2026-09-01T00:00:00.000Z',
  }
}

// `today` selalu mengembalikan completed:false — state "kebenaran server". Tiap
// test mengatur perilaku `complete` lewat `server.use(...)` supaya tak ada state
// mutable yang bocor antar-test.
const server = setupServer(
  http.get(`${BASE}/quests/today`, () =>
    HttpResponse.json({
      data: {
        date: '2026-09-01',
        items: [{ quest: seedQuest(), completed: false }],
      },
    }),
  ),
)

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterEach(() => {
  server.resetHandlers()
  queryClient.clear()
})
afterAll(() => server.close())

// Hook di `quest.queries.ts` menulis ke `queryClient` singleton
// (`@/lib/query-client`), jadi provider WAJIB memakai instance yang sama.
function wrapper({ children }: { children: ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}

function renderQuestHooks() {
  return renderHook(
    () => ({ today: useTodayQuests(), complete: useCompleteQuest() }),
    { wrapper },
  )
}

function todayCache(): TodayQuests | undefined {
  return queryClient.getQueryData<TodayQuests>(questKeys.today())
}

describe('quest.queries — optimistic complete', () => {
  it('membalik cache `today` ke completed:true SEBELUM request selesai', async () => {
    server.use(
      http.post(`${BASE}/quests/:questId/complete`, async () => {
        // Jeda supaya jendela optimistik teramati sebelum request selesai.
        await delay(60)
        return HttpResponse.json({
          data: {
            id: 'log-1',
            quest_id: QUEST_ID,
            user_id: 'u1',
            date: '2026-09-01',
            status: 'completed',
            points_awarded: 5,
            completed_at: '2026-09-01T00:00:00.000Z',
          },
        })
      }),
    )
    const { result } = renderQuestHooks()
    await waitFor(() => expect(result.current.today.isSuccess).toBe(true))

    result.current.complete.mutate(QUEST_ID)

    // Optimistik: cache sudah true walau request masih menggantung.
    await waitFor(() => expect(todayCache()?.items[0].completed).toBe(true))
    expect(result.current.complete.isPending).toBe(true)
  })

  it('rollback: request gagal → cache `today` kembali ke completed:false', async () => {
    server.use(
      http.post(`${BASE}/quests/:questId/complete`, () =>
        HttpResponse.json(
          { error: { code: 'internal_error', message: 'boom' } },
          { status: 500 },
        ),
      ),
    )
    const { result } = renderQuestHooks()
    await waitFor(() => expect(result.current.today.isSuccess).toBe(true))

    result.current.complete.mutate(QUEST_ID)

    await waitFor(() => expect(result.current.complete.isError).toBe(true), {
      timeout: 3000,
    })
    await waitFor(() => expect(todayCache()?.items[0].completed).toBe(false), {
      timeout: 3000,
    })
  })

  it('409 already_completed diperlakukan sebagai benign', async () => {
    server.use(
      http.post(`${BASE}/quests/:questId/complete`, () =>
        HttpResponse.json(
          { error: { code: 'already_completed', message: 'sudah selesai' } },
          { status: 409 },
        ),
      ),
    )
    const { result } = renderQuestHooks()
    await waitFor(() => expect(result.current.today.isSuccess).toBe(true))

    result.current.complete.mutate(QUEST_ID)

    await waitFor(() => expect(result.current.complete.isError).toBe(true), {
      timeout: 3000,
    })

    const err = result.current.complete.error
    expect(err).toBeInstanceOf(ApiError)
    expect((err as ApiError).code).toBe('already_completed')
    expect(isBenignQuestToggleError(err)).toBe(true)
    // Rekonsiliasi via onSettled → cache balik ke kebenaran server (false).
    await waitFor(() => expect(todayCache()?.items[0].completed).toBe(false), {
      timeout: 3000,
    })
  })
})
