package middleware

import (
	"net/http"
	"strings"

	"github.com/arvinpaundra/cent/user/core/format"
	"github.com/arvinpaundra/cent/user/core/token"
	"github.com/arvinpaundra/cent/user/domain/auth/service"
	"github.com/arvinpaundra/cent/user/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Authentication struct {
	rdb *redis.Client
	db  *gorm.DB
}

func NewAuthentication(rdb *redis.Client, db *gorm.DB) Authentication {
	return Authentication{
		rdb: rdb,
		db:  db,
	}
}

func (m Authentication) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")

		if bearerToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, format.Unauthorized("unauthenticated user"))
			return
		}

		handler := service.NewAuthenticateHandler(
			auth.NewUserReaderRepository(m.db),
			auth.NewUserCacheRepository(m.rdb),
			token.NewJWT(viper.GetString("JWT_SECRET")),
		)

		sanitizeToken := strings.Replace(bearerToken, "Bearer ", "", 1)

		res, err := handler.Handle(c, sanitizeToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, format.Unauthorized("unauthenticated user"))
			return
		}

		c.Set("session", res)

		c.Next()

	}
}
