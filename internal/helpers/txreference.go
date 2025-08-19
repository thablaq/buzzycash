package helpers

import (
	"fmt"
	"math/rand"
	"time"
)


func GenerateTransactionReference(prefix ...string) string {
	p := "TRF"
	if len(prefix) > 0 {
		p = prefix[0]
	}

	// Random alphanumeric string
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomPart := make([]byte, 9)
	for i := range randomPart {
		randomPart[i] = charset[rand.Intn(len(charset))]
	}

	timestamp := time.Now().UnixMilli()
	reference := fmt.Sprintf("%s|%d|%s", p, timestamp, string(randomPart))
	return reference
}
