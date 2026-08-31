// Package httpx berisi helper HTTP lintas-module: menulis response, decode
// request, dan memetakan error ke status code. Tujuannya supaya bentuk
// response (sukses & error) SERAGAM di seluruh API.
package httpx

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON menulis status + body sebagai JSON dengan header yang benar.
// Header di-set SEBELUM WriteHeader.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("httpx: encode response: %v", err)
	}
}

// dataEnvelope membungkus payload sukses dalam { "data": ... }.
type dataEnvelope struct {
	Data any `json:"data"`
}

// Data menulis payload sukses dengan envelope { "data": <data> }.
func Data(w http.ResponseWriter, status int, data any) {
	JSON(w, status, dataEnvelope{Data: data})
}

// NoContent menulis 204 tanpa body.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
