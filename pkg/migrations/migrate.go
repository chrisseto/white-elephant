//go:generate bash -c "go run github.com/kevinburke/go-bindata/go-bindata -pkg migrations *.sql"
package migrations

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cockroachdb"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"
	"go.uber.org/zap"
)

type logger struct {
	*zap.Logger
}

func (l *logger) Printf(fmt string, args ...interface{}) {
	zap.L().Sugar().Infof(fmt, args...)
}

func (l *logger) Verbose() bool {
	return true
}

func Up(ctx context.Context, db *sql.DB) error {
	// Bind our bindata to migrate's interface
	s, err := bindata.WithInstance(bindata.Resource(AssetNames(), Asset))
	if err != nil {
		return errors.Wrap(err, "creating resource instance")
	}

	crdb, err := cockroachdb.WithInstance(db, &cockroachdb.Config{
		MigrationsTable: "migrations",
		DatabaseName:    "defaultdb",
	})
	if err != nil {
		return errors.Wrap(err, "creating crdb instance")
	}

	m, err := migrate.NewWithInstance("internal", s, "database", crdb)
	if err != nil {
		return errors.Wrap(err, "creating migrator")
	}

	m.Log = &logger{zap.L()}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		m.GracefulStop <- true
	}()

	err = m.Up()

	// This function is idempotent; ignore NoChange errors.
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return errors.Wrap(err, "migrating up")
}
