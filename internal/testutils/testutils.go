package testutils

import (
	"github.com/Krajiyah/nimble-interview-backend/internal/migrations"
	"github.com/Krajiyah/nimble-interview-backend/internal/utils"
)

func NewUnitDeps() (utils.Deps, string, error) {
	unitDeps, fileName, err := utils.NewUnitDeps()
	if err != nil {
		return nil, "", err
	}

	if err := migrations.Migrate(unitDeps.DB()); err != nil {
		return nil, "", err
	}

	return unitDeps, fileName, nil
}
