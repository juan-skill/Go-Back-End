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

var (
	// UserName  data to connect with Database Cliente
	UserName = env.GetString("DATABASE_USERNAME", "test")
	// HostName data to connect with Database Cliente
	HostName = env.GetString("DATABASE_HOSTNAME", "localhost")
	// Port data to connect with Database Cliente
	Port = env.GetString("DATABASE_PORT", "26257")
	// DatabaseName data to connect with Database Cliente
	DatabaseName = env.GetString("DATABASE_NAME", "testdb")
	// DriverName data to connect with Database Cliente
	DriverName = env.GetString("DATABASE_DRIVER", "postgres")
)

// NewSQLClient function will create a new sql client
func NewSQLClient() *sql.DB {
	database, err := sql.Open(
		DriverName,
		fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=disable", UserName, HostName, Port, DatabaseName),
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
