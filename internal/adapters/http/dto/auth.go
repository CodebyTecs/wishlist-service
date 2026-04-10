package dto

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}
