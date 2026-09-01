package quest

import (
	"context"
	"errors"
	"testing"
	"time"
)

// --- fakes -----------------------------------------------------------

func logKey(questID string, d time.Time) string {
	return questID + "|" + d.Format("2006-01-02")
}

// fakeRepo = Repository in-memory. Sengaja TIDAK memfilter user_id supaya
// cabang ErrNotOwner di service bisa diuji.
type fakeRepo struct {
	quests map[string]Quest
	logs   map[string]QuestLog
	seq    int
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{quests: map[string]Quest{}, logs: map[string]QuestLog{}}
}

func (f *fakeRepo) CreateQuest(_ context.Context, q Quest) error {
	if q.CreatedAt.IsZero() {
		q.CreatedAt = time.Now()
	}
	f.quests[q.ID] = q
	return nil
}

func (f *fakeRepo) GetQuest(_ context.Context, _ /*userID*/, id string) (Quest, error) {
	q, ok := f.quests[id]
	if !ok {
		return Quest{}, ErrQuestNotFound
	}
	return q, nil
}

func (f *fakeRepo) ListQuestsByUser(_ context.Context, userID string) ([]Quest, error) {
	var out []Quest
	for _, q := range f.quests {
		if q.UserID == userID {
			out = append(out, q)
		}
	}
	return out, nil
}

func (f *fakeRepo) UpdateQuest(_ context.Context, q Quest) error {
	if _, ok := f.quests[q.ID]; !ok {
		return ErrQuestNotFound
	}
	f.quests[q.ID] = q
	return nil
}

func (f *fakeRepo) ArchiveQuest(_ context.Context, _ /*userID*/, id string) error {
	q, ok := f.quests[id]
	if !ok {
		return ErrQuestNotFound
	}
	q.Active = false
	f.quests[id] = q
	return nil
}

func (f *fakeRepo) CreateLog(_ context.Context, l QuestLog) error {
	k := logKey(l.QuestID, l.Date)
	if _, ok := f.logs[k]; ok {
		return ErrAlreadyCompleted
	}
	l.CompletedAt = time.Now()
	f.logs[k] = l
	return nil
}

func (f *fakeRepo) GetLog(_ context.Context, _ /*userID*/, questID string, date time.Time) (QuestLog, error) {
	l, ok := f.logs[logKey(questID, date)]
	if !ok {
		return QuestLog{}, ErrNotCompleted
	}
	return l, nil
}

func (f *fakeRepo) DeleteLog(_ context.Context, _ /*userID*/, questID string, date time.Time) error {
	k := logKey(questID, date)
	if _, ok := f.logs[k]; !ok {
		return ErrNotCompleted
	}
	delete(f.logs, k)
	return nil
}

func (f *fakeRepo) ListLogsByUserAndDate(_ context.Context, userID string, date time.Time) ([]QuestLog, error) {
	var out []QuestLog
	for _, l := range f.logs {
		if l.UserID == userID && l.Date.Format("2006-01-02") == date.Format("2006-01-02") {
			out = append(out, l)
		}
	}
	return out, nil
}

func (f *fakeRepo) ListActiveDates(_ context.Context, userID string, from, to time.Time) ([]time.Time, error) {
	seen := map[string]time.Time{}
	for _, l := range f.logs {
		if l.UserID == userID && !l.Date.Before(from) && !l.Date.After(to) {
			seen[l.Date.Format("2006-01-02")] = l.Date
		}
	}
	var out []time.Time
	for _, d := range seen {
		out = append(out, d)
	}
	return out, nil
}

// fakeScorer mencatat panggilan.
type fakeScorer struct {
	completed   int
	uncompleted int
	lastPoints  int
	lastDate    time.Time
}

func (s *fakeScorer) OnQuestCompleted(_ context.Context, _, _ string, points int, date time.Time) error {
	s.completed++
	s.lastPoints = points
	s.lastDate = date
	return nil
}

func (s *fakeScorer) OnQuestUncompleted(_ context.Context, _, _ string, points int, date time.Time) error {
	s.uncompleted++
	s.lastPoints = points
	s.lastDate = date
	return nil
}

func newSvc() (*service, *fakeRepo, *fakeScorer) {
	repo := newFakeRepo()
	sc := &fakeScorer{}
	return newService(repo, sc), repo, sc
}

