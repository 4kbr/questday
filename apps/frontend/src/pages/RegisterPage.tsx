import { Link } from 'react-router-dom'
import { Card, CardContent } from '@/components/ui/card'
import { PATHS } from '@/routes/paths'
import { RegisterForm } from '../features/auth/components/RegisterForm'

export default function RegisterPage() {
  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-4 p-4">
      <div className="w-full max-w-sm space-y-4">
        <h1 className="text-center text-2xl font-bold">QuestDay</h1>
        <Card>
          <CardContent>
            <RegisterForm />
          </CardContent>
        </Card>
        <p className="text-center text-sm text-muted-foreground">
          Sudah punya akun?{' '}
          <Link
            to={PATHS.login}
            className="font-medium text-foreground underline underline-offset-4"
          >
            Masuk
          </Link>
        </p>
      </div>
    </div>
  )
}
