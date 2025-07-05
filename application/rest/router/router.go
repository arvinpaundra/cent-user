package router

import (
	"github.com/arvinpaundra/cent/user/application/rest/handler"
	"github.com/arvinpaundra/cent/user/application/rest/middleware"
	"github.com/arvinpaundra/cent/user/application/rest/router/auth"
	"github.com/arvinpaundra/cent/user/core/validator"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type router struct {
	g   *gin.Engine
	db  *gorm.DB
	rdb *redis.Client
	vld *validator.Validator
	hdl handler.Handler
}

func Register(g *gin.Engine, db *gorm.DB, rdb *redis.Client, vld *validator.Validator) {
	h := handler.NewHandler(db, rdb, vld)

	g.Use(middleware.Cors())
	g.Use(gin.Recovery())
	g.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/metrics"},
	}))

	r := router{g, db, rdb, vld, h}

	r.public()
	r.private()
}

func (r *router) public() {
	v1 := r.g.Group("/api/v1")

	auth.PublicRoute(v1, r.hdl)
}

func (r *router) private() {
	// v1 := r.g.Group("/api/v1")
}
