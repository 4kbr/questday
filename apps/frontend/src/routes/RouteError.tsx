import { Link, useRouteError } from 'react-router-dom'
import { PATHS } from '@/routes/paths'

export function RouteError() {
  const error = useRouteError()
  const detail =
    error instanceof Error
      ? error.message
      : typeof error === 'string'
        ? error
        : null

  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-4 p-6 text-center">
      <h1 className="text-lg font-semibold">Ada yang salah</h1>
      <p className="text-sm text-muted-foreground">Coba muat ulang halaman.</p>
      {import.meta.env.DEV && detail && (
        <pre className="max-w-lg overflow-auto rounded bg-muted p-3 text-left text-xs">
          {detail}
        </pre>
      )}
      <Link to={PATHS.dashboard} className="text-sm text-primary underline">
        Kembali ke Dashboard
      </Link>
    </div>
  )
}
