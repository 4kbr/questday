package user

import (
	"net/http"

	"questday/internal/platform/httpx"
	"questday/internal/platform/middleware"
	"questday/internal/platform/validator"
)

// handler tipis: decode/validate -> service -> tulis lewat httpx. Nol logika
// bisnis. Error domain dipetakan lewat httpx.WriteError (registry ADR-023).
type handler struct {
	svc *service
	v   *validator.Validator
}

func newHandler(svc *service, v *validator.Validator) *handler {
	return &handler{svc: svc, v: v}
}

// register menangani POST /auth/register.
func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := httpx.DecodeAndValidate(w, r, &req, h.v); err != nil {
		httpx.ValidationFailed(w, err.Error())
		return
	}

	res, err := h.svc.Register(r.Context(), req)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	// Kontrak: register mengembalikan 200 (bukan 201).
	httpx.Data(w, http.StatusOK, res)
}

// login menangani POST /auth/login.
func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := httpx.DecodeAndValidate(w, r, &req, h.v); err != nil {
		httpx.ValidationFailed(w, err.Error())
		return
	}

	res, err := h.svc.Login(r.Context(), req)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.Data(w, http.StatusOK, res)
}

// me menangani GET /me. userID diambil dari context (diisi Authenticator).
func (h *handler) me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFrom(r.Context())
	if !ok {
		httpx.Unauthorized(w, "token tidak valid")
		return
	}

	res, err := h.svc.Profile(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.Data(w, http.StatusOK, res)
}

// updateMe menangani PATCH /me. Mengembalikan AuthResponse (token baru, ADR-022).
func (h *handler) updateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFrom(r.Context())
	if !ok {
		httpx.Unauthorized(w, "token tidak valid")
		return
	}

	var req UpdateProfileRequest
	if err := httpx.DecodeAndValidate(w, r, &req, h.v); err != nil {
		httpx.ValidationFailed(w, err.Error())
		return
	}

	res, err := h.svc.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.Data(w, http.StatusOK, res)
}
