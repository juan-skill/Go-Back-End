package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// createListServer create a list of servers
/*
func createListServer() (servers []*string) {

	server1, _ := NewServer("server1", "B", "US", "Amazon.com, Inc.")
	server2, _ := NewServer("server2", "A+", "US", "Amazon.com, Inc.")
	server3, _ := NewServer("server3", "A", "US", "Amazon.com, Inc.")

	servers = make([]*string, 3)

	return append(servers, &server1.ServerID, &server2.ServerID, &server3.ServerID)
}
*/

func TestNewDomainSuccess(t *testing.T) {
	c := require.New(t)

	//servers := createListServer()

	domain, err := NewDomain(false, false, "google.com", "A+", "B", "https://server.com/icon.png", "Title of the page")

	c.NoError(err)
	c.NotEmpty(domain.DomainName)
	c.NotEmpty(domain.DomainID)
	c.NotNil(domain.CreationDate)
	c.NotNil(domain.UpdateDate)
}

func TestNewDomainWithWrongParams(t *testing.T) {
	c := require.New(t)

	//servers := createListServer()

	_, err := NewDomain(false, false, "google.com", "", "B", "https://server.com/icon.png", "Title of the page")
	c.NoError(err)

	_, err = NewDomain(false, false, "google.com", "A+", "", "https://server.com/icon.png", "Title of the page")
	c.NoError(err)

	_, err = NewDomain(false, false, "google.com", "A+", "B", "", "Title of the page")
	c.EqualError(ErrEmptyLogo, err.Error())

	_, err = NewDomain(false, false, "google.com", "A+", "B", "https://server.com/icon.png", "")
	c.EqualError(ErrEmptyTitle, err.Error())

	_, err = NewDomain(false, false, "", "A+", "B", "https://server.com/icon.png", "Title of the page")
	c.EqualError(ErrEmptyDomainName, err.Error())
}

func BenchmarkNewDomain(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := NewDomain(false, false, "google.com", "A+", "B", "https://server.com/icon.png", "Title of the page")
		if err != nil {
			b.Fatal(err)
		}
	}
}
