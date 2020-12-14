package testrandom

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	c := require.New(t)

	number, err := RandomInt(5, 10)
	c.NoError(err)
	c.NotEmpty(number)

	num := int64(0)

	number, err = RandomInt(-5, 10)
	c.Equal(num, number)
	c.Error(err)
	c.EqualError(ErrPositiveNumber, err.Error())

	number, err = RandomInt(10, 5)
	c.Equal(num, number)
	c.Error(err)
	c.EqualError(ErrMinValue, err.Error())
}

func TestRandomSSLRating(t *testing.T) {
	c := require.New(t)

	rate := RandomSSLRating("A")
	c.NotEmpty(rate)

	rate = RandomSSLRating("")
	c.NotEmpty(rate)
}

func TestRandomServerNumber(t *testing.T) {
	c := require.New(t)

	num := RandomServerNumber()
	c.NotEmpty(num)
}

func TestRandomNameDomain(t *testing.T) {
	c := require.New(t)

	nameDomain := RandomNameDomain()
	c.NotEmpty(nameDomain)
}

func BenchmarkRandomNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := RandomInt(1, 4)
		if err != nil {
			b.Fatal(err)
		}
	}
}
