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

	//domain, err := models.NewDomain(false, false, "google.com", "A+", testrandom.RandomSSLRating("B"), "https://server.com/icon.png", "Title of the page")
	domain, err := models.NewDomain(false, false, "google.com", "", "", "https://server.com/icon.png", "Title of the page")
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
		FromDomain: domain1,
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
		FromDomain: nil,
	}

	_, err := store.TransferTx(ctx, arg)
	c.Error(err)
}

func TestTransferTxPreSSL(t *testing.T) {
	c := require.New(t)

	// Iniciar la base de datos
	InitCockroach()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	// iniciar las operaciones de transacciones
	store := NewStore()

	// cargar los ultimos registros hace una hora
	_, err := ReloadRecord(ctx)
	c.NoError(err)

	// crear un dominio que no está en la base de datos
	domain1 := getNewDomain(t)

	// guardamos una copia en la tabla cache
	record, err := NewRecord(domain1)
	c.NoError(err)
	c.NotEmpty(record)

	// reasingnar el attributo gradeSSL en la base de datos
	arg := TransferTxParams{
		FromDomain: domain1,
	}

	result, err := store.TransferTx(ctx, arg)
	c.NoError(err)
	c.Equal(domain1.SSLGrade, result.FromDomain.SSLGrade)

	// guardamos una copia en la tabla cache con el nuevo estado gradeSSL
	record, err = NewRecord(result.ToDomain)
	c.NoError(err)
	c.NotEmpty(record)

	// reasignar el attributo previoGradeSSL
	argPre := TransferTxParamsPreSSL{
		FromDomain: domain1,
	}

	result1, err := store.TransferTxPreSSL(ctx, argPre)
	c.NoError(err)
	c.Equal(domain1.PreviousSSLGrade, result1.FromDomain.PreviousSSLGrade)

	// guardar en la tabla cache este nuevo registro
	record, err = NewRecord(domain1)
	c.NoError(err)
	c.NotEmpty(record)
}

func TestTransferTxServerChanged(t *testing.T) {
	c := require.New(t)

	// Iniciar la base de datos
	InitCockroach()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	// iniciar las operaciones de transacciones
	store := NewStore()

	// cargar los ultimos registros hace una hora
	_, err := ReloadRecord(ctx)
	c.NoError(err)

	// crear un dominio que no está en la base de datos
	domain1 := getNewDomain(t)

	// guardamos una copia en la tabla cache
	record, err := NewRecord(domain1)
	c.NoError(err)
	c.NotEmpty(record)

	// reasingnar el attributo gradeSSL en la base de datos
	arg := TransferTxParams{
		FromDomain: domain1,
	}

	result, err := store.TransferTx(ctx, arg)
	c.NoError(err)
	c.Equal(domain1.SSLGrade, result.FromDomain.SSLGrade)

	// guardamos una copia en la tabla cache con el nuevo estado gradeSSL
	record, err = NewRecord(result.ToDomain)
	c.NoError(err)
	c.NotEmpty(record)

	// reasignar el attributo previoGradeSSL
	argPre := TransferTxParamsPreSSL{
		FromDomain: domain1,
	}

	result1, err := store.TransferTxPreSSL(ctx, argPre)
	c.NoError(err)
	c.Equal(domain1.PreviousSSLGrade, result1.FromDomain.PreviousSSLGrade)

	// guardar en la tabla cache este nuevo registro
	record, err = NewRecord(result1.ToDomain)
	c.NoError(err)
	c.NotEmpty(record)

	// reasignar el attributo serverchanged
	argServer := TransferTxParamsServerChange{
		FromDomain: domain1,
	}

	result2, err := store.TransferTxServerChange(ctx, argServer)
	c.NoError(err)
	c.NotEmpty(result2)
	//c.Equal(domain1.PreviousSSLGrade, result2.FromDomain.PreviousSSLGrade)

	// guardar en la tabla cache este nuevo registro
	record, err = NewRecord(result2.ToDomain)
	c.NoError(err)
	c.NotEmpty(record)
}

/*



 */

// TestDelete
func TestRomasnPere(t *testing.T) {
	c := require.New(t)

	InitCockroach()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	domains, err := GetDomains(ctx, "")
	c.NoError(err)

	for _, domain := range domains {
		err = DeleteDomain(ctx, domain.DomainID)
		c.Nil(err)
	}

	_, err = GetDomains(ctx, "")
	c.NoError(err)
}

func BenchmarkTransferTx(b *testing.B) {
	InitCockroach()

	store := NewStore()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	for i := 0; i < b.N; i++ {
		domains, err := GetDomains(ctx, "")
		if err != nil {
			b.Fatal(err)
		}

		arg := TransferTxParams{
			FromDomain: &domains[0],
		}

		_, err = store.TransferTx(ctx, arg)
		if err != nil {
			b.Fatal(err)
		}
	}
}
