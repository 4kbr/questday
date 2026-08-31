package user

// TODO:
//   type User struct {
//       ID           string
//       Email        string
//       PasswordHash string      // JANGAN pernah kirim ke response
//       DisplayName  string
//       Timezone     string      // IANA, mis. "Asia/Jakarta"; default saat register
//       CreatedAt    time.Time
//   }
//
//   var (
//       ErrEmailTaken        = errors.New("email sudah terpakai")
//       ErrInvalidCredential = errors.New("email atau password salah")
//       ErrUserNotFound      = errors.New("user tidak ditemukan")
//   )
