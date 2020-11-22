package cockroachdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// we do not call any function of lib/pq directly in the code
	_ "github.com/lib/pq"
	"github.com/other_project/crockroach/internal/logs"
	"github.com/other_project/crockroach/shared/env"
)

// NewSQLClient function will create a new sql client
func NewSQLClient() *sql.DB {
	address := env.GetString("COCKROACH_ADDRESS", "test@localhost:26257")
	db := env.GetString("COCKROACH_DATABASE", "testdb")

	//"postgresql://test@localhost:26257/testdb?sslmode=disable"
	// postgres://<username>:<password>@<host>:<port>/<database>?<parameters>

	database, err := sql.Open(
		"postgres",
		fmt.Sprintf("postgresql://%s/%s?sslmode=disable", address, db),
	)
	if err != nil {
		logs.Log().Errorf("Error connecting to the database:  %s", err.Error())
		return nil
	}

	/* Control number of session: problem chech coverage
	database.SetMaxOpenConns(20)
	database.SetMaxIdleConns(20)
	database.SetConnMaxLifetime(time.Minute * 5)
	*/

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	err = database.PingContext(ctx)
	if err != nil {
		logs.Log().Errorf("cannot connect to mysql: %s", err.Error())
		return nil
	}

	return database
}
