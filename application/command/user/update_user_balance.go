package user

type UpdateUserBalance struct {
	UserId int64 `validate:"required"`
	Amount float64 `validate:"required"`
}
