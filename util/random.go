package util

import(
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init(){
	rand.Seed(time.Now().UnixNano())
}

//Generating a random number.
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

//Generating a random string.
func RandomString(n int) string {

	// Standard GO function to build a string.
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i ++ {
		// Picks a random letter.
		c := alphabet[rand.Intn(k)]

		// Writes it to the string.
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random name.
func RandomOwner() string {
	return RandomString(6)
}

// Random money generator.
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// Random Currency
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}





