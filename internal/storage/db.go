package storage

import (
	"context"
	"database/sql"

	"github.com/other_project/crockroach/models"
)

var (
	// Default to make the interface PostStorage global
	Default DBTX
)

// DBTX interface
type DBTX interface {
	StoreServer(ctx context.Context, server *models.Server) (*models.Server, error)
	GetServer(ctx context.Context, serverID string) (*models.Server, error)
	UpdateServer(ctx context.Context, serverID, sslgrade string) (*models.Server, error)
	DeleteServer(ctx context.Context, serverID string) error
	GetServers(ctx context.Context) ([]models.Server, error)
	StoreDomain(ctx context.Context, domain *models.Domain) (*models.Domain, error)
	GetDomain(ctx context.Context, domainID string) (*models.Domain, error)
	UpdateDomain(ctx context.Context, serverID, sslgrade string) (*models.Domain, error)
	DeleteDomain(ctx context.Context, domainID string) error
	GetDomains(ctx context.Context) ([]models.Domain, error)
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
func StoreServer(ctx context.Context, server *models.Server) (*models.Server, error) {
	return Default.StoreServer(ctx, server)
}

// GetServer function will retrieve a server in the database.
func GetServer(ctx context.Context, serverID string) (*models.Server, error) {
	return Default.GetServer(ctx, serverID)
}

// UpdateServer function will update a server struct
func UpdateServer(ctx context.Context, serverID, sslgrade string) (*models.Server, error) {
	return Default.UpdateServer(ctx, serverID, sslgrade)
}

// DeleteServer function will delete a server struct
func DeleteServer(ctx context.Context, serverID string) error {
	return Default.DeleteServer(ctx, serverID)
}

// GetServers function will list all the servers structures
func GetServers(ctx context.Context) ([]models.Server, error) {
	return Default.GetServers(ctx)
}

// StoreDomain function will store a domain in the database.
func StoreDomain(ctx context.Context, domain *models.Domain) (*models.Domain, error) {
	return Default.StoreDomain(ctx, domain)
}

// GetDomain function will retrieve a domain in the database.
func GetDomain(ctx context.Context, domainID string) (*models.Domain, error) {
	return Default.GetDomain(ctx, domainID)
}

// UpdateDomain function will update a domain struct
func UpdateDomain(ctx context.Context, serverID, sslgrade string) (*models.Domain, error) {
	return Default.UpdateDomain(ctx, serverID, sslgrade)
}

// DeleteDomain function will delete a domain struct
func DeleteDomain(ctx context.Context, domainID string) error {
	return Default.DeleteDomain(ctx, domainID)
}

// GetDomains function will list all the domains structures
func GetDomains(ctx context.Context) ([]models.Domain, error) {
	return Default.GetDomains(ctx)
}

func init() {
	Default = &Queries{}
	CockroachClient = sql.DB{}
}
