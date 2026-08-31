package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims adalah isi token QuestDay: userID + timezone user (ADR-013) plus
// klaim standar JWT.
type Claims struct {
	UserID   string `json:"uid"`
	Timezone string `json:"tz"`
	jwt.RegisteredClaims
}

// Issuer membuat token untuk userID + timezone dengan masa berlaku ttl.
type Issuer interface {
	Issue(userID, timezone string) (string, error)
}

// Verifier memvalidasi token dan mengembalikan Claims.
type Verifier interface {
	Verify(token string) (Claims, error)
}

// JWT mengimplementasi Issuer & Verifier memakai secret + ttl dari config.
type JWT struct {
	secret []byte
	ttl    time.Duration
}

var (
	_ Issuer   = (*JWT)(nil)
	_ Verifier = (*JWT)(nil)
)

// NewJWT membuat JWT dengan secret HS256 dan masa berlaku ttl.
func NewJWT(secret string, ttl time.Duration) *JWT {
	return &JWT{secret: []byte(secret), ttl: ttl}
}

// Issue menandatangani token HS256 berisi userID + timezone.
func (j *JWT) Issue(userID, timezone string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Timezone: timezone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("auth: sign token: %w", err)
	}
	return signed, nil
}

// Verify memvalidasi token (signature, kadaluwarsa, algoritma) dan mengembalikan
// Claims. Hanya HS256 yang diterima — alg:none / algoritma lain ditolak.
func (j *JWT) Verify(token string) (Claims, error) {
	keyFunc := func(_ *jwt.Token) (any, error) {
		return j.secret, nil
	}

	parsed, err := jwt.ParseWithClaims(token, &Claims{}, keyFunc, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return Claims{}, fmt.Errorf("auth: verify token: %w", err)
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return Claims{}, fmt.Errorf("auth: verify token: klaim tidak valid")
	}
	return *claims, nil
}
