package user

// dto.go = bentuk data HTTP (request/response) module user. Terpisah dari
// entitas domain. Acuan bentuk: contracts/openapi.yaml.

// RegisterRequest adalah body POST /auth/register.
//
// Timezone opsional; default "Asia/Jakarta" diisi di service (T1.2). Tag
// omitempty,timezone supaya string ngawur tetap ditolak tapi kosong boleh.
type RegisterRequest struct {
	Email       string `json:"email"        validate:"required,email"`
	Password    string `json:"password"     validate:"required,min=8"`
	DisplayName string `json:"display_name" validate:"required"`
	Timezone    string `json:"timezone"     validate:"omitempty,timezone"`
}

// LoginRequest adalah body POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateProfileRequest adalah body PATCH /me. Pointer supaya "tak dikirim" beda
// dari "dikirim kosong". Email tidak dapat diubah di MVP.
type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name" validate:"omitempty,min=1"`
	Timezone    *string `json:"timezone"     validate:"omitempty,timezone"`
}

// AuthResponse dikembalikan register, login, dan PATCH /me (ADR-022): token baru
// plus profil user.
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse adalah profil publik user. Tidak pernah memuat PasswordHash.
type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Timezone    string `json:"timezone"`
}

// toUserResponse adalah satu-satunya jalan User -> HTTP. Sengaja tak menyalin
// PasswordHash.
func toUserResponse(u User) UserResponse {
	return UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		Timezone:    u.Timezone,
	}
}
