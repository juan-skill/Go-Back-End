package models

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

var (
	// ErrEmptyAddress for empty address
	ErrEmptyAddress = errors.New("address cannot be empty")
	// ErrEmptySSLGrade for empty SSLGrade
	ErrEmptySSLGrade = errors.New("ssl grade cannot be empty")
	// ErrEmptyCountry for empty country
	ErrEmptyCountry = errors.New("country cannot be empty")
	// ErrEmptyOwner for empty Owner
	ErrEmptyOwner = errors.New("owner cannot be empty")
)

// Server model structure for server
type Server struct {
	ServerID     string     `json:"server_id"`
	Address      string     `json:"address"`
	SSLGrade     string     `json:"ssl_grade"`
	Country      string     `json:"country"`
	Owner        string     `json:"owner"`
	CreationDate *time.Time `json:"creation_date"`
	UpdateDate   *time.Time `json:"update_date"`
}

// NewServer Initialize a new server
func NewServer(address, sslGrade, country, owner string) (server *Server, err error) {
	if address == "" {
		return nil, ErrEmptyAddress
	}

	if sslGrade == "" {
		return nil, ErrEmptySSLGrade
	}

	if country == "" {
		return nil, ErrEmptyCountry
	}

	if owner == "" {
		return nil, ErrEmptyOwner
	}

	serverID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	created := time.Now()
	updated := time.Now()

	server = &Server{
		ServerID:     serverID.String(),
		Address:      address,
		SSLGrade:     sslGrade,
		Country:      country,
		Owner:        owner,
		CreationDate: &created,
		UpdateDate:   &updated,
	}

	return server, nil
}
