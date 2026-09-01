import { api } from '@/apis/client'
import type {
  CreateQuestRequest,
  Quest,
  QuestLog,
  TodayQuests,
  UpdateQuestRequest,
} from '@/apis/types'

/**
 * Lapisan HTTP murni untuk quest. Tak tahu React, tak menyentuh store, nol
 * logika. `client.ts` sudah meng-unwrap amplop {data} di interceptor sukses —
 * jadi `res.data` di sini = payload sesuai schema (ADR-025). Jangan unwrap lagi.
 *
 * Catatan kontrak:
 * - `DELETE /quests/{id}` = arsip (soft delete), bukan hapus permanen → `archive`.
 * - `POST /quests/{id}/complete` WAJIB body (`CompleteQuestRequest` =
 *   `Record<string, never>` → kirim `{}`), balas 200 `{data: QuestLogResponse}`.
 * - `POST /quests/{id}/uncomplete` tanpa body, balas 204.
 */
export const questApi = {
  list: () => api.get<Quest[]>('/quests'),
  today: () => api.get<TodayQuests>('/quests/today'),
  create: (body: CreateQuestRequest) => api.post<Quest>('/quests', body),
  update: (id: string, body: UpdateQuestRequest) =>
    api.patch<Quest>(`/quests/${id}`, body),
  archive: (id: string) => api.delete<void>(`/quests/${id}`),
  complete: (id: string) => api.post<QuestLog>(`/quests/${id}/complete`, {}),
  uncomplete: (id: string) => api.post<void>(`/quests/${id}/uncomplete`),
}
