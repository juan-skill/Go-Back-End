package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/other_project/crockroach/internal/logs"
	"github.com/other_project/crockroach/models"
)

const (
	createDomain = `
	INSERT INTO domains (
		id,
		domain_name,
		serverChanged,
		sslgrade,
		previousslgrade,
		logo,
		title,
		isdown,
		creationDate,
		updateDate
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
	) RETURNING *;
	`

	listDomains = `
	SELECT * FROM domains
	ORDER BY id
	LIMIT $1
	OFFSET $2
	`

	listDomainsByDate = `
	SELECT id, domain_name, serverchanged, sslgrade, previousslgrade, logo, title, isdown, creationdate, updatedate 
	FROM domains
	WHERE domains.updatedate >= now() - '1 hours'::INTERVAL 
    AND domains.updatedate <= now()
	`

	getDomain = `
	SELECT * FROM domains
	WHERE id = $1
	ORDER BY sslgrade DESC
	LIMIT 1
	`

	updateDomain = `
	UPDATE domains
	SET sslgrade = $2, updatedate = now()
	WHERE id = $1
	RETURNING *
	`

	updateDomainPrevioSSL = `
	UPDATE domains
	SET previousslgrade = $2, updatedate = now()
	WHERE id = $1
	RETURNING *
	`

	deleteDomain = `
	DELETE FROM domains
	WHERE id = $1
	`
)

var (
	// ErrInvalidDomain to ensure if exists domain
	ErrInvalidDomain = errors.New("invalid domain object")
	// ErrEmptyDomainID when it's
	ErrEmptyDomainID = errors.New("cannot be empty domain_id")
	// ErrDomainNotFound to ensure that domain are returned
	ErrDomainNotFound = errors.New("domain was not found")
	// ErrEmptyDomain when it does not exist a domain
	ErrEmptyDomain = errors.New("domain name cannot be empty")
)

// StoreDomain function will store a domain struct
func (q *Queries) StoreDomain(ctx context.Context, domain *models.Domain) (*models.Domain, error) {
	if domain == nil {
		logs.Log().Errorf("cannot store domain in database %s ", ErrInvalidDomain.Error())
		return nil, ErrInvalidDomain
	}

	row := CockroachClient.QueryRowContext(ctx, createDomain, domain.DomainID, domain.DomainName, domain.ServerChanged, domain.SSLGrade, domain.PreviousSSLGrade, domain.Logo, domain.Title, domain.IsDown, domain.CreationDate, domain.UpdateDate)
	if row.Err() != nil {
		logs.Log().Errorf("Query error %s", row.Err())
		return nil, ErrInvalidQuery
	}

	item := new(models.Domain)

	err := row.Scan(
		&item.DomainID,
		&item.DomainName,
		&item.ServerChanged,
		&item.SSLGrade,
		&item.PreviousSSLGrade,
		&item.Logo,
		&item.Title,
		&item.IsDown,
		&item.CreationDate,
		&item.UpdateDate,
	)
	if err != nil {
		logs.Log().Errorf("Scan error %s", err.Error())
		return nil, ErrScanRow
	}

	return item, nil
}

// GetDomain function will get a domain struct by domainID
func (q *Queries) GetDomain(ctx context.Context, domainID string) (*models.Domain, error) {
	if domainID == "" {
		logs.Log().Errorf("cannot store domain in database %s ", ErrEmptyDomainID.Error())
		return nil, ErrEmptyDomainID
	}

	row := CockroachClient.QueryRowContext(ctx, getDomain, domainID)
	if row.Err() != nil {
		logs.Log().Errorf("Query error %s", row.Err())
		return nil, ErrInvalidQuery
	}

	item := new(models.Domain)

	err := row.Scan(
		&item.DomainID,
		&item.DomainName,
		&item.ServerChanged,
		&item.SSLGrade,
		&item.PreviousSSLGrade,
		&item.Logo,
		&item.Title,
		&item.IsDown,
		&item.CreationDate,
		&item.UpdateDate,
	)
	if err != nil {
		logs.Log().Errorf("Scan error %s", err.Error())
		return nil, ErrScanRow
	}

	return item, nil
}

// GetDomains function will get a list of domains
func (q *Queries) GetDomains(ctx context.Context, time string) ([]models.Domain, error) {
	var rows *sql.Rows
	var err error

	if time == "" {
		rows, err = CockroachClient.QueryContext(ctx, listDomains, Limit, Offset)
	} else {
		rows, err = CockroachClient.QueryContext(ctx, listDomainsByDate)
	}

	if err != nil {
		logs.Log().Errorf("Query error %s", err.Error())
		return nil, ErrInvalidQuery
	}

	items := []models.Domain{}

	for rows.Next() {
		item := new(models.Domain)
		err := rows.Scan(
			&item.DomainID,
			&item.DomainName,
			&item.ServerChanged,
			&item.SSLGrade,
			&item.PreviousSSLGrade,
			&item.Logo,
			&item.Title,
			&item.IsDown,
			&item.CreationDate,
			&item.UpdateDate,
		)

		if err != nil {
			logs.Log().Errorf("Scan error %s", err.Error())
			return nil, err
		}

		item.Servers, err = q.GetServers(ctx, item.DomainID)
		if err != nil {
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

// UpdateDomain function will update a domain struct
func (q *Queries) UpdateDomain(ctx context.Context, sslgrade, previouSSL string, domain *models.Domain) (*models.Domain, error) {
	var row *sql.Row

	if domain == nil {
		logs.Log().Errorf("cannot be empty domain_id attribute %s ", ErrEmptyDomain)
		return nil, ErrEmptyDomain
	}

	if sslgrade != "" && previouSSL == "" {
		row = CockroachClient.QueryRowContext(ctx, updateDomain, domain.DomainID, sslgrade)
	} else if sslgrade == "" && previouSSL != "" {
		row = CockroachClient.QueryRowContext(ctx, updateDomainPrevioSSL, domain.DomainID, previouSSL)
	}

	if row.Err() != nil {
		logs.Log().Errorf("Query error %s", row.Err())
		return nil, ErrInvalidQuery
	}

	item := *domain

	err := row.Scan(
		&item.DomainID,
		&item.DomainName,
		&item.ServerChanged,
		&item.SSLGrade,
		&item.PreviousSSLGrade,
		&item.Logo,
		&item.Title,
		&item.IsDown,
		&item.CreationDate,
		&item.UpdateDate,
	)
	if err != nil {
		logs.Log().Errorf("Scan error %s", err.Error())
		return nil, ErrScanRow
	}

	return &item, err
}

// DeleteDomain function will update a domain struct
func (q *Queries) DeleteDomain(ctx context.Context, domainID string) error {
	if domainID == "" {
		logs.Log().Errorf("cannot be empty domain_id attribute %s ", ErrEmptyDomainID.Error())
		return ErrEmptyDomainID
	}

	row, err := CockroachClient.ExecContext(ctx, deleteDomain, domainID)
	if err != nil {
		logs.Log().Errorf("Query error %s", err.Error())
		return ErrInvalidQuery
	}

	result, _ := row.RowsAffected()
	if result == 0 {
		logs.Log().Errorf("Query error %s", ErrZeroRowsAffected.Error())
		return ErrZeroRowsAffected
	}

	return nil
}
