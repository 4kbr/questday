package quest

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ScoreAwarder = port KELUAR milik module quest (ADR-005). Diimplementasi
// module scoring, disuntik oleh server. quest TIDAK meng-import scoring.
type ScoreAwarder interface {
	OnQuestCompleted(ctx context.Context, userID, questID string, points int, date time.Time) error
	OnQuestUncompleted(ctx context.Context, userID, questID string, points int, date time.Time) error
}

// service = use case module quest. Tak tahu HTTP, tak menyentuh SQL. Tak pernah
// memanggil time.Now() — "hari ini" datang sebagai argumen dari handler (ADR-006).
type service struct {
	repo   Repository
	scorer ScoreAwarder
}

func newService(repo Repository, scorer ScoreAwarder) *service {
	return &service{repo: repo, scorer: scorer}
}

// CreateQuest membuat definisi quest baru milik userID.
func (s *service) CreateQuest(ctx context.Context, userID string, in CreateQuestRequest) (Quest, error) {
	diff := Difficulty(in.Difficulty)
	if !diff.Valid() {
		return Quest{}, fmt.Errorf("quest: difficulty tidak valid: %q", in.Difficulty)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return Quest{}, fmt.Errorf("quest: generate id: %w", err)
	}

	q := Quest{
		ID:         id.String(),
		UserID:     userID,
		Title:      in.Title,
		Note:       in.Note,
		Category:   Category(in.Category),
		Difficulty: diff,
		Recurrence: RecurrenceDaily,
		Active:     true,
	}
	if err := s.repo.CreateQuest(ctx, q); err != nil {
		return Quest{}, err
	}
	// created_at diisi DB; ambil ulang supaya response konsisten.
	return s.repo.GetQuest(ctx, userID, q.ID)
}

// ListQuests mengembalikan semua quest milik user (aktif & arsip).
func (s *service) ListQuests(ctx context.Context, userID string) ([]Quest, error) {
	return s.repo.ListQuestsByUser(ctx, userID)
}

// UpdateQuest menerapkan patch parsial ke quest milik user.
func (s *service) UpdateQuest(ctx context.Context, userID, id string, in UpdateQuestRequest) (Quest, error) {
	q, err := s.repo.GetQuest(ctx, userID, id)
	if err != nil {
		return Quest{}, err
	}

	if in.Title != nil {
		q.Title = *in.Title
	}
	if in.Note != nil {
		q.Note = *in.Note
	}
	if in.Category != nil {
		q.Category = Category(*in.Category)
	}
	if in.Difficulty != nil {
		diff := Difficulty(*in.Difficulty)
		if !diff.Valid() {
			return Quest{}, fmt.Errorf("quest: difficulty tidak valid: %q", *in.Difficulty)
		}
		q.Difficulty = diff
	}
	if in.Active != nil {
		q.Active = *in.Active
	}

	if err := s.repo.UpdateQuest(ctx, q); err != nil {
		return Quest{}, err
	}
	return s.repo.GetQuest(ctx, userID, id)
}

// ArchiveQuest men-soft-delete quest (active=false). quest_logs lama tetap ada
// untuk perhitungan streak.
func (s *service) ArchiveQuest(ctx context.Context, userID, id string) error {
	return s.repo.ArchiveQuest(ctx, userID, id)
}

// GetToday menggabungkan quest aktif user dengan log tanggal localDate.
func (s *service) GetToday(ctx context.Context, userID string, localDate time.Time) (TodayQuestsResponse, error) {
	quests, err := s.repo.ListQuestsByUser(ctx, userID)
	if err != nil {
		return TodayQuestsResponse{}, err
	}
	logs, err := s.repo.ListLogsByUserAndDate(ctx, userID, localDate)
	if err != nil {
		return TodayQuestsResponse{}, err
	}

	done := make(map[string]bool, len(logs))
	for _, l := range logs {
		done[l.QuestID] = true
	}

	items := make([]TodayQuestItem, 0, len(quests))
	for _, q := range quests {
		if !q.Active {
			continue
		}
		items = append(items, TodayQuestItem{Quest: toQuestResponse(q), Completed: done[q.ID]})
	}

	return TodayQuestsResponse{Date: localDate.Format(dateLayout), Items: items}, nil
}

// CompleteQuest menandai quest selesai untuk localDate lalu memicu penambahan
// poin lewat scorer. Log + poin belum atomik (keputusan di T3.10 / ADR-016).
func (s *service) CompleteQuest(ctx context.Context, userID, questID string, localDate time.Time) (QuestLog, error) {
	q, err := s.repo.GetQuest(ctx, userID, questID)
	if err != nil {
		return QuestLog{}, err
	}
	if q.UserID != userID {
		return QuestLog{}, ErrNotOwner
	}
	if !q.Active {
		return QuestLog{}, ErrQuestNotFound
	}

	id, err := uuid.NewV7()
	if err != nil {
		return QuestLog{}, fmt.Errorf("quest: generate log id: %w", err)
	}

	points := q.Points()
	log := QuestLog{
		ID:            id.String(),
		QuestID:       questID,
		UserID:        userID,
		Date:          localDate,
		Status:        LogStatusCompleted,
		PointsAwarded: points,
	}
	if err := s.repo.CreateLog(ctx, log); err != nil {
		return QuestLog{}, err
	}

	// ADR-016: dua write lintas-module tak dalam satu transaksi. Kalau scorer
	// gagal, kompensasi manual — hapus log yang baru dibuat supaya tak ada log
	// yatim (poin tak masuk tapi quest tampak selesai).
	if err := s.scorer.OnQuestCompleted(ctx, userID, questID, points, localDate); err != nil {
		if delErr := s.repo.DeleteLog(ctx, userID, questID, localDate); delErr != nil {
			return QuestLog{}, fmt.Errorf("quest: award score: %w (kompensasi gagal: %v)", err, delErr)
		}
		return QuestLog{}, fmt.Errorf("quest: award score: %w", err)
	}

	return s.repo.GetLog(ctx, userID, questID, localDate)
}

// UncompleteQuest membatalkan penyelesaian localDate lalu me-rollback poin.
func (s *service) UncompleteQuest(ctx context.Context, userID, questID string, localDate time.Time) error {
	q, err := s.repo.GetQuest(ctx, userID, questID)
	if err != nil {
		return err
	}
	if q.UserID != userID {
		return ErrNotOwner
	}

	log, err := s.repo.GetLog(ctx, userID, questID, localDate)
	if err != nil {
		return err // ErrNotCompleted bila tak ada
	}

	if err := s.repo.DeleteLog(ctx, userID, questID, localDate); err != nil {
		return err
	}

	if err := s.scorer.OnQuestUncompleted(ctx, userID, questID, log.PointsAwarded, localDate); err != nil {
		return fmt.Errorf("quest: rollback score: %w", err)
	}
	return nil
}
