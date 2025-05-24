package resthttp

import (
	"net/http"

	"github.com/arvinpaundra/cent/user/core/format"
	"github.com/arvinpaundra/cent/user/core/token"
	"github.com/arvinpaundra/cent/user/domain/auth/constant"
	"github.com/arvinpaundra/cent/user/domain/auth/dto/request"
	"github.com/arvinpaundra/cent/user/domain/auth/service"
	"github.com/arvinpaundra/cent/user/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (cont *Controller) Register(c *gin.Context) {
	var payload request.Register

	_ = c.ShouldBindJSON(&payload)

	verrs := cont.validator.Validate(payload)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	handler := service.NewRegisterHandler(
		auth.NewUserReaderRepository(cont.db),
		auth.NewUserWriterRepository(cont.db),
	)

	err := handler.Handle(c, payload)
	if err != nil {
		switch err {
		case constant.ErrEmailAlreadyTaken:
			c.JSON(http.StatusConflict, format.Conflict(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusCreated, format.SuccessCreated("success register", nil))
}

func (cont *Controller) Login(c *gin.Context) {
	var payload request.Login

	_ = c.ShouldBindJSON(&payload)

	verrs := cont.validator.Validate(payload)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	handler := service.NewLoginHandler(
		auth.NewUserReaderRepository(cont.db),
		auth.NewSessionWriterRepository(cont.db),
		token.NewJWT(viper.GetString("JWT_SECRET")),
	)

	res, err := handler.Handle(c, payload)
	if err != nil {
		switch err {
		case constant.ErrUserNotFound, constant.ErrWrongEmailOrPassword:
			c.JSON(http.StatusUnprocessableEntity, format.UnprocessableEntity(constant.ErrWrongEmailOrPassword.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("success login", res))
}
