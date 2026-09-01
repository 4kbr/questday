package user

import (
	"context"
	"errors"
	"strings"
	"testing"

	"questday/internal/platform/auth"
)

// fakeRepo adalah Repository in-memory untuk test service.
type fakeRepo struct {
	byID    map[string]User
	byEmail map[string]User
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{byID: map[string]User{}, byEmail: map[string]User{}}
}

func (f *fakeRepo) Create(_ context.Context, u User) error {
	if _, ok := f.byEmail[u.Email]; ok {
		return ErrEmailTaken
	}
	f.byID[u.ID] = u
	f.byEmail[u.Email] = u
	return nil
}

func (f *fakeRepo) GetByEmail(_ context.Context, email string) (User, error) {
	u, ok := f.byEmail[email]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}

func (f *fakeRepo) GetByID(_ context.Context, id string) (User, error) {
	u, ok := f.byID[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}

func (f *fakeRepo) Update(_ context.Context, u User) error {
	if _, ok := f.byID[u.ID]; !ok {
		return ErrUserNotFound
	}
	f.byID[u.ID] = u
	f.byEmail[u.Email] = u
	return nil
}

// fakeIssuer mencatat argumen terakhir Issue.
type fakeIssuer struct {
	lastUserID string
	lastTZ     string
	calls      int
}

func (i *fakeIssuer) Issue(userID, timezone string) (string, error) {
	i.lastUserID = userID
	i.lastTZ = timezone
	i.calls++
	return "token-" + userID + "-" + timezone, nil
}

var _ auth.Issuer = (*fakeIssuer)(nil)

func newSvc() (*service, *fakeRepo, *fakeIssuer) {
	repo := newFakeRepo()
	iss := &fakeIssuer{}
	return newService(repo, iss), repo, iss
}

func TestRegister_HappyPath_DefaultTimezone(t *testing.T) {
	svc, repo, iss := newSvc()

	res, err := svc.Register(context.Background(), RegisterRequest{
		Email:       "a@b.c",
		Password:    "rahasia123",
		DisplayName: "A",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if res.Token == "" {
		t.Fatal("token kosong")
	}
	if res.User.Timezone != defaultTimezone {
		t.Fatalf("timezone = %q, mau %q", res.User.Timezone, defaultTimezone)
	}
	if iss.lastTZ != defaultTimezone {
		t.Fatalf("issuer tz = %q, mau %q", iss.lastTZ, defaultTimezone)
	}
	stored := repo.byEmail["a@b.c"]
	if stored.PasswordHash == "" || stored.PasswordHash == "rahasia123" {
		t.Fatalf("password harus di-hash, dapat %q", stored.PasswordHash)
	}
	if stored.ID != res.User.ID || stored.ID == "" {
		t.Fatalf("ID tidak konsisten: %q vs %q", stored.ID, res.User.ID)
	}
}

func TestRegister_EmailTaken(t *testing.T) {
	svc, _, _ := newSvc()
	req := RegisterRequest{Email: "a@b.c", Password: "rahasia123", DisplayName: "A"}

	if _, err := svc.Register(context.Background(), req); err != nil {
		t.Fatalf("register pertama: %v", err)
	}
	_, err := svc.Register(context.Background(), req)
	if !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("err = %v, mau ErrEmailTaken", err)
	}
}

func TestLogin_WrongPasswordAndUnknownEmail_SameError(t *testing.T) {
	svc, _, _ := newSvc()
	svc.Register(context.Background(), RegisterRequest{
		Email: "a@b.c", Password: "rahasia123", DisplayName: "A",
	})

	_, errBadPass := svc.Login(context.Background(), LoginRequest{Email: "a@b.c", Password: "salah"})
	_, errNoUser := svc.Login(context.Background(), LoginRequest{Email: "x@y.z", Password: "rahasia123"})

	if !errors.Is(errBadPass, ErrInvalidCredential) {
		t.Fatalf("password salah -> %v, mau ErrInvalidCredential", errBadPass)
	}
	if !errors.Is(errNoUser, ErrInvalidCredential) {
		t.Fatalf("email tak ada -> %v, mau ErrInvalidCredential", errNoUser)
	}
}

func TestLogin_HappyPath(t *testing.T) {
	svc, _, _ := newSvc()
	svc.Register(context.Background(), RegisterRequest{
		Email: "a@b.c", Password: "rahasia123", DisplayName: "A", Timezone: "Asia/Makassar",
	})

	res, err := svc.Login(context.Background(), LoginRequest{Email: "a@b.c", Password: "rahasia123"})
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if res.User.Timezone != "Asia/Makassar" {
		t.Fatalf("tz = %q", res.User.Timezone)
	}
}

func TestProfile_NoPasswordLeak(t *testing.T) {
	svc, _, _ := newSvc()
	reg, _ := svc.Register(context.Background(), RegisterRequest{
		Email: "a@b.c", Password: "rahasia123", DisplayName: "A",
	})

	got, err := svc.Profile(context.Background(), reg.User.ID)
	if err != nil {
		t.Fatalf("Profile: %v", err)
	}
	// UserResponse tak punya field PasswordHash secara tipe; pastikan juga
	// tak ada string hash yang bocor lewat field lain.
	for _, v := range []string{got.ID, got.Email, got.DisplayName, got.Timezone} {
		if strings.HasPrefix(v, "$2a$") || strings.HasPrefix(v, "$2b$") {
			t.Fatalf("hash bocor di response: %q", v)
		}
	}
}

func TestUpdateProfile_NewTokenCarriesNewTimezone(t *testing.T) {
	svc, _, iss := newSvc()
	reg, _ := svc.Register(context.Background(), RegisterRequest{
		Email: "a@b.c", Password: "rahasia123", DisplayName: "A",
	})
	callsBefore := iss.calls

	newTZ := "Asia/Makassar"
	res, err := svc.UpdateProfile(context.Background(), reg.User.ID, UpdateProfileRequest{Timezone: &newTZ})
	if err != nil {
		t.Fatalf("UpdateProfile: %v", err)
	}
	if res.User.Timezone != newTZ {
		t.Fatalf("response tz = %q, mau %q", res.User.Timezone, newTZ)
	}
	if iss.calls != callsBefore+1 {
		t.Fatalf("Issue dipanggil %d kali, mau %d (token baru wajib)", iss.calls, callsBefore+1)
	}
	if iss.lastTZ != newTZ {
		t.Fatalf("token baru bawa tz %q, mau %q", iss.lastTZ, newTZ)
	}

	// display_name tak dikirim -> tak berubah.
	if res.User.DisplayName != "A" {
		t.Fatalf("display_name berubah jadi %q padahal tak dikirim", res.User.DisplayName)
	}
}

func TestUpdateProfile_UserNotFound(t *testing.T) {
	svc, _, _ := newSvc()
	_, err := svc.UpdateProfile(context.Background(), "tidak-ada", UpdateProfileRequest{})
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("err = %v, mau ErrUserNotFound", err)
	}
}
