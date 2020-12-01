package models

import (
	"testing"

	"github.com/other_project/crockroach/shared/testrandom"
	"github.com/stretchr/testify/require"
)

func storeServerTest(t *testing.T) *Domain {
	c := require.New(t)

	domain, err := NewDomain(false, false, "google.com", "A+", testrandom.RandomSSLRating("B"), "https://server.com/icon.png", "Title of the page")
	c.NoError(err)
	c.NotEmpty(domain)
	c.Equal("google.com", domain.DomainName)

	serverNumber := testrandom.RandomServerNumber()
	var i int64

	for i = 0; i < serverNumber; i++ {
		// create a new Server
		server, err := NewServer("server1", testrandom.RandomSSLRating(""), "US", "Amazon.com, Inc.", domain)
		c.NoError(err)
		c.NotNil(server)

		domain.Servers = append(domain.Servers, server)
	}

	return domain
}

func TestNewLogDomainSuccess(t *testing.T) {
	c := require.New(t)

	domain := storeServerTest(t)

	log, err := NewLogDomainStatus("google.com", "C", domain)

	c.NoError(err)
	c.NotEmpty(log.LogDomainStatusID)
	c.NotEmpty(log.Domain)
	c.Equal("google.com", domain.DomainName)
	c.NotNil(log.UpdateDate)
}

func TestNewLogDomainWithWrongParams(t *testing.T) {
	c := require.New(t)

	domain := storeServerTest(t)

	_, err := NewLogDomainStatus("", "C", domain)
	c.EqualError(ErrEmptyDomainName, err.Error())

	_, err = NewLogDomainStatus("google.com", "", domain)
	c.NoError(err)

	_, err = NewLogDomainStatus("google.com", "E", nil)
	c.EqualError(ErrEmptyDomain, err.Error())
}

func BenchmarkNewLogDomain(b *testing.B) {
	for n := 0; n < b.N; n++ {
		domain, err := NewDomain(false, false, "google.com", "A+", testrandom.RandomSSLRating("B"), "https://server.com/icon.png", "Title of the page")
		if err != nil {
			b.Fatal(err)
		}

		serverNumber := testrandom.RandomServerNumber()
		var i int64

		for i = 0; i < serverNumber; i++ {
			server, erro := NewServer("server1", testrandom.RandomSSLRating(""), "US", "Amazon.com, Inc.", domain)
			if erro != nil {
				b.Fatal(erro)
			}

			domain.Servers = append(domain.Servers, server)
		}

		_, err = NewLogDomainStatus("google.com", "E", domain)
		if err != nil {
			b.Fatal(err)
		}
	}
}
