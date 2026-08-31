// Package httpx berisi helper HTTP lintas-module: menulis response, decode
// request, dan memetakan error ke status code. Tujuannya supaya bentuk
// response (sukses & error) SERAGAM di seluruh API.
package httpx

// TODO (response sukses):
//   // JSON menulis status + body sebagai JSON dengan header yang benar.
//   func JSON(w http.ResponseWriter, status int, v any)
//
//   // Kalau mau pakai envelope konsisten, mis. { "data": ... }:
//   func Data(w http.ResponseWriter, status int, data any)
//
//   // NoContent untuk 204.
//   func NoContent(w http.ResponseWriter)
