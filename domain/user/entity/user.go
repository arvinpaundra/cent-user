package entity

type User struct {
	ID       int64
	Email    string
	Password *string
	Fullname string
	Slug     *string
	Image    *string
}