func seedQuest(f *fakeRepo, id, userID string, diff Difficulty) {
	f.quests[id] = Quest{
		ID: id, UserID: userID, Title: "T", Category: "c",
		Difficulty: diff, Recurrence: RecurrenceDaily, Active: true,
		CreatedAt: time.Now(),
	}
}

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// --- tests -----------------------------------------------------------

func TestCreateQuest_InvalidDifficulty(t *testing.T) {
	svc, _, _ := newSvc()
	_, err := svc.CreateQuest(context.Background(), "u1", CreateQuestRequest{
		Title: "Lari", Category: "olahraga", Difficulty: "ekstrem",
	})
	if err == nil {
		t.Fatal("difficulty ngawur harus ditolak")
	}
}

func TestCompleteQuest_HappyPath(t *testing.T) {
	svc, repo, sc := newSvc()
	seedQuest(repo, "q1", "u1", DifficultyMedium)

	d := date(2026, 9, 1)
	log, err := svc.CompleteQuest(context.Background(), "u1", "q1", d)
	if err != nil {
		t.Fatalf("CompleteQuest: %v", err)
	}
	if log.PointsAwarded != 10 {
		t.Fatalf("PointsAwarded = %d, mau 10", log.PointsAwarded)
	}
	if sc.completed != 1 {
		t.Fatalf("scorer.OnQuestCompleted dipanggil %d kali, mau 1", sc.completed)
	}
	if sc.lastPoints != 10 {
		t.Fatalf("scorer dapat %d poin, mau 10 (== Quest.Points())", sc.lastPoints)
	}
	if !sc.lastDate.Equal(d) {
		t.Fatalf("scorer dapat tanggal %v, mau %v", sc.lastDate, d)
	}
}

func TestCompleteQuest_NotOwner_ScorerNotCalled(t *testing.T) {
	svc, repo, sc := newSvc()
	seedQuest(repo, "q1", "owner", DifficultyEasy)

	_, err := svc.CompleteQuest(context.Background(), "penyusup", "q1", date(2026, 9, 1))
	if !errors.Is(err, ErrNotOwner) {
		t.Fatalf("err = %v, mau ErrNotOwner", err)
	}
	if sc.completed != 0 {
		t.Fatalf("scorer tidak boleh dipanggil untuk quest orang lain (dipanggil %d kali)", sc.completed)
	}
}

func TestCompleteQuest_TwiceSameDate(t *testing.T) {
	svc, repo, _ := newSvc()
	seedQuest(repo, "q1", "u1", DifficultyHard)
	d := date(2026, 9, 1)

	if _, err := svc.CompleteQuest(context.Background(), "u1", "q1", d); err != nil {
		t.Fatalf("complete pertama: %v", err)
	}
	_, err := svc.CompleteQuest(context.Background(), "u1", "q1", d)
	if !errors.Is(err, ErrAlreadyCompleted) {
		t.Fatalf("err = %v, mau ErrAlreadyCompleted", err)
	}
}

func TestUncompleteQuest_NoLog(t *testing.T) {
	svc, repo, sc := newSvc()
	seedQuest(repo, "q1", "u1", DifficultyEasy)

	err := svc.UncompleteQuest(context.Background(), "u1", "q1", date(2026, 9, 1))
	if !errors.Is(err, ErrNotCompleted) {
		t.Fatalf("err = %v, mau ErrNotCompleted", err)
	}
	if sc.uncompleted != 0 {
		t.Fatalf("scorer tidak boleh dipanggil (dipanggil %d kali)", sc.uncompleted)
	}
}

func TestUncompleteQuest_RollsBackPoints(t *testing.T) {
	svc, repo, sc := newSvc()
	seedQuest(repo, "q1", "u1", DifficultyMedium)
	d := date(2026, 9, 1)
	svc.CompleteQuest(context.Background(), "u1", "q1", d)

	if err := svc.UncompleteQuest(context.Background(), "u1", "q1", d); err != nil {
		t.Fatalf("UncompleteQuest: %v", err)
	}
	if sc.uncompleted != 1 || sc.lastPoints != 10 {
		t.Fatalf("rollback: uncompleted=%d lastPoints=%d, mau 1 & 10", sc.uncompleted, sc.lastPoints)
	}
	// log hari itu harus hilang
	if _, err := svc.CompleteQuest(context.Background(), "u1", "q1", d); err != nil {
		t.Fatalf("setelah uncomplete harus bisa complete lagi: %v", err)
	}
}

