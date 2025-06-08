package rest

import (
	"net/http"

	"github.com/arvinpaundra/cent/user/core/format"
	"github.com/arvinpaundra/cent/user/core/token"
	"github.com/arvinpaundra/cent/user/domain/auth/constant"
	authcmd"github.com/arvinpaundra/cent/user/application/command/auth"
	"github.com/arvinpaundra/cent/user/domain/auth/service"
	authinfra"github.com/arvinpaundra/cent/user/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (cont Controller) Register(c *gin.Context) {
	var payload authcmd.Register

	_ = c.ShouldBindJSON(&payload)

	verrs := cont.validator.Validate(payload)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	handler := service.NewRegisterHandler(
		authinfra.NewUserReaderRepository(cont.db),
		authinfra.NewUserWriterRepository(cont.db),
		authinfra.NewOutboxWriterRepository(cont.db),
		authinfra.NewUnitOfWork(cont.db),
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

func (cont Controller) Login(c *gin.Context) {
	var payload authcmd.Login

	_ = c.ShouldBindJSON(&payload)

	verrs := cont.validator.Validate(payload)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	handler := service.NewLoginHandler(
		authinfra.NewUserReaderRepository(cont.db),
		authinfra.NewSessionWriterRepository(cont.db),
		authinfra.NewUserCacheRepository(cont.rdb),
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

func (cont Controller) RefreshToken(c *gin.Context) {
	var payload authcmd.RefreshToken

	_ = c.ShouldBindJSON(&payload)

	verrs := cont.validator.Validate(payload)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	handler := service.NewRefreshTokenHandler(
		authinfra.NewUserReaderRepository(cont.db),
		authinfra.NewSessionReaderRepository(cont.db),
		authinfra.NewSessionWriterRepository(cont.db),
		authinfra.NewUnitOfWork(cont.db),
		authinfra.NewUserCacheRepository(cont.rdb),
		token.NewJWT(viper.GetString("JWT_SECRET")),
	)

	res, err := handler.Handle(c, payload)
	if err != nil {
		switch err {
		case constant.ErrUserNotFound, constant.ErrSessionNotFound, constant.ErrTokenInvalid:
			c.JSON(http.StatusUnprocessableEntity, format.UnprocessableEntity(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusCreated, format.SuccessCreated("success refresh token", res))
}
