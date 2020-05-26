package model

type User struct {
	BaseModel
	Name    string `json:"name"`
	PhoneNo string `json:"phone_no"`
	Email   string `json:"email"`
}
