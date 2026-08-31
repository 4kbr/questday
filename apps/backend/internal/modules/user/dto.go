package user

// TODO:
//   type RegisterRequest struct {
//       Email       string `json:"email"       validate:"required,email"`
//       Password    string `json:"password"    validate:"required,min=8"`
//       DisplayName string `json:"displayName" validate:"required"`
//       Timezone    string `json:"timezone"    validate:"required"` // default "Asia/Jakarta"
//   }
//   type LoginRequest struct {
//       Email    string `json:"email"    validate:"required,email"`
//       Password string `json:"password" validate:"required"`
//   }
//   type AuthResponse struct { Token string; User UserResponse }
//   type UserResponse struct { ID string; Email string; DisplayName string; Timezone string }
//   // mapper: toUserResponse(User) UserResponse (tanpa PasswordHash).
