package dbinfra

import (
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/profile/ent"
)

func NewEntClient(dsn string) (*ent.Client, error) {
	drv, err := sql.Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, errors.New("failed to open database connection", "database error", errcode.ErrInternalFailure)
	}
	if drv == nil {
		return nil, errors.New("database driver is nil", "database error", errcode.ErrInternalFailure)
	}

	db := drv.DB()
	if err := db.Ping(); err != nil {
		return nil, errors.New("failed to ping database", "database error", errcode.ErrInternalFailure)
	}

	client := ent.NewClient(ent.Driver(drv))
	if client == nil {
		return nil, errors.New("failed to create ent client", "database error", errcode.ErrInternalFailure)
	}

	return client, nil
}
