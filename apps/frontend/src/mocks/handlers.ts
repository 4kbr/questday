import { http, HttpResponse } from 'msw'
import type {
  ApiErrorBody,
  AuthResponse,
  CreateQuestRequest,
  HealthResponse,
  LeaderboardEntry,
  Quest,
  QuestLog,
  Score,
  Streak,
  TodayQuests,
  UpdateProfileRequest,
  UpdateQuestRequest,
  User,
} from '@/apis/types'

/**
 * Base URL untuk semua handler. Sama dengan `baseURL` axios di `apis/client.ts`
 * sehingga request yang keluar cocok dengan path lengkap yang di-intercept.
 * Contoh: `http://localhost:8080/api/v1`.
 */
const BASE = import.meta.env.VITE_API_BASE_URL

/**
 * Amplop error konsisten dengan kontrak: `{"error":{code,message}}`.
 * Dipakai ulang oleh handler fitur di phase berikutnya (F1.10 / F2.11 / F3.8).
 */
export function errorBody(code: string, message: string): ApiErrorBody {
  return { error: { code, message } }
}

/**
 * Bangun response error ber-amplop dengan status HTTP tertentu.
 * Error di-amplop `{error:{code,message}}`; sukses di-amplop `{data}` lewat
 * `dataResponse` (ADR-025).
 */
export function errorResponse(status: number, code: string, message: string) {
  return HttpResponse.json(errorBody(code, message), { status })
}

/**
 * Amplop sukses konsisten dengan ADR-025: `{"data": <payload>}`.
 * `client.ts` meng-unwrap sekali di interceptor sukses, jadi handler sukses
 * WAJIB membungkus payload di sini.
 */
export function dataResponse<T>(payload: T, status = 200) {
  return HttpResponse.json({ data: payload }, { status })
}

// --- Fake in-memory store supaya alur auth koheren tanpa backend --------------

/** User ter-seed; password-nya `'password123'`. */
const seededUser: User = {
  id: '00000000-0000-4000-8000-000000000001',
  email: 'demo@questday.test',
  display_name: 'Demo',
  timezone: 'Asia/Jakarta',
}
const seededPassword = 'password123'

/** Email yang sudah "terdaftar" lewat register — supaya register ulang → 409. */
const registeredEmails = new Set<string>([seededUser.email])

const MOCK_TOKEN = 'mock-jwt-token'

type LoginBody = { email?: string; password?: string }
type RegisterBody = {
  email?: string
  password?: string
  display_name?: string
  timezone?: string
}

// --- Fake in-memory store quest (F2.11) -------------------------------------
//
// Mock WAJIB stateful: complete/uncomplete harus benar-benar mengubah
// `GET /quests/today`, kalau tidak optimistic update (F2.6) tak pernah teruji.

/** Poin per difficulty — satu sumber, cocok dengan `Quest.Points()` backend. */
const POINTS_BY_DIFFICULTY = { easy: 5, medium: 10, hard: 20 } as const

const seedNow = new Date().toISOString()

/** 3 quest ter-seed, semua aktif, milik `seededUser`. Uuid tetap. */
const quests: Quest[] = [
  {
    id: '10000000-0000-4000-8000-000000000001',
    user_id: seededUser.id,
    title: 'Baca 10 halaman buku',
    note: 'Buku non-fiksi',
    category: 'Belajar',
    difficulty: 'easy',
    recurrence: 'daily',
    points: POINTS_BY_DIFFICULTY.easy,
    active: true,
    created_at: seedNow,
  },
  {
    id: '10000000-0000-4000-8000-000000000002',
    user_id: seededUser.id,
    title: 'Olahraga 30 menit',
    category: 'Kesehatan',
    difficulty: 'medium',
    recurrence: 'daily',
    points: POINTS_BY_DIFFICULTY.medium,
    active: true,
    created_at: seedNow,
  },
  {
    id: '10000000-0000-4000-8000-000000000003',
    user_id: seededUser.id,
    title: 'Tulis jurnal harian',
    category: 'Refleksi',
    difficulty: 'hard',
    recurrence: 'daily',
    points: POINTS_BY_DIFFICULTY.hard,
    active: true,
    created_at: seedNow,
  },
]

/** Set id quest yang selesai "hari ini" (MVP: hanya hari ini). */
const todayCompleted = new Set<string>()

/** Tanggal "hari ini" versi mock — boleh pakai jam asli, ini fake backend. */
function todayISO(): string {
  return new Date().toISOString().slice(0, 10)
}

