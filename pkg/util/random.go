package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomInt(min, max int64) int64 {
	return min + rng.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rng.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomUsername() string {
	return fmt.Sprintf("user_%s", RandomString(5))
}

func RandomEmail() string {
	return fmt.Sprintf("%s@example.com", RandomString(6))
}

func RandomPassword() string {
	return RandomString(10)
}

func RandomFullName() string {
	return fmt.Sprintf("%s %s", RandomString(6), RandomString(6))
}
