package app

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

const (
	defaultAttempts = 5
	defaultTimeout  = time.Second
)

func init() {
	l, err := logger.NewLogrusLogger("info")
	if err != nil {
		l.Fatal("failed to create new logger", map[string]interface{}{
			"error": err,
		})
	}

	databaseURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(databaseURL) == 0 {
		l.Fatal("migrate: environment variable not declared: PG_URL", map[string]interface{}{})
	}

	var (
		attempts = defaultAttempts
		db       *sql.DB
	)

	for attempts > 0 {
		db, err = sql.Open("postgres", databaseURL)
		if err == nil {
			break
		}

		l.Info("migrate: postgres is trying to connect..", map[string]interface{}{
			"attempts": attempts,
		})

		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		l.Fatal("failed to connect to database", map[string]interface{}{
			"error": err,
		})
	}

	defer db.Close()

	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		l.Fatal("failed to get current directory", map[string]interface{}{
			"error": err,
		})
	}

	migrationsDir := filepath.Join(dir, "migrations")

	if err := goose.Up(db, migrationsDir); err != nil {
		l.Fatal("failed to run migrations", map[string]interface{}{
			"error": err,
		})
	}

	l.Info("the migrations up attempt was successful", map[string]interface{}{})

	// // Rollback
	// if err := goose.Down(db, migrationsDir); err != nil {
	// 	l.Fatal("failed to run migrations", map[string]interface{}{
	// 		"error": err,
	// 	})
	// }

	// l.Info("the migrations down attempt was successful", map[string]interface{}{})
}
