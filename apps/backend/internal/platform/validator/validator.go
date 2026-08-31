// Package validator membungkus library validasi (mis. go-playground/validator)
// jadi satu instance yang dipakai bersama, plus penerjemahan error validasi
// menjadi pesan yang enak dibaca.
package validator

// TODO:
//   // New membuat validator siap pakai (register custom rule bila perlu).
//   func New() *Validator
//
//   // Struct memvalidasi v berdasarkan tag `validate:"..."` pada DTO.
//   // Kembalikan error yang sudah rapi (field -> pesan).
//   func (v *Validator) Struct(s any) error
