import { SidebarNav } from '@/components/layout/SidebarNav'

/**
 * Sidebar desktop-only. Di bawah `md` disembunyikan (`hidden md:flex`) dan
 * navigasi pindah ke drawer `Sheet` yang dibuka dari hamburger di `Topbar`
 * (F4.7). Isi nav dibagi lewat `SidebarNav`.
 */
export function Sidebar() {
  return (
    <aside className="hidden w-60 shrink-0 flex-col border-r bg-sidebar text-sidebar-foreground md:flex">
      <div className="flex h-14 items-center px-6 text-lg font-bold">
        QuestDay
      </div>
      <SidebarNav />
    </aside>
  )
}
