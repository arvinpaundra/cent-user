package rest

import (
	"github.com/arvinpaundra/cent/user/core/validator"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Controller struct {
	db        *gorm.DB
	rdb       *redis.Client
	validator *validator.Validator
}

func NewController(db *gorm.DB, rdb *redis.Client, validator *validator.Validator) Controller {
	return Controller{
		db:        db,
		rdb:       rdb,
		validator: validator,
	}
}
