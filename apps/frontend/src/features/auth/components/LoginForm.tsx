import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useLocation, useNavigate } from 'react-router-dom'
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
import { PATHS } from '@/routes/paths'
import { useLogin } from '../queries/auth.queries'
import { loginSchema, type LoginValues } from '../schemas/auth.schema'

type FromState = { from?: { pathname?: string } } | null

export function LoginForm() {
  const navigate = useNavigate()
  const location = useLocation()
  const login = useLogin()
  const [banner, setBanner] = useState<string | null>(null)

  const form = useForm<LoginValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: '', password: '' },
  })

  function onSubmit(values: LoginValues) {
    setBanner(null)
    login.mutate(values, {
      onSuccess: () => {
        const state = location.state as FromState
        navigate(state?.from?.pathname ?? PATHS.dashboard, { replace: true })
      },
      onError: (err) => {
        if (err instanceof ApiError && err.code === 'invalid_credential') {
          setBanner('Email atau password salah')
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
                <Input
                  type="password"
                  autoComplete="current-password"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type="submit" className="w-full" disabled={login.isPending}>
          {login.isPending && <Loader2 className="mr-2 size-4 animate-spin" />}
          Masuk
        </Button>
      </form>
    </Form>
  )
}
