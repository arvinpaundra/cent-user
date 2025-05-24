package request

type Register struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
	Fullname string `json:"fullname" validate:"required,min=3,max=100"`
}
