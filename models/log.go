package models

import (
	"errors"
	"time"
)

var (
	// ErrEmptyDomain when it does not exist a domain
	ErrEmptyDomain = errors.New("domain name cannot be empty")
	// Count when it creates a record
	Count int64
)

// LogDomainStatus model structure for server
type LogDomainStatus struct {
	LogDomainStatusID int64      `json:"log_id"`
	Domain            *Domain    `json:"domain_id"`
	DomainName        string     `json:"domain_name"`
	SSLGrade          string     `json:"ssl_grade"`
	ServerChanged     bool       `json:"server_changed"`
	UpdateDate        *time.Time `json:"update_date"`
}

// NewLogDomainStatus Initialize a new log
func NewLogDomainStatus(domainName, sslGrade string, domain *Domain) (logDomain *LogDomainStatus, err error) {
	if domainName == "" {
		return nil, ErrEmptyDomainName
	}

	if domain == nil {
		return nil, ErrEmptyDomain
	}

	updated := time.Now()
	Count = Count + 1

	logDomain = &LogDomainStatus{
		LogDomainStatusID: Count,
		Domain:            domain,
		DomainName:        domainName,
		SSLGrade:          sslGrade,
		ServerChanged:     false,
		UpdateDate:        &updated,
	}

	return logDomain, nil
}

func init() {
	Count = 0
}
