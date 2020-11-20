package storage

import (
	"testing"
	"time"

	"github.com/other_project/crockroach/models"
	"github.com/other_project/crockroach/shared/cockroachdb"
	"github.com/stretchr/testify/require"
)

func InitCockroach() {
	// CockroachClient creates a connection with the CockroachDB
	CockroachClient = *cockroachdb.NewSQLClient()
}

func storeServerTest(t *testing.T) *models.Server {
	c := require.New(t)

	InitCockroach()

	// create a new Server
	server, err := models.NewServer("server1", "B", "US", "Amazon.com, Inc.")
	c.NoError(err)
	c.NotNil(server)

	server1, err := StoreServer(server)
	c.NoError(err)
	c.NotEmpty(server1)
	c.Equal(server.ServerID, server1.ServerID)
	c.Equal(server.Address, server1.Address)
	c.Equal(server.SSLGrade, server1.SSLGrade)
	c.Equal(server.Country, server1.Country)

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
}

func BenchmarkStoreServer(b *testing.B) {
	InitCockroach()

	for i := 0; i < b.N; i++ {
		// create a new Server
		server, err := models.NewServer("server1", "B", "US", "Amazon.com, Inc.")
		if err != nil {
			b.Fatal(err)
		}

		_, err = StoreServer(server)
		if err != nil {
			b.Fatal(err)
		}
	}
}
