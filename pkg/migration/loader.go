package migration

import (
	"errors"
	"fmt"

	"ariga.io/atlas-provider-gorm/gormschema"
	"gorm.io/gorm"
)

var ErrLoadGormSchema = errors.New("failed to load gorm schema")

func Load(mdls ...any) (string, error) {
	// Load all migration files here
	stmts, err := gormschema.New("postgres", gormschema.WithConfig(&gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})).
		Load(mdls...)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrLoadGormSchema, err)
	}
	return stmts, nil
}
