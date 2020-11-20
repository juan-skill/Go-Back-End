package storage

import (
	"testing"
	"time"

	"github.com/other_project/crockroach/models"
	"github.com/stretchr/testify/require"
)

func storeDomainTest(t *testing.T) *models.Domain {
	c := require.New(t)

	// create a new Server
	server := storeServerTest(t)

	// create a new Domain
	domain, err := models.NewDomain(server.ServerID, false, false, "A+", "B", "https://server.com/icon.png", "Title of the page")
	c.NoError(err)
	c.NotNil(domain)

	domain1, err := StoreDomain(domain)
	c.NoError(err)
	c.NotEmpty(domain1)
	c.Equal(domain.DomainID, domain1.DomainID)
	c.Equal(domain.ServerChanged, domain1.ServerChanged)
	c.Equal(domain.SSLGrade, domain1.SSLGrade)
	c.Equal(domain.PreviousSSLGrade, domain1.PreviousSSLGrade)
	c.Equal(domain.Logo, domain1.Logo)
	c.Equal(domain.Title, domain1.Title)
	c.Equal(domain.IsDown, domain1.IsDown)
	c.Equal(domain.ServerID, domain1.ServerID)

	return domain
}

func TestStoreDomain(t *testing.T) {
	storeDomainTest(t)
}

func TestStoreDomainFailure(t *testing.T) {
	c := require.New(t)

	domain, err := StoreDomain(nil)
	c.Error(err)
	c.Nil(domain)
	c.EqualError(ErrInvalidDomain, err.Error())

	domain, err = StoreDomain(domain)
	c.Error(err)
	c.Nil(domain)
	c.EqualError(ErrInvalidDomain, err.Error())
}

func TestGetDomain(t *testing.T) {
	c := require.New(t)

	// create a new Server
	domain := storeDomainTest(t)

	domain1, err := GetDomain(domain.DomainID)
	c.NoError(err)
	c.NotEmpty(domain1)

	c.Equal(domain.ServerID, domain1.ServerID)
	c.Equal(domain.DomainID, domain1.DomainID)
	c.Equal(domain.ServerChanged, domain1.ServerChanged)
	c.Equal(domain.SSLGrade, domain1.SSLGrade)
	c.Equal(domain.PreviousSSLGrade, domain1.PreviousSSLGrade)
	c.Equal(domain.Logo, domain1.Logo)
	c.Equal(domain.Title, domain1.Title)
	c.Equal(domain.IsDown, domain1.IsDown)
	c.Equal(domain.ServerID, domain1.ServerID)
	c.WithinDuration(*domain.CreationDate, *domain1.CreationDate, time.Second)
}

func TestGetDomainFailure(t *testing.T) {
	c := require.New(t)

	InitCockroach()

	domain, err := GetDomain("")
	c.Error(err)
	c.Nil(domain)
	c.EqualError(ErrEmptyDomainID, err.Error())

	//
	domain, err = GetDomain("cae0ae1d-45bd-4dda-b938-cfb34569052b")
	c.Error(err)
	c.Nil(domain)
	c.EqualError(ErrScanRow, err.Error())
}

func TestUpdateDomain(t *testing.T) {
	c := require.New(t)

	// create a new Server
	domain := storeDomainTest(t)

	domain1, err := UpdateDomain(domain.DomainID, "A")
	c.NoError(err)
	c.NotEmpty(domain1)

	c.NotEqual(domain.SSLGrade, domain1.SSLGrade)
	c.Equal(domain1.SSLGrade, "A")
	c.Equal(domain.ServerID, domain1.ServerID)
	c.Equal(domain.DomainID, domain1.DomainID)
	c.Equal(domain.ServerChanged, domain1.ServerChanged)
	c.Equal(domain.PreviousSSLGrade, domain1.PreviousSSLGrade)
	c.Equal(domain.Logo, domain1.Logo)
	c.Equal(domain.Title, domain1.Title)
	c.Equal(domain.IsDown, domain1.IsDown)
	c.Equal(domain.ServerID, domain1.ServerID)
	c.WithinDuration(*domain.CreationDate, *domain.CreationDate, time.Second)
}

func TestUpdateDomainFailure(t *testing.T) {
	c := require.New(t)

	InitCockroach()

	domain, err := UpdateDomain("", "")
	c.Error(err)
	c.Nil(domain)
	c.EqualError(ErrEmptyDomainID, err.Error())

	domain, err = UpdateDomain("cae0ae1d-45bd-4dda-b938-cfb34569052b", "")
	c.Error(err)
	c.Nil(domain)
	c.EqualError(ErrEmptySSLGrade, err.Error())

	domain, err = UpdateDomain("", "A")
	c.Error(err)
	c.Nil(domain)
	c.EqualError(ErrEmptyServerID, err.Error())

	domain, err = UpdateDomain("cae0ae1d-45bd-4dda-b938-cfb34569052b", "B")
	c.Error(err)
	c.Nil(domain)
	c.EqualError(ErrScanRow, err.Error())
}

func TestDeleteDomain(t *testing.T) {
	c := require.New(t)

	// create a new domain
	domain := storeDomainTest(t)

	err := DeleteDomain(domain.DomainID)
	c.NoError(err)

	server1, err := GetDomain(domain.DomainID)
	c.Error(err)
	c.Empty(server1)
	c.EqualError(err, ErrScanRow.Error())
}

func TestDeleteDomainFailure(t *testing.T) {
	c := require.New(t)

	InitCockroach()

	err := DeleteDomain("")
	c.Error(err)
	c.EqualError(ErrEmptyDomainID, err.Error())

	err = DeleteDomain("cae0ae1d-45bd-4dda-b939-cfb34569052b")
	c.Error(err)
	c.EqualError(ErrZeroRowsAffected, err.Error())
}

func TestGetDomains(t *testing.T) {
	c := require.New(t)

	for i := 0; i < 10; i++ {
		// create a new Server
		storeDomainTest(t)
	}

	domains, err := GetDomains()
	c.NoError(err)

	for _, domain := range domains {
		c.NotEmpty(domain)
	}
}

func BenchmarkStoreDomain(b *testing.B) {
	InitCockroach()

	for i := 0; i < b.N; i++ {
		// create a new Server
		server, err := models.NewServer("server1", "B", "US", "Amazon.com, Inc.")
		if err != nil {
			b.Fatal(err)
		}

		server, err = StoreServer(server)
		if err != nil {
			b.Fatal(err)
		}

		// create a new Domain
		domain, err := models.NewDomain(server.ServerID, false, false, "A+", "B", "https://server.com/icon.png", "Title of the page")
		if err != nil {
			b.Fatal(err)
		}

		_, err = StoreDomain(domain)
		if err != nil {
			b.Fatal(err)
		}
	}
}
