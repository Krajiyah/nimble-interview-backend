package middlewares

import (
	"net/http"

	"github.com/Krajiyah/nimble-interview-backend/internal/models"
	"github.com/Krajiyah/nimble-interview-backend/internal/utils"
	"github.com/labstack/echo/v4"
)

const (
	JwtRequestHeader = "X-TOKEN"
	UserContextKey   = "user"
)

func UserAuthMiddleware(deps utils.Deps) echo.MiddlewareFunc {
	return func(f echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := deps.Logger().WithContext(c.Request().Context()).WithField("middleware", "UserAuthMiddleware")

			token := c.Request().Header.Get(JwtRequestHeader)
			if token == "" {
				logger.Warn("missing token in header")
				return c.JSON(http.StatusForbidden, utils.NewErrorResponse(utils.InvalidAuthInfo))
			}

			user, err := utils.ValidateJWT(deps.DB(), token)
			if err != nil {
				logger.Warn("invalid auth token")
				return c.JSON(http.StatusForbidden, utils.NewErrorResponse(utils.InvalidAuthInfo))
			}

			c.Set(UserContextKey, user)
			return f(c)
		}
	}
}

func RequireUser(c echo.Context) *models.User {
	return c.Get(UserContextKey).(*models.User)
}
