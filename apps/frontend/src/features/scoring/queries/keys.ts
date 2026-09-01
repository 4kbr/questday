/**
 * Satu-satunya sumber query key untuk scoring. DILARANG menulis array literal
 * `['scoring', ...]` di file lain — satu salah ketik dan invalidasi diam-diam
 * gagal. F3.6 memanggil key-key ini dari `features/quest` (hanya konstanta key,
 * bukan hook/komponen).
 */
export const scoringKeys = {
  all: () => ['scoring'] as const,
  score: () => [...scoringKeys.all(), 'score'] as const,
  streak: () => [...scoringKeys.all(), 'streak'] as const,
  leaderboard: (limit: number) =>
    [...scoringKeys.all(), 'leaderboard', limit] as const,
}
