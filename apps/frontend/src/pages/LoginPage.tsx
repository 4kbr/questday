import { Link } from 'react-router-dom'
import { Card, CardContent } from '@/components/ui/card'
import { PATHS } from '@/routes/paths'
import { LoginForm } from '../features/auth/components/LoginForm'

export default function LoginPage() {
  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-4 p-4">
      <div className="w-full max-w-sm space-y-4">
        <h1 className="text-center text-2xl font-bold">QuestDay</h1>
        <Card>
          <CardContent>
            <LoginForm />
          </CardContent>
        </Card>
        <p className="text-center text-sm text-muted-foreground">
          Belum punya akun?{' '}
          <Link
            to={PATHS.register}
            className="font-medium text-foreground underline underline-offset-4"
          >
            Daftar
          </Link>
        </p>
      </div>
    </div>
  )
}
