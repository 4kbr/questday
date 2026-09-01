import type { components } from './schema.gen'

type S = components['schemas']

export type User = S['UserResponse']
export type AuthResponse = S['AuthResponse']
export type RegisterRequest = S['RegisterRequest']
export type LoginRequest = S['LoginRequest']
export type Quest = S['QuestResponse']
export type CreateQuestRequest = S['CreateQuestRequest']
export type UpdateQuestRequest = S['UpdateQuestRequest']
export type TodayQuests = S['TodayQuestsResponse']
export type Score = S['ScoreResponse']
export type Streak = S['StreakResponse']
export type LeaderboardEntry = S['LeaderboardEntry']
export type ApiErrorBody = S['ErrorResponse']

// Optional aliases — schema keys confirmed present in schema.gen.ts
export type QuestLog = S['QuestLogResponse']
export type CompleteQuestRequest = S['CompleteQuestRequest']
export type UpdateProfileRequest = S['UpdateProfileRequest']
export type HealthResponse = S['HealthResponse']
