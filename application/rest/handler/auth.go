package handler

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

func (h Handler) Register(c *gin.Context) {
	var payload authcmd.Register

	_ = c.ShouldBindJSON(&payload)

	verrs := h.validator.Validate(payload)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	svc := service.NewRegister(
		authinfra.NewUserReaderRepository(h.db),
		authinfra.NewUserWriterRepository(h.db),
		authinfra.NewOutboxWriterRepository(h.db),
		authinfra.NewUnitOfWork(h.db),
	)

	err := svc.Exec(c, payload)
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

func (h Handler) Login(c *gin.Context) {
	var payload authcmd.Login

	_ = c.ShouldBindJSON(&payload)

	verrs := h.validator.Validate(payload)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	svc := service.NewLogin(
		authinfra.NewUserReaderRepository(h.db),
		authinfra.NewSessionWriterRepository(h.db),
		authinfra.NewUserCacheRepository(h.rdb),
		token.NewJWT(viper.GetString("JWT_SECRET")),
	)

	res, err := svc.Exec(c, payload)
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

func (h Handler) RefreshToken(c *gin.Context) {
	var payload authcmd.RefreshToken

	_ = c.ShouldBindJSON(&payload)

	verrs := h.validator.Validate(payload)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	svc := service.NewRefreshToken(
		authinfra.NewUserReaderRepository(h.db),
		authinfra.NewSessionReaderRepository(h.db),
		authinfra.NewSessionWriterRepository(h.db),
		authinfra.NewUnitOfWork(h.db),
		authinfra.NewUserCacheRepository(h.rdb),
		token.NewJWT(viper.GetString("JWT_SECRET")),
	)

	res, err := svc.Exec(c, payload)
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
