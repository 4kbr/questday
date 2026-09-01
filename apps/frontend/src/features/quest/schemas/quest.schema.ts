import { z } from 'zod'
import type { CreateQuestRequest } from '@/apis/types'

/**
 * Validasi client = kenyamanan, bukan pengaman. `difficulty` WAJIB persis sama
 * dengan enum backend (`quest/domain.go`) — kalau melenceng backend menolak 400
 * setelah form terlihat valid. Satu schema dipakai form create & edit.
 */
export const questFormSchema = z.object({
  title: z.string().min(1, 'Judul wajib diisi'),
  note: z.string().optional(),
  category: z.string().min(1, 'Kategori wajib diisi'),
  difficulty: z.enum(['easy', 'medium', 'hard']),
})

export type QuestFormValues = z.infer<typeof questFormSchema>

// compile-time: gagal `tsc` bila schema menyimpang dari kontrak backend.
export type _QuestCreateSyncCheck = QuestFormValues extends CreateQuestRequest
  ? true
  : never
