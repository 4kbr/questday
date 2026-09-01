/**
 * Satu-satunya sumber query key untuk quest. DILARANG menulis array literal
 * `['quests', ...]` di file lain — satu salah ketik dan invalidasi diam-diam
 * gagal. Phase 3 (F3.6) bergantung penuh pada key ini.
 */
export const questKeys = {
  all: () => ['quests'] as const,
  list: () => [...questKeys.all(), 'list'] as const,
  today: () => [...questKeys.all(), 'today'] as const,
  detail: (id: string) => [...questKeys.all(), 'detail', id] as const,
}
