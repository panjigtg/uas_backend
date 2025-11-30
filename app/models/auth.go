package models

type RegistReq struct {
	Username string `json:"username" validate:"required,min=6,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required,min=6"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserProfile struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	FullName    string   `json:"full_name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type LoginResponse struct {
	User         UserProfile `json:"user"`
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
}

type RefreshReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshResp struct {
	AccessToken string `json:"access_token"`
}
