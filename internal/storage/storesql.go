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
		rbErr := tx.Rollback()
		if rbErr != nil {
			logs.Log().Errorf("tx err: %v, rb err: %v", err, rbErr)
			return err
		}

		return err
	}

	return tx.Commit()
}

// TransferTxParamsServers contains the input parameters of the transfer transaction
type TransferTxParamsServers struct {
	FromDomain *models.Domain `json:"from_domain"`
}

// TransferTxResultServers is the result of the transfer transaction
type TransferTxResultServers struct {
	FromDomain *models.Domain `json:"from_domain"`
}

// TransferTxServers performs a update servers attribute
func (store *Store) TransferTxServers(ctx context.Context, arg TransferTxParamsServers) (TransferTxResultServers, error) {
	if arg.FromDomain == nil {
		return TransferTxResultServers{}, ErrEmptyDomainID
	}

	var result TransferTxResultServers

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.FromDomain, err = q.StoreDomain(ctx, arg.FromDomain)
		if err != nil {
			logs.Log().Errorf(`error Store Domain %s`, err.Error())
			return err
		}

		serversNumber := len(result.FromDomain.Servers)

		if len(result.FromDomain.Servers) == 0 {
			logs.Log().Errorf(`error Not Found servers of the domain`)
			return err
		}

		for i := 0; i < serversNumber; i++ {
			server := result.FromDomain.Servers[i]

			_, erro := q.StoreServer(ctx, server, result.FromDomain)
			if erro != nil {
				logs.Log().Errorf(`error in create a server of the domain: %s`, erro.Error())
				return erro
			}
		}

		return nil
	})

	return result, err
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

// TransferTxParamsInitialize contains the input parameters of the transfer transaction
type TransferTxParamsInitialize struct {
	FromDomain *models.Domain `json:"from_domain"`
}

// TransferTxResultInitialize is the result of the transfer transaction
type TransferTxResultInitialize struct {
	FromServers  []*models.Server `json:"from_servers"`
	ConsultTable []*models.Domain `json:"from_domains"`
	ToDomain     *models.Domain   `json:"to_domain"`
}

// TransferTxInitialize performs a update ssl_grade attribute, transfer from ssl_grade server to ssl_grade domain
func (store *Store) TransferTxInitialize(ctx context.Context, arg TransferTxParamsInitialize) (TransferTxResultInitialize, error) {
	if arg.FromDomain == nil {
		return TransferTxResultInitialize{}, ErrEmptyDomainID
	}

	var result TransferTxResultInitialize

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.FromServers, err = q.GetServers(ctx, arg.FromDomain.DomainID)
		if err != nil {
			logs.Log().Errorf(`error getservers by domain %s`, err.Error())
			return err
		}

		if result.FromServers == nil {
			return ErrEmptyServerByDomain
		}

		// servidores ordenados para tomar el primer servidor ssl_grade
		arg.FromDomain.Servers = result.FromServers

		// consultamos la tabla para saber los últimos registros
		result.ConsultTable, err = q.GetRecordByName(arg.FromDomain)
		if err != nil {
			return err
		}

		var sslGrade string
		var serverChanged bool
		if len(result.ConsultTable) == 0 {
			sslGrade = ""
			serverChanged = false
		} else {
			// el último registro de hace una hora para comparar
			lastRecord := result.ConsultTable[len(result.ConsultTable)-1]
			sslGrade = lastRecord.SSLGrade

			serverChanged = compareTwoDomains(arg.FromDomain, lastRecord)
		}

		result.ToDomain, err = q.UpdateDomain(ctx, arg.FromDomain.Servers[0].SSLGrade, sslGrade, arg.FromDomain, serverChanged)
		if err != nil {
			logs.Log().Errorf(`error updatedomain %s`, err.Error())
			return err
		}

		return nil
	})

	return result, err
}

// compareTwoDomains compare the current domain between a domain 1 hour ago
func compareTwoDomains(current, lastRecord *models.Domain) bool {
	serverChanged := false

	if current != lastRecord {
		count := 0

		for i := 0; i < len(lastRecord.Servers); i++ {
			if current.Servers[i].Address != lastRecord.Servers[i].Address {
				break
			}

			if current.Servers[i].Country != lastRecord.Servers[i].Country {
				break
			}

			if current.Servers[i].Owner != lastRecord.Servers[i].Owner {
				break
			}

			if current.Servers[i].SSLGrade != lastRecord.Servers[i].SSLGrade {
				break
			}
			count++
		}

		if count != len(lastRecord.Servers) {
			serverChanged = true
		}
	}

	return serverChanged
}
