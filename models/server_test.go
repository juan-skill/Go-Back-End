package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewServerinSuccess(t *testing.T) {
	c := require.New(t)

	domain, err := NewDomain(false, false, "A+", "B", "https://server.com/icon.png", "Title of the page")
	c.NoError(err)
	c.NotEmpty(domain)

	server, err := NewServer("server1", "B", "US", "Amazon.com, Inc.", domain)
	c.NoError(err)
	c.NotEmpty(server.ServerID)
	c.NotEmpty(domain.DomainID)
	c.NotNil(server.CreationDate)
	c.NotNil(server.UpdateDate)
}

func TestNewServerWithWrongParams(t *testing.T) {
	c := require.New(t)

	domain, err := NewDomain(false, false, "A+", "B", "https://server.com/icon.png", "Title of the page")
	c.NoError(err)
	c.NotEmpty(domain)

	_, err = NewServer("", "B", "US", "Amazon.com, Inc.", domain)
	c.EqualError(ErrEmptyAddress, err.Error())

	_, err = NewServer("server1", "", "US", "Amazon.com, Inc.", domain)
	c.EqualError(ErrEmptySSLGrade, err.Error())

	_, err = NewServer("server1", "B", "", "Amazon.com, Inc.", domain)
	c.EqualError(ErrEmptyCountry, err.Error())

	_, err = NewServer("server1", "B", "US", "", domain)
	c.EqualError(ErrEmptyOwner, err.Error())
}

func BenchmarkNewServer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		domain, err := NewDomain(false, false, "A+", "B", "https://server.com/icon.png", "Title of the page")
		if err != nil {
			b.Fatal(err)
		}

		_, err = NewServer("server1", "B", "US", "Amazon.com, Inc.", domain)
		if err != nil {
			b.Fatal(err)
		}
	}
}
