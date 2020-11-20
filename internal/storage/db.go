package storage

import (
	"database/sql"

	"github.com/other_project/crockroach/models"
)

var (
	// Default to make the interface PostStorage global
	Default DBTX
)

// DBTX interface
type DBTX interface {
	StoreServer(server *models.Server) (*models.Server, error)
	GetServer(serverID string) (*models.Server, error)
	UpdateServer(serverID, sslgrade string) (*models.Server, error)
	DeleteServer(serverID string) error
	GetServers() ([]models.Server, error)
	StoreDomain(domain *models.Domain) (*models.Domain, error)
	GetDomain(domainID string) (*models.Domain, error)
	UpdateDomain(serverID, sslgrade string) (*models.Domain, error)
	DeleteDomain(domainID string) error
	GetDomains() ([]models.Domain, error)
	/*
		ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
		PrepareContext(context.Context, string) (*sql.Stmt, error)
		QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
		QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	*/
}

// NewQueries function create a new instance of
/*
func NewQueries(db DBTX) *Queries {
	return &Queries{
		db: db,
	}
}
*/

// Queries structure allow us extend the functionality
type Queries struct {
	//db DBTX
}

// WithTx return a instance with transaction connection
/*
func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}

}
*/

// StoreServer function will store a server in the database.
func StoreServer(server *models.Server) (*models.Server, error) {
	return Default.StoreServer(server)
}

// GetServer function will retrieve a server in the database.
func GetServer(serverID string) (*models.Server, error) {
	return Default.GetServer(serverID)
}

// UpdateServer function will update a server struct
func UpdateServer(serverID, sslgrade string) (*models.Server, error) {
	return Default.UpdateServer(serverID, sslgrade)
}

// DeleteServer function will delete a server struct
func DeleteServer(serverID string) error {
	return Default.DeleteServer(serverID)
}

// GetServers function will list all the servers structures
func GetServers() ([]models.Server, error) {
	return Default.GetServers()
}

// StoreDomain function will store a domain in the database.
func StoreDomain(domain *models.Domain) (*models.Domain, error) {
	return Default.StoreDomain(domain)
}

// GetDomain function will retrieve a domain in the database.
func GetDomain(domainID string) (*models.Domain, error) {
	return Default.GetDomain(domainID)
}

// UpdateDomain function will update a domain struct
func UpdateDomain(serverID, sslgrade string) (*models.Domain, error) {
	return Default.UpdateDomain(serverID, sslgrade)
}

// DeleteDomain function will delete a domain struct
func DeleteDomain(domainID string) error {
	return Default.DeleteDomain(domainID)
}

// GetDomains function will list all the domains structures
func GetDomains() ([]models.Domain, error) {
	return Default.GetDomains()
}

func init() {
	Default = &Queries{}
	CockroachClient = sql.DB{}
}
