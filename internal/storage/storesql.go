package storage

import (
	"context"
	"errors"

	"github.com/other_project/crockroach/internal/logs"
	"github.com/other_project/crockroach/models"
)

var (
	// ErrEmptyServerByDomain when it trying to fill ssl_grade field
	ErrEmptyServerByDomain = errors.New("there are not servers in this domain")
)

// Store provides all functions to execute SQL queries and transactions
type Store struct {
	*Queries
}

// NewStore creates a new store
func NewStore() *Store {
	return &Store{
		Queries: NewQueries(),
	}
}

// ExecTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := CockroachClient.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := NewQueries()

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			logs.Log().Errorf("tx err: %v, rb err: %v", err, rbErr)
			return err
		}

		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromDomainID string `json:"from_domain_id"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	FromDomain *models.Domain `json:"from_domain"`
	ToDomain   *models.Domain `json:"to_domain"`
}

// TransferTx performs a update ssl_grade attribute, transfer from ssl_grade server to ssl_grade domain
// It get the domain, get the servers of the domain, and update ssl_grade attribute of domain within a database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	if arg.FromDomainID == "" {
		return TransferTxResult{}, ErrEmptyDomainID
	}

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.FromDomain, err = q.GetDomain(ctx, arg.FromDomainID)
		if err != nil {
			logs.Log().Errorf(`error getdomain %s`, err.Error())
			return err
		}

		result.FromDomain.Servers, err = q.GetServers(ctx, arg.FromDomainID)
		if err != nil {
			logs.Log().Errorf(`error getservers by domain %s`, err.Error())
			return err
		}

		if result.FromDomain.Servers == nil {
			return ErrEmptyServerByDomain
		}

		result.ToDomain, err = q.UpdateDomain(ctx, arg.FromDomainID, result.FromDomain.Servers[0].SSLGrade)
		if err != nil {
			logs.Log().Errorf(`error updatedomain %s`, err.Error())
			return err
		}

		return nil
	})

	return result, err
}
