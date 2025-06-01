package route

import (
	"net/http"

	"github.com/arvinpaundra/cent/user/api/middleware"
	"github.com/arvinpaundra/cent/user/api/route/auth"
	restapp "github.com/arvinpaundra/cent/user/application/rest"
	"github.com/arvinpaundra/cent/user/core/validator"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Routes struct {
	g    *gin.Engine
	db   *gorm.DB
	rdb  *redis.Client
	vld  *validator.Validator
	cont restapp.Controller
}

func NewRoutes(g *gin.Engine, db *gorm.DB, rdb *redis.Client, vld *validator.Validator) *Routes {
	controller := restapp.NewController(db, rdb, vld)

	g.Use(middleware.Cors())
	g.Use(gin.Recovery())
	g.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/metrics"},
	}))

	return &Routes{
		g:    g,
		db:   db,
		rdb:  rdb,
		vld:  vld,
		cont: controller,
	}
}

func (r *Routes) WithPublic() *Routes {
	v1 := r.g.Group("/api/v1")

	auth.PublicRoute(v1, r.cont)

	return r
}

func (r *Routes) WithPrivate() *Routes {
	v1 := r.g.Group("/api/v1")

	authentication := middleware.NewAuthentication(r.rdb, r.db)

	test := v1.Group("/tests", authentication.Authenticate())

	test.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}

func (r *Routes) WithInternal() *Routes {
	return r
}
