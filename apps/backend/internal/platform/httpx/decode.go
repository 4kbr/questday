package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"questday/internal/platform/validator"
)

const maxBodyBytes = 1 << 20 // 1 MiB

// Decode membaca JSON body ke dst: batasi ukuran, tolak field asing, dan
// terjemahkan error JSON rusak jadi pesan yang ramah. Tidak menulis response —
// caller yang memutuskan. Param w hanya untuk MaxBytesReader.
func Decode(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return translateDecodeError(err)
	}

	// Tolak JSON value kedua di body.
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("body hanya boleh berisi satu objek JSON")
	}

	return nil
}

// DecodeAndValidate menjalankan Decode lalu v.Struct(dst).
func DecodeAndValidate(w http.ResponseWriter, r *http.Request, dst any, v *validator.Validator) error {
	if err := Decode(w, r, dst); err != nil {
		return err
	}
	return v.Struct(dst)
}

func translateDecodeError(err error) error {
	var (
		syntaxErr   *json.SyntaxError
		typeErr     *json.UnmarshalTypeError
		maxBytesErr *http.MaxBytesError
	)

	switch {
	case errors.Is(err, io.EOF):
		return fmt.Errorf("body kosong")
	case errors.As(err, &syntaxErr):
		return fmt.Errorf("JSON tidak valid")
	case errors.Is(err, io.ErrUnexpectedEOF):
		return fmt.Errorf("JSON tidak valid")
	case errors.As(err, &typeErr):
		return fmt.Errorf("tipe field %s salah", typeErr.Field)
	case errors.As(err, &maxBytesErr):
		return fmt.Errorf("body terlalu besar")
	case strings.Contains(err.Error(), "unknown field"):
		return fmt.Errorf("%s", strings.TrimPrefix(err.Error(), "json: "))
	default:
		return fmt.Errorf("gagal membaca body: %w", err)
	}
}
