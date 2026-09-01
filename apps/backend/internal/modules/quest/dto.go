package quest

import "time"

// dto.go = bentuk data HTTP (request/response) module quest. Terpisah dari
// entitas domain. Acuan bentuk: contracts/openapi.yaml.

// CreateQuestRequest = body POST /quests.
type CreateQuestRequest struct {
	Title      string `json:"title"      validate:"required"`
	Note       string `json:"note"`
	Category   string `json:"category"   validate:"required"`
	Difficulty string `json:"difficulty" validate:"required,oneof=easy medium hard"`
}

// UpdateQuestRequest = body PATCH /quests/{questId}. Semua pointer supaya
// "tak dikirim" beda dari "dikirim kosong".
type UpdateQuestRequest struct {
	Title      *string `json:"title"      validate:"omitempty,min=1"`
	Note       *string `json:"note"`
	Category   *string `json:"category"   validate:"omitempty,min=1"`
	Difficulty *string `json:"difficulty" validate:"omitempty,oneof=easy medium hard"`
	Active     *bool   `json:"active"`
}

// CompleteQuestRequest = body POST /quests/{questId}/complete. Kosong untuk MVP;
// poin sepenuhnya dari Quest.Points().
type CompleteQuestRequest struct{}

// QuestResponse = bentuk Quest di API. `points` = hasil Quest.Points() (ADR-026).
type QuestResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Title      string `json:"title"`
	Note       string `json:"note"`
	Category   string `json:"category"`
	Difficulty string `json:"difficulty"`
	Recurrence string `json:"recurrence"`
	Points     int    `json:"points"`
	Active     bool   `json:"active"`
	CreatedAt  string `json:"created_at"`
}

// QuestLogResponse = bentuk QuestLog di API.
type QuestLogResponse struct {
	ID            string `json:"id"`
	QuestID       string `json:"quest_id"`
	UserID        string `json:"user_id"`
	Date          string `json:"date"` // "2026-08-31", tanggal lokal user
	Status        string `json:"status"`
	PointsAwarded int    `json:"points_awarded"`
	CompletedAt   string `json:"completed_at"`
}

// TodayQuestItem = satu baris di TodayQuestsResponse.
type TodayQuestItem struct {
	Quest     QuestResponse `json:"quest"`
	Completed bool          `json:"completed"`
}

// TodayQuestsResponse = GET /quests/today: quest aktif user + status hari ini.
type TodayQuestsResponse struct {
	Date  string           `json:"date"` // tanggal lokal user
	Items []TodayQuestItem `json:"items"`
}

const dateLayout = "2006-01-02"

func toQuestResponse(q Quest) QuestResponse {
	return QuestResponse{
		ID:         q.ID,
		UserID:     q.UserID,
		Title:      q.Title,
		Note:       q.Note,
		Category:   string(q.Category),
		Difficulty: string(q.Difficulty),
		Recurrence: string(q.Recurrence),
		Points:     q.Points(),
		Active:     q.Active,
		CreatedAt:  q.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func toQuestLogResponse(l QuestLog) QuestLogResponse {
	return QuestLogResponse{
		ID:            l.ID,
		QuestID:       l.QuestID,
		UserID:        l.UserID,
		Date:          l.Date.Format(dateLayout),
		Status:        string(l.Status),
		PointsAwarded: l.PointsAwarded,
		CompletedAt:   l.CompletedAt.UTC().Format(time.RFC3339),
	}
}
