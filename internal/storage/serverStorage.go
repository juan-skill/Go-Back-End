package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/other_project/crockroach/internal/logs"
	"github.com/other_project/crockroach/models"
	"github.com/other_project/crockroach/shared/env"
)

const (
	createServer = `
	INSERT INTO servers (
		id,
		address,
		sslgrade,
		country,
		owner,
		domain_id,
		creationDate,
		updateDate 
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8
	) RETURNING *;
	`

	getServer = `
	SELECT servers.id, servers.address, servers.sslgrade, servers.country, servers.owner, servers.creationdate, servers.updatedate, domains.id, domains.serverchanged, domains.sslgrade, domains.previousslgrade, domains.logo, domains.title, domains.isdown, domains.creationdate, domains.updatedate 
	FROM servers 
	INNER JOIN domains ON domains.id = servers.domain_id
	WHERE servers.id = $1 LIMIT 1
	`

	listServers = `
	SELECT servers.id, servers.address, servers.sslgrade, servers.country, servers.owner, servers.creationdate, servers.updatedate, domains.id, domains.serverchanged, domains.sslgrade, domains.previousslgrade, domains.logo, domains.title, domains.isdown, domains.creationdate, domains.updatedate 
	FROM servers 
	INNER JOIN domains ON domains.id = servers.domain_id
	ORDER BY servers.id
	LIMIT $1
	OFFSET $2
	`
	updateServer = `
	UPDATE servers
	SET sslgrade = $2
	WHERE id = $1
	RETURNING *
	`

	deleteServer = `
	DELETE FROM servers
	WHERE id = $1
	`
)

var (
	// CockroachClient creates a connection with the CockroachDB
	CockroachClient = sql.DB{}
	// ErrInvalidServer to ensure if exists server
	ErrInvalidServer = errors.New("invalid server object")
	// ErrEmptyServerID in
	ErrEmptyServerID = errors.New("cannot be empty server_id")
	// ErrInvalidQuery when the query is launch
	ErrInvalidQuery = errors.New("cannot query the database")
	// ErrEmptySSLGrade to ensure if exist ssl_grade
	ErrEmptySSLGrade = errors.New("cannot be empty ssl_grade ")
	// ErrScanRow when the row  is Scan copies the columns from the matched row
	ErrScanRow = errors.New("cannot scan query result of set")
	// ErrServerNotFound to ensure that server are returned
	ErrServerNotFound = errors.New("server was not found")
	// ErrZeroRowsAffected when try to delete a record does not exists
	ErrZeroRowsAffected = errors.New("cannot record that does not exist")
	// ErrEmptyList there are not element
	ErrEmptyList = errors.New("there are not elements")
	// Limit fdfddfgdf
	Limit = env.GetInt64("LIMIT_QUERY", 5)
	// Offset fdf
	Offset = env.GetInt64("OFFSET_QUERY", 5)
)

// StoreServer function will store a server struct
func (q *Queries) StoreServer(ctx context.Context, server *models.Server) (*models.Server, error) {
	if server == nil {
		logs.Log().Errorf("cannot store server in database %s ", ErrInvalidServer.Error())
		return nil, ErrInvalidServer
	}

	row := CockroachClient.QueryRowContext(ctx, createServer, server.ServerID, server.Address, server.SSLGrade, server.Country, server.Owner, server.Domain.DomainID, server.CreationDate, server.UpdateDate)
	if row.Err() != nil {
		logs.Log().Errorf("Query error %s", row.Err())
		return nil, ErrInvalidQuery
	}

	item := new(models.Server)
	item.Domain = new(models.Domain)

	err := row.Scan(
		&item.ServerID,
		&item.Address,
		&item.SSLGrade,
		&item.Country,
		&item.Owner,
		&item.Domain.DomainID,
		&item.CreationDate,
		&item.UpdateDate)
	if err != nil {
		logs.Log().Errorf("Scan error %s", err.Error())
		return nil, ErrScanRow
	}

	if item.Domain.DomainID == server.Domain.DomainID {
		item.Domain = server.Domain
	}

	return item, nil
}

