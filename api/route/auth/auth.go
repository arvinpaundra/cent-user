package auth

import (
	"github.com/arvinpaundra/cent/user/api/middleware"
	restapp "github.com/arvinpaundra/cent/user/application/rest"
	"github.com/gin-gonic/gin"
)

func PublicRoute(g *gin.RouterGroup, cont restapp.Controller) {
	auth := g.Group("/auth")

	auth.POST("/register", cont.Register)
	auth.POST("/login", cont.Login)
	auth.POST("/refresh-tokens", cont.RefreshToken)
}

func PrivateRoute(g *gin.RouterGroup, mdlwr middleware.Authentication, cont restapp.Controller) {
}
