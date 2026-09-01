package scoring

import (
	"net/http"
	"strconv"

	"questday/internal/platform/httpx"
	"questday/internal/platform/middleware"
)

const (
	defaultLeaderboardLimit = 20
	maxLeaderboardLimit     = 100
)

// handler tipis: ambil userID/param -> service -> tulis lewat httpx.
type handler struct {
	svc *service
}

func newHandler(svc *service) *handler {
	return &handler{svc: svc}
}

func (h *handler) score(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFrom(r.Context())
	if !ok {
		httpx.Unauthorized(w, "token tidak valid")
		return
	}
	res, err := h.svc.GetScore(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.Data(w, http.StatusOK, res)
}

func (h *handler) streak(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFrom(r.Context())
	if !ok {
		httpx.Unauthorized(w, "token tidak valid")
		return
	}
	res, err := h.svc.GetStreak(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.Data(w, http.StatusOK, res)
}

func (h *handler) leaderboard(w http.ResponseWriter, r *http.Request) {
	if _, ok := middleware.UserIDFrom(r.Context()); !ok {
		httpx.Unauthorized(w, "token tidak valid")
		return
	}

	res, err := h.svc.Leaderboard(r.Context(), leaderboardLimit(r))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.Data(w, http.StatusOK, res)
}

// leaderboardLimit membaca ?limit=, default 20, clamp ke [1, 100].
func leaderboardLimit(r *http.Request) int {
	raw := r.URL.Query().Get("limit")
	if raw == "" {
		return defaultLeaderboardLimit
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return defaultLeaderboardLimit
	}
	if n > maxLeaderboardLimit {
		return maxLeaderboardLimit
	}
	return n
}
