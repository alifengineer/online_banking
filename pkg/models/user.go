package models

type User struct {
	Guid      string `json:"guid"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetUserByIDRequest struct {
	Phone  string `json:"phone"`
	UserId string `json:"user_id"`
}

type GetUserByIDResponse struct {
	User *User `json:"user"`
}

type CreateUserRequest struct {
	User *User `json:"user"`
}