// Batas hari: instant yang sama, dua timezone -> localDate beda -> log jatuh di
// tanggal yang benar sesuai argumen (service tak pernah pakai time.Now()).
func TestCompleteQuest_DayBoundary(t *testing.T) {
	svc, repo, _ := newSvc()
	seedQuest(repo, "qJKT", "u1", DifficultyEasy)
	seedQuest(repo, "qUTC", "u1", DifficultyEasy)

	instant := time.Date(2026, 8, 31, 17, 0, 0, 0, time.UTC) // 2026-09-01 00:00 WIB
	jkt, _ := time.LoadLocation("Asia/Jakarta")

	local := func(loc *time.Location) time.Time {
		w := instant.In(loc)
		return time.Date(w.Year(), w.Month(), w.Day(), 0, 0, 0, 0, time.UTC)
	}

	logJKT, err := svc.CompleteQuest(context.Background(), "u1", "qJKT", local(jkt))
	if err != nil {
		t.Fatalf("complete JKT: %v", err)
	}
	logUTC, err := svc.CompleteQuest(context.Background(), "u1", "qUTC", local(time.UTC))
	if err != nil {
		t.Fatalf("complete UTC: %v", err)
	}

	if got := logJKT.Date.Format("2006-01-02"); got != "2026-09-01" {
		t.Fatalf("log WIB jatuh di %s, mau 2026-09-01", got)
	}
	if got := logUTC.Date.Format("2006-01-02"); got != "2026-08-31" {
		t.Fatalf("log UTC jatuh di %s, mau 2026-08-31", got)
	}
}

// erroringScorer selalu gagal — untuk menguji kompensasi ADR-016.
type erroringScorer struct{ calls int }

func (s *erroringScorer) OnQuestCompleted(context.Context, string, string, int, time.Time) error {
	s.calls++
	return errors.New("scorer meledak")
}
func (s *erroringScorer) OnQuestUncompleted(context.Context, string, string, int, time.Time) error {
	return errors.New("scorer meledak")
}

// ADR-016: kalau OnQuestCompleted gagal, CompleteQuest harus menghapus log yang
// baru dibuat supaya tak ada log yatim.
func TestCompleteQuest_ScorerFails_NoOrphanLog(t *testing.T) {
	repo := newFakeRepo()
	sc := &erroringScorer{}
	svc := newService(repo, sc)
	seedQuest(repo, "q1", "u1", DifficultyMedium)
	d := date(2026, 9, 1)

	_, err := svc.CompleteQuest(context.Background(), "u1", "q1", d)
	if err == nil {
		t.Fatal("CompleteQuest harus mengembalikan error saat scorer gagal")
	}
	if sc.calls != 1 {
		t.Fatalf("scorer dipanggil %d kali, mau 1", sc.calls)
	}
	if _, ok := repo.logs[logKey("q1", d)]; ok {
		t.Fatal("log yatim tertinggal — kompensasi ADR-016 gagal")
	}
	// harus bisa complete lagi setelah kompensasi
	svc2 := newService(repo, &fakeScorer{})
	if _, err := svc2.CompleteQuest(context.Background(), "u1", "q1", d); err != nil {
		t.Fatalf("complete ulang setelah kompensasi: %v", err)
	}
}

func TestGetToday_MarksCompleted(t *testing.T) {
	svc, repo, _ := newSvc()
	seedQuest(repo, "q1", "u1", DifficultyEasy)
	seedQuest(repo, "q2", "u1", DifficultyHard)
	// q3 diarsipkan -> tak muncul
	seedQuest(repo, "q3", "u1", DifficultyEasy)
	q3 := repo.quests["q3"]
	q3.Active = false
	repo.quests["q3"] = q3

	d := date(2026, 9, 1)
	svc.CompleteQuest(context.Background(), "u1", "q1", d)

	res, err := svc.GetToday(context.Background(), "u1", d)
	if err != nil {
		t.Fatalf("GetToday: %v", err)
	}
	if res.Date != "2026-09-01" {
		t.Fatalf("Date = %q", res.Date)
	}
	if len(res.Items) != 2 {
		t.Fatalf("items = %d, mau 2 (q3 diarsipkan)", len(res.Items))
	}
	done := map[string]bool{}
	for _, it := range res.Items {
		done[it.Quest.ID] = it.Completed
	}
	if !done["q1"] || done["q2"] {
		t.Fatalf("status salah: q1=%v q2=%v (mau true, false)", done["q1"], done["q2"])
	}
}
