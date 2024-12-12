package migrator

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/foreground-eclipse/song-library/internal/config"
	"github.com/foreground-eclipse/song-library/internal/logger"
	"github.com/golang-migrate/migrate/v4"
)

func Migrate(log *logger.Logger, cfg *config.Config) error {
	const op = "cmd.migrator.migrator"

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer db.Close()

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working dir path: %w", err)
	}
	dir = fmt.Sprintf("%s/internal/migrations", dir)
	dir = filepath.ToSlash(dir)

	m, err := migrate.New(fmt.Sprintf("file://%s", dir), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.LogInfo("trying to apply migrations")

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.LogInfo("no migrations to apply")
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	defer m.Close()

	log.LogInfo("migrations applied")

	return nil
}
