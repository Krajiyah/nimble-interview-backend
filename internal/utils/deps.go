package utils

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Deps interface {
	DB() *gorm.DB
	Logger() *logrus.Logger
}

type ProdDeps struct {
	db     *gorm.DB
	logger *logrus.Logger
}

type UnitDeps struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewProdDeps() (Deps, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	db, err := NewProdDB()
	if err != nil {
		return nil, err
	}

	return &ProdDeps{db: db, logger: logger}, nil
}

func (deps *ProdDeps) DB() *gorm.DB           { return deps.db }
func (deps *ProdDeps) Logger() *logrus.Logger { return deps.logger }

func NewUnitDeps() (Deps, string, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	db, fileName, err := NewUnitDB()
	if err != nil {
		return nil, "", err
	}

	return &UnitDeps{db: db, logger: logger}, fileName, nil
}

func (deps *UnitDeps) DB() *gorm.DB           { return deps.db }
func (deps *UnitDeps) Logger() *logrus.Logger { return deps.logger }
