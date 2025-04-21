package domain

import (
	"context"
)

type Datastore interface { //nolint:iface // It is not necessary to check for errors at this moment.
	HealthCheckRepository
}

type HealthCheckRepository interface { //nolint:iface // It is not necessary to check for errors at this moment.
	Ping(ctx context.Context) error
}
