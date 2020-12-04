package storage

import (
	"context"
	"testing"
	"time"

	"github.com/other_project/crockroach/models"
	"github.com/other_project/crockroach/shared/testrandom"
	"github.com/stretchr/testify/require"
)

func newDomainTest(t *testing.T) *models.Domain {
	c := require.New(t)

	InitCockroach()

	domain, err := models.NewDomain(false, false, "google.com", "A+", testrandom.RandomSSLRating("B"), "https://server.com/icon.png", "Title of the page")
	c.NoError(err)
	c.NotEmpty(domain)
	c.Equal("google.com", domain.DomainName)

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

		server1, err := StoreServer(ctx, server, domain)
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

	return domain
}

func TestNewDomain(t *testing.T) {
	newDomainTest(t)
}

func TestNewRecord(t *testing.T) {
	c := require.New(t)

	domain := newDomainTest(t)

	record, err := NewRecord(domain)
	c.NoError(err)
	c.NotEmpty(record)

	serverNumber := testrandom.RandomServerNumber()
	var i int64

	for i = 0; i < serverNumber; i++ {
		domain = newDomainTest(t)
		record, err = NewRecord(domain)
		c.NoError(err)
		c.NotEmpty(record)
	}

	_, err = NewRecord(nil)
	c.EqualError(models.ErrEmptyDomain, err.Error())
}

func TestReloadRecord(t *testing.T) {
	c := require.New(t)

	InitCockroach()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	objects, err := ReloadRecord(ctx)
	c.NoError(err)
	c.NotEmpty(objects)
}

func TestGetRecordByName(t *testing.T) {
	c := require.New(t)

	InitCockroach()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err := ReloadRecord(ctx)
	c.NoError(err)

	domain := newDomainTest(t)

	record, err := NewRecord(domain)
	c.NoError(err)
	c.NotEmpty(record)

	records, err := GetRecordByName(domain)
	c.NoError(err)
	c.NotEmpty(records)

	records, err = GetRecordByName(nil)
	c.EqualError(models.ErrEmptyDomain, err.Error())
	c.Nil(records)
}
