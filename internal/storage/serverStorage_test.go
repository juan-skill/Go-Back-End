package storage

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/other_project/crockroach/internal/logs"
	"github.com/other_project/crockroach/models"
	"github.com/other_project/crockroach/shared/cockroachdb"
	"github.com/other_project/crockroach/shared/testrandom"

	"github.com/stretchr/testify/require"
)

func InitCockroach() {
	_ = logs.InitLogger()
	// CockroachClient creates a connection with the CockroachDB
	CockroachClient = *cockroachdb.NewSQLClient()

	tx, err := CockroachClient.Begin()
	if err != nil {
		//logs.Log().Errorf("cannot create transaction")
		_ = tx.Rollback()
		return
	}

	sqlStm := `CREATE TABLE IF NOT EXISTS domains (
						id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
						serverChanged bool NOT NULL,
						sslgrade STRING NULL,
						previousslgrade STRING NULL,
						logo STRING NOT NULL,
						title STRING NULL,
						isdown bool NULL,
						creationDate TIMESTAMPTZ NOT NULL DEFAULT (now()),
						updateDate TIMESTAMPTZ NOT NULL DEFAULT (now())
					)
					`

	err = createTables(tx, sqlStm)
	if err != nil {
		_ = tx.Rollback()
		return
	}

	sqlStm = `CREATE TABLE IF NOT EXISTS servers (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				address STRING NOT NULL,
				sslgrade STRING NULL,
				country STRING NULL,
				owner STRING NULL,
				domain_id UUID REFERENCES domains (id) ON DELETE CASCADE,
				creationDate TIMESTAMPTZ NOT NULL DEFAULT (now()),
				updateDate TIMESTAMPTZ NOT NULL DEFAULT (now())
			)
			`

	err = createTables(tx, sqlStm)
	if err != nil {
		_ = tx.Rollback()
		return
	}

	_ = tx.Commit()
}

func createTables(tx *sql.Tx, query string) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	sqlStm, err := tx.PrepareContext(ctx, query)
	if err != nil {
		logs.Log().Errorf(err.Error())
		return err
	}

	result, err := sqlStm.ExecContext(ctx)
	if err != nil {
		logs.Log().Errorf(err.Error())
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		logs.Log().Errorf(err.Error())
		return err
	}

	err = sqlStm.Close()
	if err != nil {
		logs.Log().Errorf(err.Error())
		return err
	}

	return nil
}

func storeServerTest(t *testing.T) *models.Server {
	c := require.New(t)

	InitCockroach()

	domain, err := models.NewDomain(false, false, "A+", testrandom.RandomSSLRating("B"), "https://server.com/icon.png", "Title of the page")
	c.NoError(err)
	c.NotEmpty(domain)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	domain, err = StoreDomain(ctx, domain)
	c.NoError(err)
	c.NotEmpty(domain)

	serverNumber := testrandom.RandomServerNumber()
	var i int64

	for i = 0; i < serverNumber; i++ {
		// create a new Server
		server, err := models.NewServer("server1", testrandom.RandomSSLRating(""), "US", "Amazon.com, Inc.", domain)
		c.NoError(err)
		c.NotNil(server)

		server1, err := StoreServer(ctx, server)
		c.NoError(err)
		c.NotEmpty(server1)

		c.Equal(server.ServerID, server1.ServerID)
		c.Equal(server.Address, server1.Address)
		c.Equal(server.SSLGrade, server1.SSLGrade)
		c.Equal(server.Country, server1.Country)

		domain.Servers = append(domain.Servers, server1)
	}
	c.Equal(domain.DomainID, domain.Servers[0].Domain.DomainID)
	c.Equal(serverNumber, int64(len(domain.Servers)))

	return domain.Servers[0]
}

func TestStoreServer(t *testing.T) {
	storeServerTest(t)
}

func TestStoreServerFailure(t *testing.T) {
	c := require.New(t)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	server, err := StoreServer(ctx, nil)
	c.Error(err)
	c.Nil(server)
	c.EqualError(ErrInvalidServer, err.Error())

	server, err = StoreServer(ctx, server)
	c.Error(err)
	c.Nil(server)
	c.EqualError(ErrInvalidServer, err.Error())
}

