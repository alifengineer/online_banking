package models

type RegisterUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type LoginUserRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UserWithAuth struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GetByCredentialsRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
