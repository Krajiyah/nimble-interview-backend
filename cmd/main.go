package main

import (
	"net/http"

	"github.com/Krajiyah/nimble-interview-backend/internal/migrations"
	"github.com/Krajiyah/nimble-interview-backend/internal/routes"
	"github.com/Krajiyah/nimble-interview-backend/internal/utils"
	"github.com/labstack/echo/v4"
)

func main() {
	deps, err := utils.NewProdDeps()
	checkError(err)
	checkError(migrations.Migrate(deps.DB())) // in production would be a CI/CD step pre-deployment

	server := utils.NewServer(routes.GetAllRoutes(deps))
	server.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, utils.NewSuccessResponse("pong"))
	})

	deps.Logger().Info("Running server...")
	checkError(utils.StartServer(server))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
