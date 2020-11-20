package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewServerinSuccess(t *testing.T) {
	c := require.New(t)

	server, err := NewServer("server1", "B", "US", "Amazon.com, Inc.")
	c.NoError(err)
	c.NotEmpty(server.ServerID)
	c.NotNil(server.CreationDate)
	c.NotNil(server.UpdateDate)
}

func TestNewServerWithWrongParams(t *testing.T) {
	c := require.New(t)

	_, err := NewServer("", "B", "US", "Amazon.com, Inc.")
	c.EqualError(ErrEmptyAddress, err.Error())

	_, err = NewServer("server1", "", "US", "Amazon.com, Inc.")
	c.EqualError(ErrEmptySSLGrade, err.Error())

	_, err = NewServer("server1", "B", "", "Amazon.com, Inc.")
	c.EqualError(ErrEmptyCountry, err.Error())

	_, err = NewServer("server1", "B", "US", "")
	c.EqualError(ErrEmptyOwner, err.Error())
}

func BenchmarkNewServer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := NewServer("server1", "B", "US", "Amazon.com, Inc.")
		if err != nil {
			b.Fatal(err)
		}
	}
}
