// Package auth menyediakan primitif keamanan lintas-module: token JWT dan
// hashing password. Module user memakainya lewat interface (port), bukan
// bergantung ke detail implementasi.
package auth

// TODO (token):
//   type Claims struct { UserID string; ... }   // isi standar + custom
//
//   // Issuer membuat token untuk userID dengan masa berlaku ttl.
//   type Issuer interface { Issue(userID string) (token string, err error) }
//
//   // Verifier memvalidasi token dan mengembalikan Claims.
//   type Verifier interface { Verify(token string) (Claims, error) }
//
//   // JWT mengimplementasi Issuer & Verifier memakai secret + ttl dari config.
//   type JWT struct { secret []byte; ttl time.Duration }
//   func NewJWT(secret string, ttl time.Duration) *JWT
