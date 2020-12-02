package storage

import (
	"context"
	"testing"
	"time"

	"github.com/other_project/crockroach/models"
	"github.com/other_project/crockroach/shared/testrandom"
	"github.com/stretchr/testify/require"
)

func storeDomainTest(t *testing.T) *models.Domain {
	c := require.New(t)

	InitCockroach()
	// create a new Server
	//server := storeServerTest(t)

	// create a new Domain

	domain, err := models.NewDomain(false, false, "google.com", "A+", testrandom.RandomSSLRating("A+"), "https://server.com/icon.png", "Title of the page")
	c.NoError(err)
	c.NotNil(domain)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	domain1, err := StoreDomain(ctx, domain)
	c.NoError(err)
	c.NotEmpty(domain1)

	c.NotEmpty(domain1)
	c.Equal(domain.DomainID, domain1.DomainID)
	c.Equal(domain.ServerChanged, domain1.ServerChanged)
	c.Equal(domain.SSLGrade, domain1.SSLGrade)
	c.Equal(domain.PreviousSSLGrade, domain1.PreviousSSLGrade)
	c.Equal(domain.Logo, domain1.Logo)
	c.Equal(domain.Title, domain1.Title)
	c.Equal(domain.IsDown, domain1.IsDown)

	return domain1
}

func TestStoreDomain(t *testing.T) {
	storeDomainTest(t)
}

func TestStoreDomainFailure(t *testing.T) {
	c := require.New(t)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	domain, err := StoreDomain(ctx, nil)
	c.Error(err)
	c.Nil(domain)
	c.EqualError(ErrInvalidDomain, err.Error())

	domain, err = StoreDomain(ctx, domain)
	c.Error(err)
	c.Nil(domain)
	c.EqualError(ErrInvalidDomain, err.Error())
}

func TestGetDomain(t *testing.T) {
	c := require.New(t)

	// create a new Server
	domain := storeDomainTest(t)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	domain1, err := GetDomain(ctx, domain.DomainID)
	c.NoError(err)
	c.NotEmpty(domain1)

	//c.Equal(domain.ServerID, domain1.ServerID)
	c.Equal(domain.DomainID, domain1.DomainID)
	c.Equal(domain.ServerChanged, domain1.ServerChanged)
	c.Equal(domain.SSLGrade, domain1.SSLGrade)
	c.Equal(domain.PreviousSSLGrade, domain1.PreviousSSLGrade)
	c.Equal(domain.Logo, domain1.Logo)
	c.Equal(domain.Title, domain1.Title)
	c.Equal(domain.IsDown, domain1.IsDown)
	//c.Equal(domain.ServerID, domain1.ServerID)
	c.WithinDuration(*domain.CreationDate, *domain1.CreationDate, time.Second)
}

func TestGetDomainFailure(t *testing.T) {
	c := require.New(t)

	InitCockroach()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	domain, err := GetDomain(ctx, "")
	c.Error(err)
	c.Nil(domain)
	c.EqualError(ErrEmptyDomainID, err.Error())

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	domain, err = GetDomain(ctx, "cae0ae1d-45bd-4dda-b938-cfb34569052b")
	c.Error(err)
	c.Nil(domain)
	c.EqualError(ErrScanRow, err.Error())
}

func TestUpdateDomain(t *testing.T) {
	c := require.New(t)

	// create a new Server
	domain := storeDomainTest(t)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	domain1, err := UpdateDomain(ctx, "A", "", domain, false)
	c.NoError(err)
	c.NotEmpty(domain1)

	c.NotEqual(domain.SSLGrade, domain1.SSLGrade)
	c.Equal(domain1.SSLGrade, "A")
	//c.Equal(domain.ServerID, domain1.ServerID)
	c.Equal(domain.DomainID, domain1.DomainID)
	c.Equal(domain.ServerChanged, domain1.ServerChanged)
	c.Equal(domain.PreviousSSLGrade, domain1.PreviousSSLGrade)
	c.Equal(domain.Logo, domain1.Logo)
	c.Equal(domain.Title, domain1.Title)
	c.Equal(domain.IsDown, domain1.IsDown)
	//c.Equal(domain.ServerID, domain1.ServerID)
	c.WithinDuration(*domain.CreationDate, *domain.CreationDate, time.Second)
}

func TestUpdateDomainFailure(t *testing.T) {
	c := require.New(t)

	InitCockroach()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	domain, err := UpdateDomain(ctx, "", "", nil, false)
	c.Error(err)
	c.Nil(domain)
	c.EqualError(ErrEmptyDomain, err.Error())

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	domain, err = UpdateDomain(ctx, "", "A", nil, false)
	c.Error(err)
	c.Nil(domain)
	c.EqualError(ErrEmptyDomain, err.Error())
}

func TestDeleteDomain(t *testing.T) {
	c := require.New(t)

	// create a new domain
	domain := storeDomainTest(t)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	err := DeleteDomain(ctx, domain.DomainID)
	c.NoError(err)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	server1, err := GetDomain(ctx, domain.DomainID)
	c.Error(err)
	c.Empty(server1)
	c.EqualError(err, ErrScanRow.Error())
}

func TestDeleteDomainFailure(t *testing.T) {
	c := require.New(t)

	InitCockroach()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	err := DeleteDomain(ctx, "")
	c.Error(err)
	c.EqualError(ErrEmptyDomainID, err.Error())

	err = DeleteDomain(ctx, "cae0ae1d-45bd-4dda-b939-cfb34569052b")
	c.Error(err)
	c.EqualError(ErrZeroRowsAffected, err.Error())

	domains, err := GetDomains(ctx, "")
	c.NoError(err)

	for _, domain := range domains {
		err = DeleteDomain(ctx, domain.DomainID)
		c.Nil(err)
	}

	_, err = GetDomains(ctx, "")
	c.NoError(err)
}

func TestGetDomains(t *testing.T) {
	c := require.New(t)

	for i := 0; i < 10; i++ {
		// create a new Server
		storeDomainTest(t)
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	domains, err := GetDomains(ctx, "")
	c.NoError(err)

	for _, domain := range domains {
		err = DeleteDomain(ctx, domain.DomainID)
		c.Nil(err)
		c.NotEmpty(domain)
	}
}

func BenchmarkStoreDomain(b *testing.B) {
	InitCockroach()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	for i := 0; i < b.N; i++ {
		// create a new Domain
		domain, err := models.NewDomain(false, false, "google.com", "A+", "B", "https://server.com/icon.png", "Title of the page")
		if err != nil {
			b.Fatal(err)
		}

		_, err = StoreDomain(ctx, domain)
		if err != nil {
			b.Fatal(err)
		}
	}
}
