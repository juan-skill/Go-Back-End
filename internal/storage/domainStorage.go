package storage

import (
	"context"
	"errors"
	"time"

	"github.com/other_project/crockroach/internal/logs"
	"github.com/other_project/crockroach/models"
)

const (
	createDomain = `
	INSERT INTO domains (
		id,
		serverChanged,
		sslgrade,
		previousslgrade,
		logo,
		title,
		isdown,
		creationDate,
		updateDate
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9
	) RETURNING *;
	`
	listDomains = `
	SELECT * FROM domains
	ORDER BY id
	LIMIT $1
	OFFSET $2
	`

	getDomain = `
	SELECT * FROM domains
	WHERE id = $1 LIMIT 1
	`
	updateDomain = `
	UPDATE domains
	SET sslgrade = $2
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
)

// StoreDomain function will store a domain struct
func (q *Queries) StoreDomain(domain *models.Domain) (*models.Domain, error) {
	if domain == nil {
		logs.Log().Errorf("cannot store domain in database %s ", ErrInvalidDomain.Error())
		return nil, ErrInvalidDomain
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	row := CockroachClient.QueryRowContext(ctx, createDomain, domain.DomainID, domain.ServerChanged, domain.SSLGrade, domain.PreviousSSLGrade, domain.Logo, domain.Title, domain.IsDown, domain.CreationDate, domain.UpdateDate)
	if row.Err() != nil {
		logs.Log().Errorf("Query error %s", row.Err())
		return nil, ErrInvalidQuery
	}

	item := new(models.Domain)

	err := row.Scan(
		&item.DomainID,
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
func (q *Queries) GetDomain(domainID string) (*models.Domain, error) {
	if domainID == "" {
		logs.Log().Errorf("cannot store domain in database %s ", ErrEmptyDomainID.Error())
		return nil, ErrEmptyDomainID
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	row := CockroachClient.QueryRowContext(ctx, getDomain, domainID)
	if row.Err() != nil {
		logs.Log().Errorf("Query error %s", row.Err())
		return nil, ErrInvalidQuery
	}

	item := new(models.Domain)

	err := row.Scan(
		&item.DomainID,
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
func (q *Queries) GetDomains() ([]models.Domain, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	rows, err := CockroachClient.QueryContext(ctx, listDomains, Limit, Offset)
	if err != nil {
		logs.Log().Errorf("Query error %s", err.Error())
		return nil, ErrInvalidQuery
	}

	items := []models.Domain{}

	for rows.Next() {
		item := new(models.Domain)
		if err := rows.Scan(
			&item.DomainID,
			&item.ServerChanged,
			&item.SSLGrade,
			&item.PreviousSSLGrade,
			&item.Logo,
			&item.Title,
			&item.IsDown,
			&item.CreationDate,
			&item.UpdateDate,
		); err != nil {
			logs.Log().Errorf("Scan error %s", err.Error())
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
func (q *Queries) UpdateDomain(domainID, sslgrade string) (*models.Domain, error) {
	if domainID == "" {
		logs.Log().Errorf("cannot be empty domain_id attribute %s ", ErrEmptyDomainID.Error())
		return nil, ErrEmptyDomainID
	}

	if sslgrade == "" {
		logs.Log().Errorf("cannot be empty sslgrade attribute %s ", ErrEmptySSLGrade.Error())
		return nil, ErrEmptySSLGrade
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	row := CockroachClient.QueryRowContext(ctx, updateDomain, domainID, sslgrade)
	if row.Err() != nil {
		logs.Log().Errorf("Query error %s", row.Err())
		return nil, ErrInvalidQuery
	}

	item := new(models.Domain)

	err := row.Scan(
		&item.DomainID,
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

	return item, err
}

// DeleteDomain function will update a domain struct
func (q *Queries) DeleteDomain(domainID string) error {
	if domainID == "" {
		logs.Log().Errorf("cannot be empty domain_id attribute %s ", ErrEmptyDomainID.Error())
		return ErrEmptyDomainID
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

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
