package models

type RegistReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	User  			UserProfile `json:"user"`
	Token 			string      `json:"token"`
	RefreshToken 	string      `json:"refresh_token"`
}