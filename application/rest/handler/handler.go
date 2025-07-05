package handler

import (
	"github.com/arvinpaundra/cent/user/core/validator"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	db        *gorm.DB
	rdb       *redis.Client
	validator *validator.Validator
}

func NewHandler(db *gorm.DB, rdb *redis.Client, validator *validator.Validator) Handler {
	return Handler{
		db:        db,
		rdb:       rdb,
		validator: validator,
	}
}
