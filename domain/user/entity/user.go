package entity

import "github.com/arvinpaundra/cent/user/core/trait"

type User struct {
	trait.Updateable

	ID       int64
	Email    string
	Fullname string
	Balance  float64
	Currency string
	Password *string
	Slug     *string
	Image    *string
}

func (e *User) IsNew() bool {
	return e.ID <= 0
}

func (e *User) UpdateBalance(amount float64) {
	e.Balance += amount
}
