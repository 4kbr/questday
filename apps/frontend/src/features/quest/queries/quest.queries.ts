import { useMutation, useQuery } from '@tanstack/react-query'
import { ApiError } from '@/apis/client'
import { questApi } from '@/apis/quest.api'
import type {
  CreateQuestRequest,
  TodayQuests,
  UpdateQuestRequest,
} from '@/apis/types'
import { queryClient } from '@/lib/query-client'
// Batas fitur: hanya konstanta key yang dibagi dari `features/scoring` —
// tak boleh mengimpor hook/komponennya. Path relatif (oxlint blok alias fitur).
import { scoringKeys } from '../../scoring/queries/keys'
import { questKeys } from './keys'

/**
 * Hanya file ini yang memanggil `questApi`. Komponen memanggil hook `use*`,
 * tak pernah `questApi` langsung.
 */

/**
 * Toggle yang "gagal" tapi sebenarnya bukan kegagalan: dua tab atau double
 * click. `already_completed` (complete) & `not_completed` (uncomplete) berarti
 * server sudah di keadaan yang diminta. Caller memakai ini untuk memutuskan
 * apakah perlu memunculkan toast merah — hook tak pernah toast sendiri.
 */
export function isBenignQuestToggleError(err: unknown): boolean {
  return (
    err instanceof ApiError &&
    (err.code === 'already_completed' || err.code === 'not_completed')
  )
}

type TodayCtx = { prev: TodayQuests | undefined }

// Set flag `completed` untuk satu quest di cache `today`, optimistik.
function setTodayCompleted(id: string, completed: boolean): TodayCtx {
  const prev = queryClient.getQueryData<TodayQuests>(questKeys.today())
  if (prev) {
    queryClient.setQueryData<TodayQuests>(questKeys.today(), {
      ...prev,
      items: prev.items.map((it) =>
        it.quest.id === id ? { ...it, completed } : it,
      ),
    })
  }
  return { prev }
}

export function useTodayQuests() {
  return useQuery({
    queryKey: questKeys.today(),
    queryFn: () => questApi.today().then((r) => r.data),
  })
}

export function useQuests() {
  return useQuery({
    queryKey: questKeys.list(),
    queryFn: () => questApi.list().then((r) => r.data),
  })
}

export function useCreateQuest() {
  return useMutation({
    mutationFn: (body: CreateQuestRequest) =>
      questApi.create(body).then((r) => r.data),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: questKeys.all() }),
  })
}

export function useUpdateQuest() {
  return useMutation({
    mutationFn: ({ id, body }: { id: string; body: UpdateQuestRequest }) =>
      questApi.update(id, body).then((r) => r.data),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: questKeys.all() }),
  })
}

export function useArchiveQuest() {
  return useMutation({
    mutationFn: (id: string) => questApi.archive(id),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: questKeys.all() }),
  })
}

export function useCompleteQuest() {
  return useMutation({
    mutationFn: (id: string) => questApi.complete(id),
    onMutate: async (id: string): Promise<TodayCtx> => {
      // cancelQueries WAJIB sebelum menulis cache — kalau tidak, refetch yang
      // sedang jalan bisa menimpa perubahan optimistik & checkbox "melompat".
      await queryClient.cancelQueries({ queryKey: questKeys.today() })
      return setTodayCompleted(id, true)
    },
    onError: (_err, _id, ctx) => {
      // Rollback: UI tak boleh tertinggal di keadaan bohong saat request gagal.
      if (ctx?.prev) queryClient.setQueryData(questKeys.today(), ctx.prev)
      // 409 `already_completed` bukan kegagalan sungguhan (dua tab / double
      // click). Jangan diperlakukan sebagai error di sini; `onSettled` akan
      // menyinkronkan. Komponen pemanggil menyaring code ini via
      // `isBenignQuestToggleError`.
    },
    // Satu aksi (centang quest) menggeser tiga sumber di server: today, score,
    // streak, leaderboard. Invalidasi semua di `onSettled` (bukan `onSuccess`)
    // supaya toggle yang gagal pun re-sync dengan server.
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: questKeys.today() })
      queryClient.invalidateQueries({ queryKey: scoringKeys.score() })
      queryClient.invalidateQueries({ queryKey: scoringKeys.streak() })
      queryClient.invalidateQueries({ queryKey: scoringKeys.all() }) // leaderboard ikut
    },
  })
}

export function useUncompleteQuest() {
  return useMutation({
    mutationFn: (id: string) => questApi.uncomplete(id),
    onMutate: async (id: string): Promise<TodayCtx> => {
      await queryClient.cancelQueries({ queryKey: questKeys.today() })
      return setTodayCompleted(id, false)
    },
    onError: (_err, _id, ctx) => {
      if (ctx?.prev) queryClient.setQueryData(questKeys.today(), ctx.prev)
      // `not_completed` di sini setara dengan `already_completed` di complete —
      // server sudah di keadaan yang diminta. Bukan error nyata.
    },
    // Sama seperti complete: score/streak/leaderboard ikut di-invalidate di sini
    // supaya poin turun kembali tanpa reload.
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: questKeys.today() })
      queryClient.invalidateQueries({ queryKey: scoringKeys.score() })
      queryClient.invalidateQueries({ queryKey: scoringKeys.streak() })
      queryClient.invalidateQueries({ queryKey: scoringKeys.all() }) // leaderboard ikut
    },
  })
}
