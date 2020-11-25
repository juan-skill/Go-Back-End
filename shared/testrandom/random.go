package testrandom

import (
	"errors"

	"crypto/rand"
	"math/big"

	"github.com/other_project/crockroach/shared/env"
)

var (
	// ErrMinValue when select the min and max value
	ErrMinValue = errors.New("cannot be min value greater than max value")
	// ErrPositiveNumber when it pass negative number
	ErrPositiveNumber = errors.New("cannot be a negative number")
	// minNumber determinate the minimum number
	minNumber = env.GetInt64("MIN_NUMBER", 1)
	// maxNumber determinate the maximum number
	maxNumber = env.GetInt64("MIN_NUMBER", 5)
	// ErrRandomNumber when it's generating
	ErrRandomNumber = errors.New("cannot generate random number")
)

func init() {
	//rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) (int64, error) {
	if min > max {
		return 0, ErrMinValue
	}

	if min < 0 || max < 0 {
		return 0, ErrPositiveNumber
	}

	number, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
	if err != nil {
		return 0, ErrRandomNumber
	}

	//return min + rand.Int63n(max-min+1), nil
	return min + number.Int64(), nil
}

// RandomSSLRating generates a random SSL rating
func RandomSSLRating(rate string) string {
	ratings := []string{"A", "A+", "A-", "B", "B+", "B-", "C", "C+", "C-", "D", "D+", "D-", "E", "E+", "E-", "F", "F+", "F-", "M", "T"}
	n := len(ratings)

	number, _ := RandomInt(0, int64(n-1))

	newrate := ratings[number]

	if newrate != rate {
		return newrate
	}

	return RandomSSLRating(rate)
}

// RandomServerNumber generates a random number of servers
func RandomServerNumber() int64 {
	number, err := RandomInt(1, 4)
	if err != nil {
		return number
	}

	number, _ = RandomInt(minNumber, maxNumber)

	return number
}
