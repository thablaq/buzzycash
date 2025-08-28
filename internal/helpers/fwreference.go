package helpers



import (
	"math/rand"
	"time"
	"fmt"
)

func GenerateFWRef() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%08d", rand.Intn(100000000))
}
