package storage

import (
	"context"
	"fmt"
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

func getDomain1(t *testing.T) *models.Domain {
	c := require.New(t)

	domain, err := models.NewDomain(false, false, testrandom.RandomNameDomain(), "", "", "https://server.com/icon.png", "Title of the page")
	c.NoError(err)

	serverNumber := testrandom.RandomServerNumber()
	var i int64

	for i = 0; i < serverNumber; i++ {
		server, erro := models.NewServer("server1", testrandom.RandomSSLRating(""), "US", "Amazon.com, Inc.", domain)
		c.NoError(erro)
		c.NotNil(server)

		domain.Servers = append(domain.Servers, server)
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	domain, err = StoreDomain(ctx, domain)
	c.NoError(err)

	serversNumber := len(domain.Servers)

	for i := 0; i < serversNumber; i++ {
		server := domain.Servers[i]

		_, erro := StoreServer(ctx, server, domain)
		c.NoError(erro)
	}

	return domain
}

func TestGetRecordByName(t *testing.T) {
	c := require.New(t)

	InitCockroach()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err := ReloadRecord(ctx)
	c.NoError(err)

	domainsNumber := testrandom.RandomServerNumber()
	var i int64

	domains := []*models.Domain{}

	for i = 0; i < domainsNumber; i++ {
		/////////////////////////////////////////////////////////////////////////////////////////////////////////
		domain := getDomain1(t)
		/////////////////////////////////////////////////////////////////////////////////////////////////////////

		record, erro := NewRecord(domain)
		c.NoError(erro)
		c.NotEmpty(record)

		domains = append(domains, domain)
	}

	//fmt.Println(domains[0].DomainName)
	records, err := GetRecordByName(domains[0])
	c.NoError(err)
	c.NotEmpty(records)

	records, err = GetRecordByName(nil)
	c.EqualError(models.ErrEmptyDomain, err.Error())
	c.Nil(records)
}

func TestGetLastDomain(t *testing.T) {
	c := require.New(t)

	InitCockroach()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err := ReloadRecord(ctx)
	c.NoError(err)

	domainsNumber := testrandom.RandomServerNumber()
	var i int64

	domains := []*models.Domain{}

	for i = 0; i < domainsNumber; i++ {
		/////////////////////////////////////////////////////////////////////////////////////////////////////////
		domain := getDomain1(t)
		/////////////////////////////////////////////////////////////////////////////////////////////////////////

		record, erro := NewRecord(domain)
		c.NoError(erro)
		c.NotEmpty(record)

		domains = append(domains, domain)
	}

	//fmt.Println(domains[0].DomainName)
	records, err := GetRecordByName(domains[0])
	c.NoError(err)
	c.NotEmpty(records)

	recordsList := GetLastDomain()
	c.NotEmpty(recordsList)

	fmt.Println("get last record ", recordsList)
}
