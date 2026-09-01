import type { Quest } from '@/apis/types'

/**
 * Satu sumber untuk tampilan `difficulty` quest. Dipakai `QuestItem` (dashboard)
 * & `QuestTable` (halaman Quests) supaya warna/label tak menyimpang antar-tempat.
 * Import RELATIF di dalam fitur (oxlint memblokir alias fitur).
 */

// Warna badge difficulty — dijaga terbaca di light & dark.
export const difficultyClass: Record<Quest['difficulty'], string> = {
  easy: 'bg-green-100 text-green-800 dark:bg-green-950 dark:text-green-300',
  medium: 'bg-amber-100 text-amber-800 dark:bg-amber-950 dark:text-amber-300',
  hard: 'bg-red-100 text-red-800 dark:bg-red-950 dark:text-red-300',
}

export const difficultyLabel: Record<Quest['difficulty'], string> = {
  easy: 'Mudah',
  medium: 'Sedang',
  hard: 'Sulit',
}
