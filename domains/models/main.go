package models

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/other_project/crockroach/servers/models"
)

var (
	// ErrEmptyServers for empty servers
	ErrEmptyServers = errors.New("servers cannot be empty")
	// ErrEmptySSLGrade for empty SSLGrade
	ErrEmptySSLGrade = errors.New("ssl grade cannot be empty")
	// ErrEmptyPSSLGrade for empty psslgrade
	ErrEmptyPSSLGrade = errors.New("previous SSL grade cannot be empty")
	// ErrEmptyLogo for empty logo
	ErrEmptyLogo = errors.New("logo cannot be empty")
	// ErrEmptyTitle for empty title
	ErrEmptyTitle = errors.New("title cannot be empty")
)

// Domain model structure for domain
type Domain struct {
	DomainID         string           `json:"domain_id"`
	Servers          []*models.Server `json:"servers"`
	ServerChanged    bool             `json:"servers_changed"`
	SSLGrade         string           `json:"ssl_grade"`
	PreviousSSLGrade string           `json:"previous_ssl_grade"`
	Logo             string           `json:"logo"`
	Title            string           `json:"title"`
	IsDown           bool             `json:"is_down"`
	CreationDate     *time.Time       `json:"creation_date"`
	UpdateDate       *time.Time       `json:"update_date"`
}

// NewDomain Initialize a new domain
func NewDomain(servers []*models.Server, serverChanged, isdown bool, sslGrade, pSSLGrade, logo, title string) (domain *Domain, err error) {
	if servers == nil {
		return nil, ErrEmptyServers
	}

	if sslGrade == "" {
		return nil, ErrEmptySSLGrade
	}

	if pSSLGrade == "" {
		return nil, ErrEmptyPSSLGrade
	}

	if logo == "" {
		return nil, ErrEmptyLogo
	}

	if title == "" {
		return nil, ErrEmptyTitle
	}

	domainID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	created := time.Now()
	updated := time.Now()

	domain = &Domain{
		DomainID:         domainID.String(),
		Servers:          servers,
		ServerChanged:    serverChanged,
		SSLGrade:         sslGrade,
		PreviousSSLGrade: pSSLGrade,
		Logo:             logo,
		Title:            title,
		CreationDate:     &created,
		UpdateDate:       &updated,
	}

	return domain, nil
}
