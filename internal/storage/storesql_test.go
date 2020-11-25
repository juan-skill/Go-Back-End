package storage

import (
	"context"
	"testing"
	"time"

	"github.com/other_project/crockroach/models"
	"github.com/other_project/crockroach/shared/testrandom"

	"github.com/stretchr/testify/require"
)

func getNewDomain(t *testing.T) *models.Domain {
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

	return domain
}

func TestTransfertx(t *testing.T) {
	c := require.New(t)

	InitCockroach()

	store := NewStore()
	domain1 := getNewDomain(t)
	arg := TransferTxParams{
		FromDomainID: domain1.DomainID,
	}

	n := 5

	errs := make(chan error)
	results := make(chan TransferTxResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		go func() {
			ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelfunc()

			result, err := store.TransferTx(ctx, arg)

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		c.NoError(err)

		result := <-results
		c.NotEmpty(result)

		// check domain
		domain := result.FromDomain
		c.NotEmpty(domain)
		c.Equal(domain1.DomainID, result.FromDomain.DomainID)
		c.Equal(len(domain1.Servers), len(result.FromDomain.Servers))
		c.NotEqual(domain1.SSLGrade, result.ToDomain.SSLGrade)
		c.NotEqual(domain1.UpdateDate, result.ToDomain.UpdateDate)
	}
}

func TestTransfertxFailure(t *testing.T) {
	c := require.New(t)

	InitCockroach()

	store := NewStore()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	arg := TransferTxParams{
		FromDomainID: "",
	}

	_, err := store.TransferTx(ctx, arg)
	c.Error(err)

	arg.FromDomainID = "cae0ae1d-45bd-4dda-b939-cfb34569052b"

	_, err = store.TransferTx(ctx, arg)
	c.Error(err)
}

func BenchmarkTransferTx(b *testing.B) {
	InitCockroach()

	store := NewStore()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	for i := 0; i < b.N; i++ {
		domains, err := GetDomains(ctx)
		if err != nil {
			b.Fatal(err)
		}

		arg := TransferTxParams{
			FromDomainID: domains[0].DomainID,
		}

		_, err = store.TransferTx(ctx, arg)
		if err != nil {
			b.Fatal(err)
		}
	}
}
