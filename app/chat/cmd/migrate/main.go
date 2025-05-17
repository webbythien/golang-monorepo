package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/monorepo/app/chat/internal/models"
	"github.com/monorepo/pkg/migration"
)

// nolint: gosec
func main() {
	stmts, err := migration.Load(models.LoadMigrationModels()...)
	if err != nil {
		log.Fatalf("Failed to load migration: %v", err)
	}

	stmts = installExtensions(stmts)
	_, err = io.WriteString(os.Stdout, stmts)
	if err != nil {
		log.Fatalf("Failed to write migration: %v", err)
	}
}

func installExtensions(stmts string) string {
	stmts = fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS pg_trgm;\n%s", stmts)
	return stmts
}
