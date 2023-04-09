package models

type HasAccessModel struct {
	Phone     string `json:"phone"`
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
