package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"questday/internal/platform/auth"
)

// defaultTimezone dipakai saat register tanpa timezone (ADR-006).
const defaultTimezone = "Asia/Jakarta"

// service memuat use case auth & profil. Tak tahu HTTP, tak menyentuh SQL.
type service struct {
	repo   Repository
	issuer auth.Issuer
}

func newService(repo Repository, issuer auth.Issuer) *service {
	return &service{repo: repo, issuer: issuer}
}

// Register membuat user baru lalu menerbitkan token.
func (s *service) Register(ctx context.Context, in RegisterRequest) (AuthResponse, error) {
	tz := in.Timezone
	if tz == "" {
		tz = defaultTimezone
	}

	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("user: register: %w", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return AuthResponse{}, fmt.Errorf("user: register: generate id: %w", err)
	}

	u := User{
		ID:           id.String(),
		Email:        in.Email,
		PasswordHash: hash,
		DisplayName:  in.DisplayName,
		Timezone:     tz,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return AuthResponse{}, err
	}

	return s.issue(u)
}

// Login mencocokkan kredensial lalu menerbitkan token. Email tak ada maupun
// password salah -> ErrInvalidCredential yang sama (jangan bocorkan email mana
// yang terdaftar).
func (s *service) Login(ctx context.Context, in LoginRequest) (AuthResponse, error) {
	u, err := s.repo.GetByEmail(ctx, in.Email)
	if err != nil {
		return AuthResponse{}, ErrInvalidCredential
	}
	if err := auth.ComparePassword(u.PasswordHash, in.Password); err != nil {
		return AuthResponse{}, ErrInvalidCredential
	}
	return s.issue(u)
}

// Profile mengembalikan profil publik user.
func (s *service) Profile(ctx context.Context, userID string) (UserResponse, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return UserResponse{}, err
	}
	return toUserResponse(u), nil
}

// UpdateProfile menyimpan perubahan display_name/timezone lalu menerbitkan token
// BARU (ADR-022): timezone ikut di JWT claims (ADR-013), jadi tanpa token baru
// backend terus memakai timezone lama sampai token lama kedaluwarsa.
func (s *service) UpdateProfile(ctx context.Context, userID string, in UpdateProfileRequest) (AuthResponse, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return AuthResponse{}, err
	}

	if in.DisplayName != nil {
		u.DisplayName = *in.DisplayName
	}
	if in.Timezone != nil {
		u.Timezone = *in.Timezone
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return AuthResponse{}, err
	}
	return s.issue(u)
}

// NamesByIDs mengembalikan map userID -> display_name. Memenuhi port
// scoring.UserDirectory (ADR-014) — dipanggil scoring lewat interface, bukan
// dengan meng-import module user.
func (s *service) NamesByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return s.repo.ListNamesByIDs(ctx, ids)
}

// issue merangkai token + profil jadi AuthResponse.
func (s *service) issue(u User) (AuthResponse, error) {
	token, err := s.issuer.Issue(u.ID, u.Timezone)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("user: issue token: %w", err)
	}
	return AuthResponse{Token: token, User: toUserResponse(u)}, nil
}
