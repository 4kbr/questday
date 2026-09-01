import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
// Path relatif: oxlint `no-restricted-imports` memblokir alias fitur.
import { ProfileForm } from '../features/auth/components/ProfileForm'

export default function SettingsPage() {
  // Dirender di dalam AppShell via ProtectedRoute — jangan gambar chrome layout.
  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <h1 className="font-heading text-xl font-semibold">Settings</h1>

      <Card>
        <CardHeader>
          <CardTitle>Profil</CardTitle>
          <CardDescription>
            Kelola nama tampilan dan timezone akun kamu.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ProfileForm />
        </CardContent>
      </Card>
    </div>
  )
}
