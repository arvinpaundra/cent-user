package auth

import (
	"github.com/arvinpaundra/cent/user/application/resthttp"
	"github.com/gin-gonic/gin"
)

func PublicRoute(g *gin.RouterGroup, cont *resthttp.Controller) {
	auth := g.Group("/auth")

	auth.POST("/register", cont.Register)
	auth.POST("/login", cont.Login)
}
