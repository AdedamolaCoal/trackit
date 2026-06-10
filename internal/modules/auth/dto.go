package auth

type RegisterRequest struct {
	Name            string `json:"name" validate:"required,min=2,max=100"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	Currency        string `json:"currency" validate:"omitempty,len=3"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,"`
}

type UpdateProfileRequest struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=100"`
	Currency string `json:"currency" validate:"omitempty,len=3"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Currency string `json:"currency"`
}

func ToUserResponse(u *User) UserResponse {
	return UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Currency: u.Currency,
	}
}
