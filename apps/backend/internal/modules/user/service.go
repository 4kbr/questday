package user

// service.go = use case auth. Bergantung ke auth.Issuer (port) untuk token dan
// helper hashing password dari platform/auth.
//
// TODO:
//   type service struct { repo Repository; issuer auth.Issuer }
//   func newService(repo Repository, issuer auth.Issuer) *service
//
//   // Register: validasi -> hash password -> simpan user -> terbitkan token.
//   func (s *service) Register(ctx, in RegisterRequest) (AuthResponse, error)
//
//   // Login: cari user by email -> cocokkan password -> terbitkan token.
//   // Kembalikan ErrInvalidCredential (pesan sama untuk email/password salah).
//   func (s *service) Login(ctx, in LoginRequest) (AuthResponse, error)
//
//   func (s *service) Profile(ctx, userID string) (UserResponse, error)
