package users

import (
	"net/http"
	"time"

	"github.com/Krajiyah/nimble-interview-backend/internal/models"
	"github.com/Krajiyah/nimble-interview-backend/internal/utils"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionDuration = time.Hour * 12
)

type SignupAPI struct {
	deps utils.Deps
}

type authInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

func newAuthResponse(user *models.User) (res utils.Response, _ error) {
	token, err := utils.NewJWT(user, sessionDuration)
	if err != nil {
		return res, err
	}
	return utils.NewSuccessResponse(authResponse{
		User:  user,
		Token: token,
	}), nil
}

func NewSignupAPI(deps utils.Deps) utils.Route {
	return &SignupAPI{deps}
}

func (api *SignupAPI) Method() string                     { return http.MethodPost }
func (api *SignupAPI) Path() string                       { return "/users" }
func (api *SignupAPI) Middlewares() []echo.MiddlewareFunc { return []echo.MiddlewareFunc{} }

func (api *SignupAPI) Handler(c echo.Context) error {
	logger := api.deps.Logger().WithContext(c.Request().Context()).WithField("api", "SignupAPI")

	var input authInput
	if err := c.Bind(&input); err != nil {
		logger.WithError(err).Warn(utils.BadRequestMsg)
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.BadRequestMsg))
	}

	if !checkAuthInput(input) {
		logger.Warn("missing parameters")
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.BadRequestMsg))
	}

	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		logger.WithError(err).Error("could not hash password")
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(utils.InternalServerError))
	}

	user, _ := models.GetUserByUsername(api.deps.DB(), input.Username)
	if user != nil {
		logger.Warn("user already exists")
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.BadRequestMsg))
	}

	db := api.deps.DB().Begin()
	user, err = models.NewUser(db, &models.User{
		Username:     input.Username,
		PasswordHash: passwordHash,
	})
	if err != nil {
		db.Rollback()
		logger.WithError(err).Error("could not create user")
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(utils.InternalServerError))
	}

	res, err := newAuthResponse(user)
	if err != nil {
		db.Rollback()
		logger.WithError(err).Error("could not build auth response")
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(utils.InternalServerError))
	}

	if err := db.Commit().Error; err != nil {
		logger.WithError(err).Error("could not commit transaction for user creation")
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(utils.InternalServerError))
	}

	logger.WithField("id", user.ID).Debug("user created")
	return c.JSON(http.StatusOK, res)
}

type LoginAPI struct {
	deps utils.Deps
}

func NewLoginAPI(deps utils.Deps) utils.Route {
	return &LoginAPI{deps}
}

func (api *LoginAPI) Method() string                     { return http.MethodPost }
func (api *LoginAPI) Path() string                       { return "/users/login" }
func (api *LoginAPI) Middlewares() []echo.MiddlewareFunc { return []echo.MiddlewareFunc{} }

func (api *LoginAPI) Handler(c echo.Context) error {
	logger := api.deps.Logger().WithContext(c.Request().Context()).WithField("api", "LoginAPI")

	var input authInput
	if err := c.Bind(&input); err != nil {
		logger.WithError(err).Warn(utils.BadRequestMsg)
		return c.JSON(http.StatusBadRequest, utils.Response{Error: utils.BadRequestMsg})
	}

	if !checkAuthInput(input) {
		logger.Warn("missing parameters")
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.BadRequestMsg))
	}

	user, err := models.GetUserByUsername(api.deps.DB(), input.Username)
	if err != nil {
		logger.WithError(err).Warn("could not find user w/ username")
		return c.JSON(http.StatusForbidden, utils.Response{Error: utils.InvalidAuthInfo})
	}

	if !checkHash(input.Password, user.PasswordHash) {
		logger.WithError(err).Warn("invalid password")
		return c.JSON(http.StatusForbidden, utils.Response{Error: utils.InvalidAuthInfo})
	}

	logger.WithField("id", user.ID).Debug("user logged in")

	res, err := newAuthResponse(user)
	if err != nil {
		logger.WithError(err).Error("could not build auth response")
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(utils.InternalServerError))
	}

	return c.JSON(http.StatusOK, res)
}

func hashPassword(passwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(b), nil
}

func checkHash(passwd string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd)) == nil
}

func checkAuthInput(input authInput) bool {
	return input.Username != "" && input.Password != ""
}
