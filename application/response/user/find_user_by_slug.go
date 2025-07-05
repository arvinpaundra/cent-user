package user

type FindUserBySlug struct {
	ID       int64   `json:"id"`
	Email    string  `json:"email"`
	Fullname string  `json:"fullname"`
	Key      string  `json:"key"`
	Image    *string `json:"image"`
}
