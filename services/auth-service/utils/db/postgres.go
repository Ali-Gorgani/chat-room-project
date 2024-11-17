package db

import (
	"context"
	"fmt"

	"ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/ent"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/logger"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewEntClient),
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewEntClient creates a new ent client and applies migrations.
func NewEntClient(lc fx.Lifecycle, config *configs.Config, logger *logger.Logger) (*ent.Client, error) {
	cfg := Config{
		Host:     config.PSQL.Host,
		Port:     config.PSQL.Port,
		User:     config.PSQL.User,
		Password: config.PSQL.Password,
		DBName:   config.PSQL.Database,
		SSLMode:  config.PSQL.SSLMode,
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	driver, err := sql.Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, err
	}
	client := ent.NewClient(ent.Driver(driver))

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Set up the migration directory for Atlas
			dir, err := migrate.NewLocalDir("utils/ent/migrate/migrations")
			if err != nil {
				return fmt.Errorf("failed to create Atlas migration directory: %w", err)
			}

			// Set up options for schema migration
			opts := []schema.MigrateOption{
				schema.WithDir(dir),                         // specify the migrations directory
				schema.WithMigrationMode(schema.ModeReplay), // replay migration mode
				schema.WithDialect(dialect.Postgres),        // specify the Postgres dialect
				schema.WithFormatter(migrate.DefaultFormatter),
			}

			// Apply migrations
			if err := client.Schema.Create(ctx, opts...); err != nil {
				return fmt.Errorf("failed to apply migrations: %w", err)
			}
			logger.Info("Migrations applied successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}
