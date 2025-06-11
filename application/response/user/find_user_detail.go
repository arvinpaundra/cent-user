package user

type FindUserDetail struct {
	ID       int64   `json:"id"`
	Email    string  `json:"email"`
	Fullname string  `json:"fullname"`
	Image    *string `json:"image"`
}
