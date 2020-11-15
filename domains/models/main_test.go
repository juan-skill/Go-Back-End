package models

import (
	"testing"

	"github.com/other_project/crockroach/servers/models"
	"github.com/stretchr/testify/require"
)

// createListServer create a list of servers
func createListServer() (servers []*models.Server) {

	server1, _ := models.NewServer("server1", "B", "US", "Amazon.com, Inc.")
	server2, _ := models.NewServer("server2", "A+", "US", "Amazon.com, Inc.")
	server3, _ := models.NewServer("server3", "A", "US", "Amazon.com, Inc.")

	servers = make([]*models.Server, 3)

	return append(servers, server1, server2, server3)
}

func TestNewDomainSuccess(t *testing.T) {
	c := require.New(t)

	servers := createListServer()

	domain, err := NewDomain(servers, false, false, "A+", "B", "https://server.com/icon.png", "Title of the page")
	c.NoError(err)
	c.NotEmpty(domain.DomainID)
	c.NotNil(domain.CreationDate)
	c.NotNil(domain.UpdateDate)
}

func TestNewDomainWithWrongParams(t *testing.T) {
	c := require.New(t)

	servers := createListServer()

	_, err := NewDomain(nil, false, false, "A+", "B", "https://server.com/icon.png", "Title of the page")
	c.EqualError(ErrEmptyServers, err.Error())

	_, err = NewDomain(servers, false, false, "", "B", "https://server.com/icon.png", "Title of the page")
	c.EqualError(ErrEmptySSLGrade, err.Error())

	_, err = NewDomain(servers, false, false, "A+", "", "https://server.com/icon.png", "Title of the page")
	c.EqualError(ErrEmptyPSSLGrade, err.Error())

	_, err = NewDomain(servers, false, false, "A+", "B", "", "Title of the page")
	c.EqualError(ErrEmptyLogo, err.Error())

	_, err = NewDomain(servers, false, false, "A+", "B", "https://server.com/icon.png", "")
	c.EqualError(ErrEmptyTitle, err.Error())
}

func BenchmarkNewDomain(b *testing.B) {

	servers := createListServer()

	for n := 0; n < b.N; n++ {
		_, err := NewDomain(servers, false, false, "A+", "B", "https://server.com/icon.png", "Title of the page")
		if err != nil {
			b.Fatal(err)
		}
	}
}