/** Bangun `QuestLog` untuk quest yang baru diselesaikan. */
function buildLog(q: Quest): QuestLog {
  return {
    id: crypto.randomUUID(),
    quest_id: q.id,
    user_id: seededUser.id,
    date: todayISO(),
    status: 'completed',
    points_awarded: POINTS_BY_DIFFICULTY[q.difficulty],
    completed_at: new Date().toISOString(),
  }
}

/** Cari quest aktif berdasar id. */
function findActiveQuest(id: string): Quest | undefined {
  return quests.find((q) => q.id === id && q.active)
}

// --- Fake in-memory store scoring (F3.8) ----------------------------------
//
// Mock scoring WAJIB terhubung dengan state quest (`todayCompleted`). Kalau
// statis, F3.6 (invalidasi lintas-fitur) terlihat "berhasil" tanpa membuktikan
// apa pun. Complete quest di mock → poin & streak & leaderboard ikut naik.

/** Streak terpanjang yang pernah terlihat sesi ini. */
let longestStreakSeen = 0

/** Jumlah poin dari semua quest aktif yang selesai hari ini. */
function sumCompletedPoints(): number {
  return quests
    .filter((q) => q.active && todayCompleted.has(q.id))
    .reduce((sum, q) => sum + POINTS_BY_DIFFICULTY[q.difficulty], 0)
}

function computeScore(): Score {
  const xp = sumCompletedPoints()
  // MOCK-ONLY level formula — fake backend. FRONTEND MUST NOT replicate this.
  const level = Math.floor(xp / 100) + 1
  const points_to_next_level = 100 - (xp % 100)
  return { total_points: xp, xp, level, points_to_next_level }
}

function computeStreak(): Streak {
  const current = todayCompleted.size > 0 ? 1 : 0
  longestStreakSeen = Math.max(longestStreakSeen, current)
  return {
    current,
    longest: longestStreakSeen,
    last_active: current > 0 ? todayISO() : null,
  }
}

/** 2 user palsu tetap untuk membuat ranking menarik. Uuid tetap. */
const FAKE_LEADERBOARD: LeaderboardEntry[] = [
  {
    rank: 0,
    user_id: '00000000-0000-4000-8000-0000000000a1',
    display_name: 'Budi',
    points: 45,
  },
  {
    rank: 0,
    user_id: '00000000-0000-4000-8000-0000000000a2',
    display_name: 'Sari',
    points: 15,
  },
]

