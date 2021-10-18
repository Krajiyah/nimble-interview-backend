package routes

import (
	"github.com/Krajiyah/nimble-interview-backend/internal/routes/users"
	"github.com/Krajiyah/nimble-interview-backend/internal/utils"
)

func GetAllRoutes(deps utils.Deps) []utils.Route {
	return []utils.Route{
		users.NewSignupAPI(deps),
		users.NewLoginAPI(deps),
	}
}
