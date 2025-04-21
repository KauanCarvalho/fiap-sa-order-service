package datastore

import (
	"context"
)

func (ds *datastore) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultConnectionTimeout)
	defer cancel()

	sqlDB, err := ds.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}
