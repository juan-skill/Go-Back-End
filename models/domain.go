package models

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

var (
	// ErrEmptyDomainName for empty servers
	ErrEmptyDomainName = errors.New("Domain name cannot be empty")
	// ErrEmptyServers for empty servers
	ErrEmptyServers = errors.New("servers cannot be empty")
	// ErrEmptyPSSLGrade for empty psslgrade
	ErrEmptyPSSLGrade = errors.New("previous SSL grade cannot be empty")
	// ErrEmptyLogo for empty logo
	ErrEmptyLogo = errors.New("logo cannot be empty")
	// ErrEmptyTitle for empty title
	ErrEmptyTitle = errors.New("title cannot be empty")
)

// Domain model structure for domain
type Domain struct {
	DomainID         string     `json:"domain_id"`
	DomainName       string     `json:"domain_name"`
	Servers          []*Server  `json:"servers"`
	ServerChanged    bool       `json:"servers_changed"`
	SSLGrade         string     `json:"ssl_grade"`
	PreviousSSLGrade string     `json:"previous_ssl_grade"`
	Logo             string     `json:"logo"`
	Title            string     `json:"title"`
	IsDown           bool       `json:"is_down"`
	CreationDate     *time.Time `json:"creation_date"`
	UpdateDate       *time.Time `json:"update_date"`
}

// NewDomain Initialize a new domain
func NewDomain(serverChanged, isdown bool, domainName, sslGrade, pSSLGrade, logo, title string) (domain *Domain, err error) {
	if domainName == "" {
		return nil, ErrEmptyDomainName
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
		DomainName:       domainName,
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
