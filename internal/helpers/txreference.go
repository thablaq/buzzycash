package helpers

import (
	"fmt"
	// "math/rand"
	// "time"
"github.com/google/uuid"
)


func GenerateTransactionReference(prefix ...string) string {
	p := "TRF"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	return fmt.Sprintf("%s|%s", p, uuid.New().String())
}