// GetServer function will get a server struct by ServerID
func (q *Queries) GetServer(ctx context.Context, serverID string) (*models.Server, error) {
	if serverID == "" {
		logs.Log().Errorf("cannot be empty server_id %s ", ErrEmptyServerID.Error())
		return nil, ErrEmptyServerID
	}

	row := CockroachClient.QueryRowContext(ctx, getServer, serverID)
	if row.Err() != nil {
		logs.Log().Errorf("Query error %s", row.Err())
		return nil, ErrInvalidQuery
	}

	item := new(models.Server)
	item.Domain = new(models.Domain)

	err := row.Scan(
		&item.ServerID,
		&item.Address,
		&item.SSLGrade,
		&item.Country,
		&item.Owner,
		&item.CreationDate,
		&item.UpdateDate,
		&item.Domain.DomainID,
		&item.Domain.ServerChanged,
		&item.Domain.SSLGrade,
		&item.Domain.PreviousSSLGrade,
		&item.Domain.Logo,
		&item.Domain.Title,
		&item.Domain.IsDown,
		&item.Domain.CreationDate,
		&item.Domain.UpdateDate)
	if err != nil {
		logs.Log().Errorf("Scan error %s", err.Error())
		return nil, ErrScanRow
	}

	return item, nil
}

// GetServers function will get a list of servers
func (q *Queries) GetServers(ctx context.Context) ([]models.Server, error) {
	rows, err := CockroachClient.QueryContext(ctx, listServers, Limit, Offset)
	if err != nil {
		logs.Log().Errorf("Query error %s", err.Error())
		return nil, ErrInvalidQuery
	}

	items := []models.Server{}

	for rows.Next() {
		item := new(models.Server)

		item.Domain = new(models.Domain)
		if err = rows.Scan(
			&item.ServerID,
			&item.Address,
			&item.SSLGrade,
			&item.Country,
			&item.Owner,
			&item.CreationDate,
			&item.UpdateDate,
			&item.Domain.DomainID,
			&item.Domain.ServerChanged,
			&item.Domain.SSLGrade,
			&item.Domain.PreviousSSLGrade,
			&item.Domain.Logo,
			&item.Domain.Title,
			&item.Domain.IsDown,
			&item.Domain.CreationDate,
			&item.Domain.UpdateDate,
		); err != nil {
			return nil, err
		}

		items = append(items, *item)
	}

	if err := rows.Close(); err != nil {
		logs.Log().Errorf("Row error close %s", err.Error())
		return nil, err
	}

	if err := rows.Err(); err != nil {
		logs.Log().Errorf("Row error %s", err.Error())
		return nil, err
	}

	return items, nil
}

// UpdateServer function will update a server struct
func (q *Queries) UpdateServer(ctx context.Context, serverID, sslgrade string) (*models.Server, error) {
	if serverID == "" {
		logs.Log().Errorf("cannot be empty server_id attribute %s ", ErrEmptyServerID.Error())
		return nil, ErrEmptyServerID
	}

	if sslgrade == "" {
		logs.Log().Errorf("cannot be empty sslgrade attribute %s ", ErrEmptySSLGrade.Error())
		return nil, ErrEmptySSLGrade
	}

	row := CockroachClient.QueryRowContext(ctx, updateServer, serverID, sslgrade)
	if row.Err() != nil {
		logs.Log().Errorf("Query error %s", row.Err())
		return nil, ErrInvalidQuery
	}

	item := new(models.Server)
	item.Domain = new(models.Domain)

	err := row.Scan(
		&item.ServerID,
		&item.Address,
		&item.SSLGrade,
		&item.Country,
		&item.Owner,
		&item.Domain.DomainID,
		&item.CreationDate,
		&item.UpdateDate)
	if err != nil {
		logs.Log().Errorf("Scan error %s", err.Error())
		return nil, ErrScanRow
	}

	item.Domain, err = q.GetDomain(ctx, item.Domain.DomainID)
	if err != nil {
		return nil, err
	}

	return item, err
}

// DeleteServer function will update a server struct
func (q *Queries) DeleteServer(ctx context.Context, serverID string) error {
	if serverID == "" {
		logs.Log().Errorf("cannot be empty server_id attribute %s ", ErrEmptyServerID.Error())
		return ErrEmptyServerID
	}

	row, err := CockroachClient.ExecContext(ctx, deleteServer, serverID)
	if err != nil {
		logs.Log().Errorf("Query error %s", err.Error())
		return ErrInvalidQuery
	}

	result, _ := row.RowsAffected()
	if result == 0 {
		return ErrZeroRowsAffected
	}

	return nil
}
