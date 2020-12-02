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

		result.ToDomain, err = q.UpdateDomain(ctx, result.FromDomain.Servers[0].SSLGrade, "", arg.FromDomain, false)
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

		result.ToDomain, err = q.UpdateDomain(ctx, "", lastRecord.SSLGrade, arg.FromDomain, false)
		if err != nil {
			logs.Log().Errorf(`error updatedomain %s`, err.Error())
			return err
		}

		return nil
	})

	return result, err
}

// TransferTxParamsServerChange contains the input parameters of the transfer transaction
type TransferTxParamsServerChange struct {
	FromDomain *models.Domain `json:"from_domain"`
}

// TransferTxResultServerChange is the result of the transfer transaction
type TransferTxResultServerChange struct {
	FromDomain   *models.Domain   `json:"from_domain"`
	ToDomain     *models.Domain   `json:"to_domain"`
	ConsultTable []*models.Domain `json:"consult_table"`
}

// TransferTxServerChange performs a update ssl_grade attribute, transfer from ssl_grade server to ssl_grade domain
func (store *Store) TransferTxServerChange(ctx context.Context, arg TransferTxParamsServerChange) (TransferTxResultServerChange, error) {
	if arg.FromDomain == nil {
		return TransferTxResultServerChange{}, ErrEmptyDomainID
	}

	var result TransferTxResultServerChange

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

		serverChanged := false

		if result.FromDomain != lastRecord {

			count := 0
			for i := 0; i < len(lastRecord.Servers); i++ {
				if result.FromDomain.Servers[i].Address != lastRecord.Servers[i].Address {
					//fmt.Println(*result.FromDomain.Servers[i], "--- \n ----", *lastRecord.Servers[i])
					break
				}

				if result.FromDomain.Servers[i].Country != lastRecord.Servers[i].Country {
					//fmt.Println(*result.FromDomain.Servers[i], "--- \n ----", *lastRecord.Servers[i])
					break
				}

				if result.FromDomain.Servers[i].Owner != lastRecord.Servers[i].Owner {
					//fmt.Println(*result.FromDomain.Servers[i], "--- \n ----", *lastRecord.Servers[i])
					break
				}

				if result.FromDomain.Servers[i].SSLGrade != lastRecord.Servers[i].SSLGrade {
					//fmt.Println(*result.FromDomain.Servers[i], "--- \n ----", *lastRecord.Servers[i])
					break
				}
				count++

				//fmt.Println()
				//fmt.Println(result.FromDomain.Servers[i], "--- \n ----", lastRecord.Servers[i])
			}

			if count != len(lastRecord.Servers) {
				serverChanged = true
				/*
					fmt.Println(count, len(lastRecord.Servers))
					fmt.Println("serverchangeed --> ", serverChanged)
					fmt.Println()
					fmt.Println("result --> ", result.FromDomain, "\n last record -->", lastRecord)
					fmt.Println()
					fmt.Println(result.FromDomain.Servers[0], "\n", lastRecord.Servers[0])
				*/
			}
		}

		result.ToDomain, err = q.UpdateDomain(ctx, "", "", arg.FromDomain, serverChanged)
		if err != nil {
			logs.Log().Errorf(`error updatedomain %s`, err.Error())
			return err
		}

		return nil
	})

	return result, err
}
