package migrations

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type migration func(*gorm.DB) error

var (
	migrations = []migration{
		initializeDB,
	}
)

func Migrate(db *gorm.DB) error {
	i := 0
	for i < len(migrations) {
		if err := migrations[i](db); err != nil {
			return errors.Wrap(err, fmt.Sprintf("could not run %dth migration", i))
		}
		i++
	}
	return nil
}