func TestGetServer(t *testing.T) {
	c := require.New(t)

	// create a new Server
	server := storeServerTest(t)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	server1, err := GetServer(ctx, server.ServerID)
	c.NoError(err)
	c.NotEmpty(server1)

	c.Equal(server1.ServerID, server.ServerID)
	c.Equal(server1.Address, server.Address)
	c.Equal(server1.SSLGrade, server.SSLGrade)
	c.Equal(server1.Country, server.Country)
	c.WithinDuration(*server.CreationDate, *server1.CreationDate, time.Second)
	c.NotEmpty(server1.Domain.DomainID)
}

func TestGetServerFailure(t *testing.T) {
	c := require.New(t)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	server, err := GetServer(ctx, "")
	c.Error(err)
	c.Nil(server)
	c.EqualError(ErrEmptyServerID, err.Error())

	//
	server, err = GetServer(ctx, "cae0ae1d-45bd-4dda-b938-cfb34569052b")
	c.Error(err)
	c.Nil(server)
	c.EqualError(ErrScanRow, err.Error())
}

func TestUpdateServer(t *testing.T) {
	c := require.New(t)

	// create a new Server
	server := storeServerTest(t)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	server1, err := UpdateServer(ctx, server.ServerID, "A")
	c.NoError(err)
	c.NotEmpty(server1)

	c.Equal(server1.ServerID, server.ServerID)
	c.Equal(server1.Address, server.Address)
	c.NotEqual(server1.SSLGrade, server.SSLGrade)
	c.Equal(server1.SSLGrade, "A")
	c.Equal(server1.Country, server.Country)
	c.WithinDuration(*server.CreationDate, *server1.CreationDate, time.Second)

	c.NotEmpty(server1.Domain.DomainID)

	err = DeleteServer(ctx, server1.ServerID)
	c.NoError(err)

	server1, err = UpdateServer(ctx, server1.ServerID, "B+")
	c.Error(err)
	c.Empty(server1)
	c.EqualError(ErrScanRow, err.Error())
}

func TestUpdateServerFailure(t *testing.T) {
	c := require.New(t)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	server, err := UpdateServer(ctx, "", "")
	c.Error(err)
	c.Nil(server)
	c.EqualError(ErrEmptyServerID, err.Error())

	server, err = UpdateServer(ctx, "cae0ae1d-45bd-4dda-b938-cfb34569052b", "")
	c.Error(err)
	c.Nil(server)
	c.EqualError(ErrEmptySSLGrade, err.Error())

	server, err = UpdateServer(ctx, "", "A")
	c.Error(err)
	c.Nil(server)
	c.EqualError(ErrEmptyServerID, err.Error())

	_, err = UpdateServer(ctx, "cae0ae1d-4dda-b938-cfb34569052b", testrandom.RandomSSLRating(""))
	c.Error(err)
	c.EqualError(ErrInvalidQuery, err.Error())
}

func TestDeleteServer(t *testing.T) {
	c := require.New(t)

	// create a new Server
	server := storeServerTest(t)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	err := DeleteServer(ctx, server.ServerID)
	c.NoError(err)

	server1, err := GetServer(ctx, server.ServerID)
	c.Error(err)
	c.Empty(server1)
	c.EqualError(err, ErrScanRow.Error())
}

func TestDeleteServerFailure(t *testing.T) {
	c := require.New(t)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	err := DeleteServer(ctx, "")
	c.Error(err)
	c.EqualError(ErrEmptyServerID, err.Error())

	err = DeleteServer(ctx, "cae0ae1d-45bd-4dda-b939-cfb34569052b")
	c.Error(err)
	c.EqualError(ErrZeroRowsAffected, err.Error())
}

func TestGetServers(t *testing.T) {
	c := require.New(t)

	for i := 0; i < 10; i++ {
		// create a new Server
		storeServerTest(t)
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	servers, err := GetServers(ctx, "")
	c.NoError(err)

	for _, server := range servers {
		c.NotEmpty(server)
	}

	_, err = GetServers(ctx, "0ae1d-45bd-4dda-b939-cfb34569952b")
	c.Error(err)
}

func BenchmarkStoreServer(b *testing.B) {
	InitCockroach()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	for i := 0; i < b.N; i++ {
		domain, err := models.NewDomain(false, false, "A+", "B", "https://server.com/icon.png", "Title of the page")
		if err != nil {
			b.Fatal(err)
		}

		domain, err = StoreDomain(ctx, domain)
		if err != nil {
			b.Fatal(err)
		}
		// create a new Server
		server, err := models.NewServer("server1", "B", "US", "Amazon.com, Inc.", domain)
		if err != nil {
			b.Fatal(err)
		}

		_, err = StoreServer(ctx, server)
		if err != nil {
			b.Fatal(err)
		}
	}
}
