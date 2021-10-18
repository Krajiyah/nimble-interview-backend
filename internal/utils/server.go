package utils

import (
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	BadRequestMsg       = "invalid parameters"
	InvalidAuthInfo     = "invalid auth info"
	InternalServerError = "internal server error"
)

type Route interface {
	Method() string
	Path() string
	Handler(echo.Context) error
	Middlewares() []echo.MiddlewareFunc
}

type Response struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
}

func NewSuccessResponse(result interface{}) Response {
	return Response{Result: result}
}

func NewErrorResponse(msg string) Response {
	return Response{Error: msg}
}

func NewServer(routes []Route) *echo.Echo {
	e := echo.New()
	i := 0
	for i < len(routes) {
		route := routes[i]
		e.Add(route.Method(), route.Path(), route.Handler, route.Middlewares()...)
		i++
	}
	return e
}

func StartServer(e *echo.Echo) error {
	return e.Start(":" + os.Getenv("PORT"))
}

func NewPaginator(c echo.Context) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
