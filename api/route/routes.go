package route

import (
	"net/http"
	"os"

	"github.com/arvinpaundra/cent/user/api/middleware"
	"github.com/arvinpaundra/cent/user/api/route/auth"
	"github.com/arvinpaundra/cent/user/application/resthttp"
	"github.com/arvinpaundra/cent/user/core/validator"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Routes struct {
	g    *gin.Engine
	db   *gorm.DB
	rdb  *redis.Client
	vld  *validator.Validator
	cont *resthttp.Controller
}

func NewRoutes(g *gin.Engine, db *gorm.DB, rdb *redis.Client, vld *validator.Validator) *Routes {
	controller := resthttp.NewController(db, rdb, vld)

	g.Use(middleware.Cors())
	g.Use(middleware.Logger())
	g.Use(gin.Recovery())
	g.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/metrics"},
	}))

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	return &Routes{
		g:    g,
		db:   db,
		rdb:  rdb,
		vld:  vld,
		cont: controller,
	}
}

func (r *Routes) GatherRoutes() {
	r.public()

	r.private()

	r.internal()
}

func (r *Routes) public() {
	v1 := r.g.Group("/api/v1")

	auth.PublicRoute(v1, r.cont)
}

func (r *Routes) private() {
	v1 := r.g.Group("/api/v1")

	authentication := middleware.NewAuthentication(r.rdb, r.db)

	test := v1.Group("/tests", authentication.Authenticate())

	test.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func (r *Routes) internal() {
}
