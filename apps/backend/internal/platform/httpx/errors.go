package httpx

// TODO (response error yang seragam):
//   // Bentuk error yang dikirim ke client, mis:
//   //   { "error": { "code": "not_found", "message": "..." } }
//   type ErrorResponse struct { ... }
//
//   // Error menulis ErrorResponse dengan status tertentu.
//   func Error(w http.ResponseWriter, status int, code, message string)
//
//   // Shortcut umum:
//   func BadRequest(w, msg)     // 400
//   func Unauthorized(w, msg)   // 401
//   func Forbidden(w, msg)      // 403
//   func NotFound(w, msg)       // 404
//   func Conflict(w, msg)       // 409
//   func Internal(w)            // 500 (jangan bocorkan detail internal)
//
//   // (Opsional) Map error domain -> status. Tiap module mendefinisikan
//   // error domainnya (mis. quest.ErrQuestNotFound); di sini/di handler
//   // dipetakan ke HTTP. Pertimbangkan errors.Is/As.
