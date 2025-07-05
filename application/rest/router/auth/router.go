package auth

import (
	"github.com/arvinpaundra/cent/user/application/rest/handler"
	"github.com/arvinpaundra/cent/user/application/rest/middleware"
	"github.com/gin-gonic/gin"
)

func PublicRoute(g *gin.RouterGroup, h handler.Handler) {
	auth := g.Group("/auth")

	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh-tokens", h.RefreshToken)
}

func PrivateRoute(g *gin.RouterGroup, mdlwr middleware.Authentication, h handler.Handler) {
}
