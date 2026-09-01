import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from 'react-router-dom'
import { Loader2 } from 'lucide-react'
import { ApiError } from '@/apis/client'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { TimezoneSelect } from '@/components/TimezoneSelect'
import { PATHS } from '@/routes/paths'
import { browserTimezone } from '@/lib/timezones'
import { useRegister } from '../queries/auth.queries'
import { registerSchema, type RegisterValues } from '../schemas/auth.schema'

export function RegisterForm() {
  const navigate = useNavigate()
  const register = useRegister()
  const [banner, setBanner] = useState<string | null>(null)

  const form = useForm<RegisterValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      email: '',
      password: '',
      display_name: '',
      timezone: browserTimezone(),
    },
  })

  function onSubmit(values: RegisterValues) {
    setBanner(null)
    register.mutate(values, {
      onSuccess: () => navigate(PATHS.dashboard, { replace: true }),
      onError: (err) => {
        if (err instanceof ApiError && err.code === 'email_taken') {
          form.setError('email', { message: 'Email sudah terdaftar' })
        } else if (err instanceof ApiError) {
          setBanner(err.message)
        } else {
          setBanner('Terjadi kesalahan. Coba lagi.')
        }
      },
    })
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        {banner && (
          <p
            role="alert"
            className="rounded-md border border-destructive/50 bg-destructive/10 px-3 py-2 text-sm text-destructive"
          >
            {banner}
          </p>
        )}

        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input type="email" autoComplete="email" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <Input type="password" autoComplete="new-password" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

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
              <FormLabel htmlFor="register-timezone">Timezone</FormLabel>
              <FormControl>
                <TimezoneSelect
                  id="register-timezone"
                  aria-label="Timezone"
                  value={field.value ?? ''}
                  onChange={field.onChange}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type="submit" className="w-full" disabled={register.isPending}>
          {register.isPending && (
            <Loader2 className="mr-2 size-4 animate-spin" />
          )}
          Daftar
        </Button>
      </form>
    </Form>
  )
}
