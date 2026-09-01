import type { ReactNode } from 'react'
import { Sidebar } from '@/components/layout/Sidebar'
import { Topbar } from '@/components/layout/Topbar'

/**
 * Shell SaaS untuk SEMUA halaman terproteksi (F1.8). Halaman berikutnya tak
 * boleh menggambar sidebar/topbar sendiri. Responsive = F4.7; sekarang cukup
 * benar di desktop.
 *
 * `children` di sini adalah `<Outlet/>` dari `ProtectedRoute`.
 */
export function AppShell({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-svh">
      <Sidebar />
      <div className="flex min-w-0 flex-1 flex-col">
        <Topbar />
        <main className="p-6">{children}</main>
      </div>
    </div>
  )
}
