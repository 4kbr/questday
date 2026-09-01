package quest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"questday/internal/platform/httpx"
	"questday/internal/platform/middleware"
	"questday/internal/platform/validator"
)

const fallbackTimezone = "Asia/Jakarta"

// handler tipis: decode/validate -> service -> tulis lewat httpx. "Hari ini"
// dihitung di sini dari timezone user (ADR-006), tak pernah di service.
type handler struct {
	svc *service
	v   *validator.Validator
}

func newHandler(svc *service, v *validator.Validator) *handler {
	return &handler{svc: svc, v: v}
}

// localDateFrom menghitung tanggal lokal user (jam 00:00) dari timezone di
// context. Timezone kosong/invalid -> fallback Asia/Jakarta. Nilai balik dipakai
// sebagai kolom DATE (wall-clock, tanpa makna zona).
func localDateFrom(r *http.Request) time.Time {
	tz, ok := middleware.TimezoneFrom(r.Context())
	if !ok || tz == "" {
		tz = fallbackTimezone
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
		if l, e := time.LoadLocation(fallbackTimezone); e == nil {
			loc = l
		}
	}
	now := time.Now().In(loc)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFrom(r.Context())
	if !ok {
		httpx.Unauthorized(w, "token tidak valid")
		return
	}

	var req CreateQuestRequest
	if err := httpx.DecodeAndValidate(w, r, &req, h.v); err != nil {
		httpx.BadRequest(w, err.Error())
		return
	}

	q, err := h.svc.CreateQuest(r.Context(), userID, req)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.Data(w, http.StatusCreated, toQuestResponse(q))
}

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFrom(r.Context())
	if !ok {
		httpx.Unauthorized(w, "token tidak valid")
		return
	}

	quests, err := h.svc.ListQuests(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	out := make([]QuestResponse, 0, len(quests))
	for _, q := range quests {
		out = append(out, toQuestResponse(q))
	}
	httpx.Data(w, http.StatusOK, out)
}

func (h *handler) today(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFrom(r.Context())
	if !ok {
		httpx.Unauthorized(w, "token tidak valid")
		return
	}

	res, err := h.svc.GetToday(r.Context(), userID, localDateFrom(r))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.Data(w, http.StatusOK, res)
}

func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFrom(r.Context())
	if !ok {
		httpx.Unauthorized(w, "token tidak valid")
		return
	}

	var req UpdateQuestRequest
	if err := httpx.DecodeAndValidate(w, r, &req, h.v); err != nil {
		httpx.BadRequest(w, err.Error())
		return
	}

	q, err := h.svc.UpdateQuest(r.Context(), userID, chi.URLParam(r, "questId"), req)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.Data(w, http.StatusOK, toQuestResponse(q))
}

func (h *handler) archive(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFrom(r.Context())
	if !ok {
		httpx.Unauthorized(w, "token tidak valid")
		return
	}

	if err := h.svc.ArchiveQuest(r.Context(), userID, chi.URLParam(r, "questId")); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (h *handler) complete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFrom(r.Context())
	if !ok {
		httpx.Unauthorized(w, "token tidak valid")
		return
	}

	log, err := h.svc.CompleteQuest(r.Context(), userID, chi.URLParam(r, "questId"), localDateFrom(r))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.Data(w, http.StatusOK, toQuestLogResponse(log))
}

func (h *handler) uncomplete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFrom(r.Context())
	if !ok {
		httpx.Unauthorized(w, "token tidak valid")
		return
	}

	if err := h.svc.UncompleteQuest(r.Context(), userID, chi.URLParam(r, "questId"), localDateFrom(r)); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.NoContent(w)
}
