package healthcheck

import (
	"context"
	"errors"
	"github.com/suyog1pathak/services/pkg/datastore"
	customerror "github.com/suyog1pathak/services/pkg/errors/service"
	"time"
)

func checkDBconnection() bool {
	ormdb, err := datastore.GetDBConnection()
	if err != nil {
		return false
	}

	db, err := ormdb.DB()
	if err != nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if db.PingContext(ctx) != nil {
		return false
	}
	return true
}

func CheckHealthCheckStatus() error {
	dbCheck := checkDBconnection()

	if !dbCheck {
		return errors.New(customerror.ErrHealthcheckDbFailed)
	}
	return errors.New("")
}
