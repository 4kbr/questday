package httpx

import (
	"errors"
	"log"
	"net/http"
	"sync"
)

// ErrorResponse adalah bentuk error yang dikirim ke client:
//
//	{ "error": { "code": "not_found", "message": "..." } }
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// Error menulis ErrorResponse dengan status, code, dan message tertentu.
func Error(w http.ResponseWriter, status int, code, message string) {
	var body ErrorResponse
	body.Error.Code = code
	body.Error.Message = message
	JSON(w, status, body)
}

// BadRequest menulis 400 dengan code "bad_request".
func BadRequest(w http.ResponseWriter, msg string) {
	Error(w, http.StatusBadRequest, "bad_request", msg)
}

// Unauthorized menulis 401 dengan code "unauthorized".
func Unauthorized(w http.ResponseWriter, msg string) {
	Error(w, http.StatusUnauthorized, "unauthorized", msg)
}

// Forbidden menulis 403 dengan code "forbidden".
func Forbidden(w http.ResponseWriter, msg string) {
	Error(w, http.StatusForbidden, "forbidden", msg)
}

// NotFound menulis 404 dengan code "not_found".
func NotFound(w http.ResponseWriter, msg string) {
	Error(w, http.StatusNotFound, "not_found", msg)
}

// Conflict menulis 409 dengan code "conflict".
func Conflict(w http.ResponseWriter, msg string) {
	Error(w, http.StatusConflict, "conflict", msg)
}

// Internal menulis 500 generik. TIDAK pernah membocorkan detail error ke client.
func Internal(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "internal_error", "terjadi kesalahan internal")
}

// --- Registry error domain -> HTTP (ADR-023) --------------------------------

type errMapping struct {
	status int
	code   string
}

var (
	errMu  sync.RWMutex
	errMap = map[error]errMapping{}
)

// RegisterErrorMapping mendaftarkan sentinel error domain ke status + code HTTP.
// Dipanggil tiap module di New(). Aman dipanggil concurrent; diisi sekali saat
// startup.
func RegisterErrorMapping(err error, status int, code string) {
	if err == nil {
		return
	}
	errMu.Lock()
	defer errMu.Unlock()
	errMap[err] = errMapping{status: status, code: code}
}

// WriteError menulis response error untuk err: cari lewat errors.Is di registry,
// fallback ke Internal(w) (log, tanpa bocor) bila tak ada yang cocok.
func WriteError(w http.ResponseWriter, err error) {
	if err == nil {
		Internal(w)
		return
	}

	errMu.RLock()
	for sentinel, m := range errMap {
		if errors.Is(err, sentinel) {
			errMu.RUnlock()
			Error(w, m.status, m.code, err.Error())
			return
		}
	}
	errMu.RUnlock()

	log.Printf("httpx: unmapped error: %v", err)
	Internal(w)
}