export const handlers = [
  http.get(`${BASE}/healthz`, () =>
    HttpResponse.json({ status: 'ok' } satisfies HealthResponse),
  ),

  http.post(`${BASE}/auth/login`, async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as LoginBody
    if (body.email === seededUser.email && body.password === seededPassword) {
      return dataResponse<AuthResponse>(
        { token: MOCK_TOKEN, user: seededUser },
        200,
      )
    }
    return errorResponse(401, 'invalid_credential', 'Email atau password salah')
  }),

  http.post(`${BASE}/auth/register`, async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as RegisterBody
    const email = body.email ?? ''
    if (registeredEmails.has(email)) {
      return errorResponse(409, 'email_taken', 'Email sudah terdaftar')
    }
    const user: User = {
      id: crypto.randomUUID(),
      email,
      display_name: body.display_name ?? '',
      timezone: body.timezone ?? 'Asia/Jakarta',
    }
    registeredEmails.add(email)
    return dataResponse<AuthResponse>({ token: MOCK_TOKEN, user }, 200)
  }),

  http.get(`${BASE}/me`, ({ request }) => {
    const auth = request.headers.get('Authorization') ?? ''
    if (auth.startsWith('Bearer ')) {
      return dataResponse<User>(seededUser, 200)
    }
    // Selaras backend asli (F4.9): middleware auth memakai code `unauthorized`
    // untuk 401 token hilang/invalid — bukan `invalid_credential` (itu khusus
    // login email/password salah).
    return errorResponse(401, 'unauthorized', 'Token tidak ada')
  }),

  // `PATCH /me` (F4.1) — update parsial profil. Mengembalikan `AuthResponse`
  // (token BARU + user) karena timezone ikut di JWT claims (ADR-013/ADR-022).
  http.patch(`${BASE}/me`, async ({ request }) => {
    const body = (await request
      .json()
      .catch(() => ({}))) as Partial<UpdateProfileRequest>
    if (body.display_name !== undefined) {
      seededUser.display_name = body.display_name
    }
    if (body.timezone !== undefined) seededUser.timezone = body.timezone
    return dataResponse<AuthResponse>(
      { token: MOCK_TOKEN, user: seededUser },
      200,
    )
  }),

  // --- Quest (F2.11) -------------------------------------------------------
  // `/quests/today` WAJIB terdaftar SEBELUM `/quests/:questId` supaya "today"
  // tidak tertangkap sebagai id.

  http.get(`${BASE}/quests`, () =>
    dataResponse<Quest[]>(quests.filter((q) => q.active)),
  ),

  http.get(`${BASE}/quests/today`, () =>
    dataResponse<TodayQuests>({
      date: todayISO(),
      items: quests
        .filter((q) => q.active)
        .map((q) => ({ quest: q, completed: todayCompleted.has(q.id) })),
    }),
  ),

  http.post(`${BASE}/quests`, async ({ request }) => {
    const body = (await request
      .json()
      .catch(() => ({}))) as Partial<CreateQuestRequest>
    if (!body.title || !body.category || !body.difficulty) {
      return errorResponse(
        400,
        'validation_failed',
        'title, category, dan difficulty wajib diisi',
      )
    }
    const quest: Quest = {
      id: crypto.randomUUID(),
      user_id: seededUser.id,
      title: body.title,
      note: body.note,
      category: body.category,
      difficulty: body.difficulty,
      recurrence: 'daily',
      points: POINTS_BY_DIFFICULTY[body.difficulty],
      active: true,
      created_at: new Date().toISOString(),
    }
    quests.push(quest)
    return dataResponse<Quest>(quest, 200)
  }),

  http.patch(`${BASE}/quests/:questId`, async ({ request, params }) => {
    const id = String(params.questId)
    const quest = findActiveQuest(id)
    if (!quest) {
      return errorResponse(404, 'quest_not_found', 'Quest tidak ditemukan')
    }
    const body = (await request
      .json()
      .catch(() => ({}))) as Partial<UpdateQuestRequest>
    if (body.title !== undefined) quest.title = body.title
    if (body.note !== undefined) quest.note = body.note
    if (body.category !== undefined) quest.category = body.category
    if (body.difficulty !== undefined) {
      quest.difficulty = body.difficulty
      quest.points = POINTS_BY_DIFFICULTY[body.difficulty]
    }
    if (body.active !== undefined) quest.active = body.active
    return dataResponse<Quest>(quest)
  }),

  http.delete(`${BASE}/quests/:questId`, ({ params }) => {
    const id = String(params.questId)
    const quest = quests.find((q) => q.id === id)
    if (!quest) {
      return errorResponse(404, 'quest_not_found', 'Quest tidak ditemukan')
    }
    quest.active = false
    todayCompleted.delete(id)
    return new HttpResponse(null, { status: 204 })
  }),

  http.post(`${BASE}/quests/:questId/complete`, ({ params }) => {
    const id = String(params.questId)
    const quest = findActiveQuest(id)
    if (!quest) {
      return errorResponse(404, 'quest_not_found', 'Quest tidak ditemukan')
    }
    if (todayCompleted.has(id)) {
      return errorResponse(
        409,
        'already_completed',
        'Quest sudah diselesaikan hari ini',
      )
    }
    todayCompleted.add(id)
    return dataResponse<QuestLog>(buildLog(quest), 200)
  }),

  http.post(`${BASE}/quests/:questId/uncomplete`, ({ params }) => {
    const id = String(params.questId)
    const quest = findActiveQuest(id)
    if (!quest) {
      return errorResponse(404, 'quest_not_found', 'Quest tidak ditemukan')
    }
    if (!todayCompleted.has(id)) {
      return errorResponse(
        409,
        'not_completed',
        'Quest belum diselesaikan hari ini',
      )
    }
    todayCompleted.delete(id)
    return new HttpResponse(null, { status: 204 })
  }),

  // --- Scoring (F3.8) ----------------------------------------------------
  // Ketiganya membaca `todayCompleted`, jadi complete/uncomplete quest di mock
  // menggeser score, streak, DAN leaderboard sekaligus (syarat F3.8).

  http.get(`${BASE}/me/score`, () => dataResponse<Score>(computeScore())),

  http.get(`${BASE}/me/streak`, () => dataResponse<Streak>(computeStreak())),

  http.get(`${BASE}/leaderboard`, ({ request }) => {
    const limitParam = new URL(request.url).searchParams.get('limit')
    const limit = limitParam ? Number(limitParam) : 20
    const rows: LeaderboardEntry[] = [
      {
        rank: 0,
        user_id: seededUser.id,
        display_name: seededUser.display_name,
        points: computeScore().total_points,
      },
      ...FAKE_LEADERBOARD.map((e) => ({ ...e })),
    ]
    rows.sort((a, b) => b.points - a.points)
    const ranked = rows.map((e, i) => ({ ...e, rank: i + 1 })).slice(0, limit)
    return dataResponse<LeaderboardEntry[]>(ranked)
  }),
]
