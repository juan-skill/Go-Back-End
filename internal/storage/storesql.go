package storage

import (
	"context"
	"errors"
	"fmt"

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
	FromDomain *models.Domain `json:"from_domain"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	FromDomain *models.Domain `json:"from_domain"`
	ToDomain   *models.Domain `json:"to_domain"`
}

// TransferTx performs a update ssl_grade attribute, transfer from ssl_grade server to ssl_grade domain
// It get the domain, get the servers of the domain, and update ssl_grade attribute of domain within a database transaction
// el grado ssl lo tiene el primer servidor por eso lo seleccionamos por defecto
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	if arg.FromDomain == nil {
		return TransferTxResult{}, ErrEmptyDomainID
	}

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.FromDomain, err = q.GetDomain(ctx, arg.FromDomain.DomainID)
		if err != nil {
			logs.Log().Errorf(`error getdomain %s`, err.Error())
			return err
		}

		result.FromDomain.Servers, err = q.GetServers(ctx, arg.FromDomain.DomainID)
		if err != nil {
			logs.Log().Errorf(`error getservers by domain %s`, err.Error())
			return err
		}

		if result.FromDomain.Servers == nil {
			return ErrEmptyServerByDomain
		}

		result.ToDomain, err = q.UpdateDomain(ctx, result.FromDomain.Servers[0].SSLGrade, "", arg.FromDomain)
		if err != nil {
			logs.Log().Errorf(`error updatedomain %s`, err.Error())
			return err
		}

		return nil
	})

	return result, err
}

// TransferTxParamsPreSSL contains the input parameters of the transfer transaction
type TransferTxParamsPreSSL struct {
	FromDomain *models.Domain `json:"from_domain"`
}

// TransferTxResultPreSSL is the result of the transfer transaction
type TransferTxResultPreSSL struct {
	FromDomain   *models.Domain   `json:"from_domain"`
	ToDomain     *models.Domain   `json:"to_domain"`
	ConsultTable []*models.Domain `json:"consult_table"`
}

// TransferTxPreSSL performs a update ssl_grade attribute, transfer from ssl_grade server to ssl_grade domain
func (store *Store) TransferTxPreSSL(ctx context.Context, arg TransferTxParamsPreSSL) (TransferTxResultPreSSL, error) {
	if arg.FromDomain == nil {
		return TransferTxResultPreSSL{}, ErrEmptyDomainID
	}

	var result TransferTxResultPreSSL

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.FromDomain, err = q.GetDomain(ctx, arg.FromDomain.DomainID)
		if err != nil {
			logs.Log().Errorf(`error getdomain %s`, err.Error())
			return err
		}

		result.FromDomain.Servers, err = q.GetServers(ctx, arg.FromDomain.DomainID)
		if err != nil {
			logs.Log().Errorf(`error getservers by domain %s`, err.Error())
			return err
		}

		if result.FromDomain.Servers == nil {
			return ErrEmptyServerByDomain
		}

		result.ConsultTable, err = q.GetRecordByName(result.FromDomain)
		if err != nil {
			return err
		}

		if len(result.ConsultTable) == 0 {
			result.ToDomain = arg.FromDomain
			return nil
		}

		lastRecord := result.ConsultTable[len(result.ConsultTable)-1]
		fmt.Println("lastrecord --> ", lastRecord)
		previoSSL := lastRecord.SSLGrade
		fmt.Println("PREviossl lastrecord --> ", previoSSL)
		//fmt.Println("ssl grade last record ", lastRecord.SSLGrade)
		/*
			var domainFounded *models.Domain

			for key, value := range result.ConsultTable {
				if value.Domain.DomainID == result.FromDomain.DomainID {
					domainFounded = result.ConsultTable[key].Domain
					break
				}
			}

			if domainFounded == nil {
				return nil
			}
		*/

		/*
			countEquality := 0
			numberServers := 0
			if len(domainFounded.Servers) == len(result.FromDomain.Servers) {
				numberServers = len(domainFounded.Servers)

				for i := 0; i < numberServers; i++ {

					if domainFounded.Servers[i].SSLGrade != result.FromDomain.Servers[i].SSLGrade {
						countEquality++
					}
				}
			}

			if countEquality == numberServers {
				return err
			}
		*/

		result.ToDomain, err = q.UpdateDomain(ctx, "", previoSSL, arg.FromDomain)
		if err != nil {
			logs.Log().Errorf(`error updatedomain %s`, err.Error())
			return err
		}

		return nil
	})

	return result, err
}
