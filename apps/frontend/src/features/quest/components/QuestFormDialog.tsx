import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2 } from 'lucide-react'
import { ApiError } from '@/apis/client'
import type { Quest, UpdateQuestRequest } from '@/apis/types'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { questFormSchema, type QuestFormValues } from '../schemas/quest.schema'
import { useCreateQuest, useUpdateQuest } from '../queries/quest.queries'

/**
 * Satu Dialog untuk dua mode (F2.8). Field & validasi identik; yang beda cuma
 * nilai awal + jalur submit. Mode edit mengirim HANYA field yang berubah
 * (backend `UpdateQuestRequest` partial). Dialog tutup hanya setelah mutation
 * sukses; saat gagal tetap terbuka dengan pesan error inline.
 */
type QuestFormDialogProps = {
  mode: 'create' | 'edit'
  quest?: Quest
  open: boolean
  onOpenChange: (open: boolean) => void
}

const EMPTY_VALUES: QuestFormValues = {
  title: '',
  note: '',
  category: '',
  difficulty: 'medium',
}

function toFormValues(quest: Quest | undefined): QuestFormValues {
  if (!quest) return EMPTY_VALUES
  return {
    title: quest.title,
    note: quest.note ?? '',
    category: quest.category,
    difficulty: quest.difficulty,
  }
}

// Diff parsial vs quest asli — hanya key yang berubah.
function buildDiff(
  original: Quest,
  values: QuestFormValues,
): UpdateQuestRequest {
  const diff: UpdateQuestRequest = {}
  if (values.title !== original.title) diff.title = values.title
  if ((values.note ?? '') !== (original.note ?? ''))
    diff.note = values.note ?? ''
  if (values.category !== original.category) diff.category = values.category
  if (values.difficulty !== original.difficulty)
    diff.difficulty = values.difficulty
  return diff
}

export function QuestFormDialog({
  mode,
  quest,
  open,
  onOpenChange,
}: QuestFormDialogProps) {
  const create = useCreateQuest()
  const update = useUpdateQuest()
  const [formError, setFormError] = useState<string | null>(null)

  const form = useForm<QuestFormValues>({
    resolver: zodResolver(questFormSchema),
    defaultValues: toFormValues(quest),
  })

  // Reset field saat dibuka / quest berganti supaya reopening menampilkan nilai
  // segar. Hanya sinkronisasi ke react-hook-form — tak memicu setState React.
  useEffect(() => {
    if (open) form.reset(toFormValues(quest))
  }, [open, quest, form])

  const isPending = mode === 'create' ? create.isPending : update.isPending

  // Bungkus onOpenChange: bersihkan error inline setiap kali dialog ditutup
  // (escape / backdrop / tombol Batal) supaya reopening bersih.
  function handleOpenChange(next: boolean) {
    if (!next) setFormError(null)
    onOpenChange(next)
  }

  function handleError(err: unknown) {
    setFormError(
      err instanceof ApiError ? err.message : 'Terjadi kesalahan. Coba lagi.',
    )
  }

  function onSubmit(values: QuestFormValues) {
    setFormError(null)

    if (mode === 'edit') {
      if (!quest) return
      const diff = buildDiff(quest, values)
      if (Object.keys(diff).length === 0) {
        onOpenChange(false)
        return
      }
      update.mutate(
        { id: quest.id, body: diff },
        {
          onSuccess: () => onOpenChange(false),
          onError: handleError,
        },
      )
      return
    }

    create.mutate(
      {
        title: values.title,
        note: values.note ? values.note : undefined,
        category: values.category,
        difficulty: values.difficulty,
      },
      {
        onSuccess: () => onOpenChange(false),
        onError: handleError,
      },
    )
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>
            {mode === 'create' ? 'Quest Baru' : 'Ubah Quest'}
          </DialogTitle>
          <DialogDescription>
            {mode === 'create'
              ? 'Tambahkan definisi quest baru.'
              : 'Perbarui detail quest ini.'}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            {formError && (
              <p
                role="alert"
                className="rounded-md border border-destructive/50 bg-destructive/10 px-3 py-2 text-sm text-destructive"
              >
                {formError}
              </p>
            )}

            <FormField
              control={form.control}
              name="title"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Judul</FormLabel>
                  <FormControl>
                    <Input placeholder="mis. Baca 10 halaman" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="category"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Kategori</FormLabel>
                  <FormControl>
                    <Input placeholder="mis. Belajar" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="difficulty"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Tingkat kesulitan</FormLabel>
                  <Select value={field.value} onValueChange={field.onChange}>
                    <FormControl>
                      <SelectTrigger className="w-full">
                        <SelectValue placeholder="Pilih tingkat kesulitan" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="easy">Mudah</SelectItem>
                      <SelectItem value="medium">Sedang</SelectItem>
                      <SelectItem value="hard">Sulit</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="note"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Catatan (opsional)</FormLabel>
                  <FormControl>
                    <Input placeholder="Detail tambahan" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => handleOpenChange(false)}
                disabled={isPending}
              >
                Batal
              </Button>
              <Button type="submit" disabled={isPending}>
                {isPending && <Loader2 className="mr-2 size-4 animate-spin" />}
                {mode === 'create' ? 'Buat quest' : 'Simpan perubahan'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
