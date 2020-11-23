package storage

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/other_project/crockroach/internal/logs"
	"github.com/other_project/crockroach/models"
	"github.com/other_project/crockroach/shared/cockroachdb"
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

	domain, err := models.NewDomain(false, false, "A+", "B", "https://server.com/icon.png", "Title of the page")
	c.NoError(err)
	c.NotEmpty(domain)

	domain, err = StoreDomain(domain)
	c.NoError(err)
	c.NotEmpty(domain)

	// create a new Server
	server, err := models.NewServer("server1", "B", "US", "Amazon.com, Inc.", domain)
	c.NoError(err)
	c.NotNil(server)

	server1, err := StoreServer(server)
	c.NoError(err)
	c.NotEmpty(server1)

	c.Equal(server.ServerID, server1.ServerID)
	c.Equal(server.Address, server1.Address)
	c.Equal(server.SSLGrade, server1.SSLGrade)
	c.Equal(server.Country, server1.Country)

	domain.Servers = append(domain.Servers, server1)

	c.Equal(domain.DomainID, domain.Servers[0].Domain.DomainID)

	return server1
}

func TestStoreServer(t *testing.T) {
	storeServerTest(t)
}

func TestStoreServerFailure(t *testing.T) {
	c := require.New(t)

	server, err := StoreServer(nil)
	c.Error(err)
	c.Nil(server)
	c.EqualError(ErrInvalidServer, err.Error())

	server, err = StoreServer(server)
	c.Error(err)
	c.Nil(server)
	c.EqualError(ErrInvalidServer, err.Error())
}

func TestGetServer(t *testing.T) {
	c := require.New(t)

	// create a new Server
	server := storeServerTest(t)

	server1, err := GetServer(server.ServerID)
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

	server, err := GetServer("")
	c.Error(err)
	c.Nil(server)
	c.EqualError(ErrEmptyServerID, err.Error())

	//
	server, err = GetServer("cae0ae1d-45bd-4dda-b938-cfb34569052b")
	c.Error(err)
	c.Nil(server)
	c.EqualError(ErrScanRow, err.Error())
}

func TestUpdateServer(t *testing.T) {
	c := require.New(t)

	// create a new Server
	server := storeServerTest(t)

	server1, err := UpdateServer(server.ServerID, "A")
	c.NoError(err)
	c.NotEmpty(server1)

	c.Equal(server1.ServerID, server.ServerID)
	c.Equal(server1.Address, server.Address)
	c.NotEqual(server1.SSLGrade, server.SSLGrade)
	c.Equal(server1.SSLGrade, "A")
	c.Equal(server1.Country, server.Country)
	c.WithinDuration(*server.CreationDate, *server1.CreationDate, time.Second)

	c.NotEmpty(server1.Domain.DomainID)
}

func TestUpdateServerFailure(t *testing.T) {
	c := require.New(t)

	server, err := UpdateServer("", "")
	c.Error(err)
	c.Nil(server)
	c.EqualError(ErrEmptyServerID, err.Error())

	server, err = UpdateServer("cae0ae1d-45bd-4dda-b938-cfb34569052b", "")
	c.Error(err)
	c.Nil(server)
	c.EqualError(ErrEmptySSLGrade, err.Error())

	server, err = UpdateServer("", "A")
	c.Error(err)
	c.Nil(server)
	c.EqualError(ErrEmptyServerID, err.Error())
}

func TestDeleteServer(t *testing.T) {
	c := require.New(t)

	// create a new Server
	server := storeServerTest(t)

	err := DeleteServer(server.ServerID)
	c.NoError(err)

	server1, err := GetServer(server.ServerID)
	c.Error(err)
	c.Empty(server1)
	c.EqualError(err, ErrScanRow.Error())
}

func TestDeleteServerFailure(t *testing.T) {
	c := require.New(t)

	err := DeleteServer("")
	c.Error(err)
	c.EqualError(ErrEmptyServerID, err.Error())

	err = DeleteServer("cae0ae1d-45bd-4dda-b939-cfb34569052b")
	c.Error(err)
	c.EqualError(ErrZeroRowsAffected, err.Error())
}

func TestGetServers(t *testing.T) {
	c := require.New(t)

	for i := 0; i < 10; i++ {
		// create a new Server
		storeServerTest(t)
	}

	servers, err := GetServers()
	c.NoError(err)

	for _, server := range servers {
		c.NotEmpty(server)
	}

	servers = nil
	c.Nil(servers)
}

func BenchmarkStoreServer(b *testing.B) {
	InitCockroach()

	for i := 0; i < b.N; i++ {
		InitCockroach()

		domain, err := models.NewDomain(false, false, "A+", "B", "https://server.com/icon.png", "Title of the page")
		if err != nil {
			b.Fatal(err)
		}

		domain, err = StoreDomain(domain)
		if err != nil {
			b.Fatal(err)
		}
		// create a new Server
		server, err := models.NewServer("server1", "B", "US", "Amazon.com, Inc.", domain)
		if err != nil {
			b.Fatal(err)
		}

		_, err = StoreServer(server)
		if err != nil {
			b.Fatal(err)
		}
	}
}
