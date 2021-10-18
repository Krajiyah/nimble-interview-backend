package messages

import (
	"net/http"

	"github.com/Krajiyah/nimble-interview-backend/internal/models"
	"github.com/Krajiyah/nimble-interview-backend/internal/routes/middlewares"
	"github.com/Krajiyah/nimble-interview-backend/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SendMessageAPI struct {
	deps utils.Deps
}

type messageInput struct {
	Data string `json:"data"`
}

func NewSendMessageAPI(deps utils.Deps) utils.Route {
	return &SendMessageAPI{deps}
}

func (api *SendMessageAPI) Method() string { return http.MethodPost }
func (api *SendMessageAPI) Path() string   { return "/messages" }
func (api *SendMessageAPI) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{middlewares.UserAuthMiddleware(api.deps)}
}

func (api *SendMessageAPI) Handler(c echo.Context) error {
	user := middlewares.RequireUser(c)
	logger := api.deps.Logger().WithContext(c.Request().Context()).WithFields(logrus.Fields{"api": "SendMessageAPI", "user": user})

	var input messageInput
	if err := c.Bind(&input); err != nil {
		logger.WithError(err).Warn(utils.BadRequestMsg)
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.BadRequestMsg))
	}

	if input.Data == "" {
		logger.Warn("missing parameters")
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.BadRequestMsg))
	}

	db := api.deps.DB().Begin()
	message, err := models.NewMessage(db, &models.Message{Data: input.Data, Username: user.Username})
	if err != nil {
		db.Rollback()
		logger.WithError(err).Error("could not create message")
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(utils.InternalServerError))
	}

	if err := db.Commit().Error; err != nil {
		logger.WithError(err).Error("could not commit transaction for user creation")
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(utils.InternalServerError))
	}

	logger.WithField("message", message).Debug("message created")
	return c.JSON(http.StatusOK, utils.NewSuccessResponse(message))
}

type GetMessagesAPI struct {
	deps utils.Deps
}

func NewGetMessagesAPI(deps utils.Deps) utils.Route {
	return &GetMessagesAPI{deps}
}

func (api *GetMessagesAPI) Method() string { return http.MethodGet }
func (api *GetMessagesAPI) Path() string   { return "/messages" }
func (api *GetMessagesAPI) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{middlewares.UserAuthMiddleware(api.deps)}
}

func (api *GetMessagesAPI) Handler(c echo.Context) error {
	user := middlewares.RequireUser(c)
	logger := api.deps.Logger().WithContext(c.Request().Context()).WithFields(logrus.Fields{"api": "GetMessagesAPI", "user": user})

	messages := []models.Message{}
	api.deps.DB().Session(&gorm.Session{QueryFields: true}).Scopes(utils.NewPaginator(c)).Find(&messages)

	logger.WithField("messageCount", len(messages)).Debug("got messages")
	return c.JSON(http.StatusOK, utils.NewSuccessResponse(messages))
}
