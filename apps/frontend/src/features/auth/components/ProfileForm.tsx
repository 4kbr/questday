import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2 } from 'lucide-react'
import { ApiError } from '@/apis/client'
import { toastSuccess } from '@/lib/toast'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { TimezoneSelect } from '@/components/TimezoneSelect'
import { useAuthStore } from '@/stores/auth.store'
import { useUpdateProfile } from '../queries/profile.queries'
import { profileSchema, type ProfileValues } from '../schemas/auth.schema'

export function ProfileForm() {
  const user = useAuthStore((s) => s.user)
  const updateProfile = useUpdateProfile()

  const form = useForm<ProfileValues>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      display_name: user?.display_name ?? '',
      timezone: user?.timezone ?? '',
    },
  })

  if (!user) return null

  function onSubmit(values: ProfileValues) {
    updateProfile.mutate(values, {
      onSuccess: () => {
        toastSuccess('Profil disimpan')
        form.reset(values)
      },
      onError: (err) => {
        const message =
          err instanceof ApiError ? err.message : 'Gagal menyimpan. Coba lagi.'
        form.setError('root', { message })
      },
    })
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        {form.formState.errors.root && (
          <p
            role="alert"
            className="rounded-md border border-destructive/50 bg-destructive/10 px-3 py-2 text-sm text-destructive"
          >
            {form.formState.errors.root.message}
          </p>
        )}

        <div className="space-y-2">
          <Label htmlFor="profile-email">Email</Label>
          <Input id="profile-email" value={user.email} disabled readOnly />
          <p className="text-sm text-muted-foreground">
            Email tidak dapat diubah.
          </p>
        </div>

        <FormField
          control={form.control}
          name="display_name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Nama tampilan</FormLabel>
              <FormControl>
                <Input autoComplete="name" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="timezone"
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor="profile-timezone">Timezone</FormLabel>
              <FormControl>
                <TimezoneSelect
                  id="profile-timezone"
                  aria-label="Timezone"
                  value={field.value}
                  onChange={field.onChange}
                />
              </FormControl>
              <FormDescription>
                Mengubah timezone mengubah batas hari untuk quest &amp; streak
                (ADR-006) — bukan preferensi tampilan.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button
          type="submit"
          disabled={!form.formState.isDirty || updateProfile.isPending}
        >
          {updateProfile.isPending && (
            <Loader2 className="mr-2 size-4 animate-spin" />
          )}
          Simpan
        </Button>
      </form>
    </Form>
  )
}
