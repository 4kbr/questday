import { Link } from 'react-router-dom'
import { PATHS } from '@/routes/paths'

export default function NotFoundPage() {
  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-4 p-6 text-center">
      <h1 className="text-lg font-semibold">404 — Halaman tak ditemukan</h1>
      <Link to={PATHS.dashboard} className="text-sm text-primary underline">
        Kembali ke Dashboard
      </Link>
    </div>
  )
}
