package boutique

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func (q Quote) String() string {
	return fmt.Sprintf("$%d.%d", q.Dollars, q.Cents)
}

// CreateQuoteFromCount takes a number of items and returns a Quote struct.
func CreateQuoteFromCount(count int) Quote {
	return CreateQuoteFromFloat(8.99 * float64(count))
}

// CreateQuoteFromFloat takes a price represented as a float and creates a Price struct.
func CreateQuoteFromFloat(value float64) Quote {
	units, fraction := math.Modf(value)
	return Quote{
		int32(units),
		int64(math.Trunc(fraction * 100)),
	}
}

// seeded determines if the random number generator is ready.
var seeded bool = false

// CreateTrackingId generates a tracking ID.
func CreateTrackingId(salt string) string {
	if !seeded {
		rand.Seed(time.Now().UnixNano())
		seeded = true
	}
	return fmt.Sprintf("%c%c-%d%s-%d%s",
		getRandomLetterCode(),
		getRandomLetterCode(),
		len(salt),
		getRandomNumber(3),
		len(salt)/2,
		getRandomNumber(7),
	)
}

// getRandomLetterCode generates a code point value for a capital letter.
func getRandomLetterCode() uint32 {
	return 65 + uint32(rand.Intn(25))
}

// getRandomNumber generates a string representation of a number with the requested number of digits.
func getRandomNumber(digits int) string {
	str := ""
	for i := 0; i < digits; i++ {
		str = fmt.Sprintf("%s%d", str, rand.Intn(10))
	}
	return str
}

func GetQuote(ctx context.Context, items []CartItem) Money {
	// 1. Generate a quote based on the total number of items to be shipped.
	quote := CreateQuoteFromCount(len(items))
	// 2. Generate a response.
	return Money{
		Currency: "USD",
		Units:    quote.Dollars,
		Nanos:    quote.Cents * 10000000,
	}
}

func ShipOrder(ctx context.Context, address Address, items []CartItem) string {
	// 1. Create a Tracking ID
	baseAddress := fmt.Sprintf("%s, %s, %s", address.StreetAddress, address.City, address.State)
	id := CreateTrackingId(baseAddress)
	// 2. Generate a response.
	return id
}
