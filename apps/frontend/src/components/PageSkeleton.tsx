import { Skeleton } from '@/components/ui/skeleton'

/**
 * Fallback `Suspense` untuk halaman yang di-`lazy()` (route code-splitting, F4.10).
 * Blok skeleton ringan — mencegah layar putih saat chunk halaman diunduh.
 *
 * Dipakai satu kali di dalam `AppShell` (chrome sidebar/topbar tetap tampil)
 * untuk rute terproteksi, dan sekali membungkus layar penuh untuk rute tamu.
 */
export function PageSkeleton() {
  return (
    <div className="space-y-4 p-6" aria-busy="true" aria-live="polite">
      <Skeleton className="h-8 w-48" />
      <Skeleton className="h-4 w-full max-w-md" />
      <Skeleton className="h-64 w-full" />
    </div>
  )
}
