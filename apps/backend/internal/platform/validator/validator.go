// Package validator membungkus library validasi (go-playground/validator)
// jadi satu instance yang dipakai bersama, plus penerjemahan error validasi
// menjadi pesan yang enak dibaca.
package validator

import (
	"errors"
	"fmt"
	"strings"
	"time"

	govalidator "github.com/go-playground/validator/v10"
)

// Validator membungkus *govalidator.Validate.
type Validator struct {
	v *govalidator.Validate
}

// New membuat validator siap pakai, lengkap dengan rule kustom "timezone"
// (nama zona IANA valid via time.LoadLocation).
func New() *Validator {
	v := govalidator.New(govalidator.WithRequiredStructEnabled())

	_ = v.RegisterValidation("timezone", func(fl govalidator.FieldLevel) bool {
		s := fl.Field().String()
		if s == "" {
			return false
		}
		_, err := time.LoadLocation(s)
		return err == nil
	})

	return &Validator{v: v}
}

// Struct memvalidasi s berdasarkan tag `validate:"..."`. Error dikembalikan
// sebagai pesan rapi "field: pesan; field: pesan".
func (val *Validator) Struct(s any) error {
	err := val.v.Struct(s)
	if err == nil {
		return nil
	}

	var verrs govalidator.ValidationErrors
	if !errors.As(err, &verrs) {
		return err
	}

	msgs := make([]string, 0, len(verrs))
	for _, fe := range verrs {
		msgs = append(msgs, fieldMessage(fe))
	}
	return errors.New(strings.Join(msgs, "; "))
}

// fieldMessage menerjemahkan satu FieldError ke pesan Bahasa Indonesia.
func fieldMessage(fe govalidator.FieldError) string {
	field := strings.ToLower(fe.Field())

	switch fe.Tag() {
	case "required":
		return field + ": wajib diisi"
	case "email":
		return field + ": format email tidak valid"
	case "min":
		return fmt.Sprintf("%s: minimal %s karakter", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s: maksimal %s karakter", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s: harus salah satu dari [%s]", field, fe.Param())
	case "timezone":
		return field + ": zona waktu tidak valid"
	default:
		return fmt.Sprintf("%s: tidak valid (%s)", field, fe.Tag())
	}
}
