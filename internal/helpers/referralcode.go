package helpers 


import (
	"math/rand"
	"strings"
	"time"
)

func GenerateReferralCode(prefix ...string) string {
	
	p := "BZ"
	if len(prefix) > 0 {
		p = prefix[0]
	}


	rand.New(rand.NewSource(time.Now().UnixNano()))

	
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const randomLength = 6


	randomPart := make([]byte, randomLength)
	for i := range randomPart {
		randomPart[i] = chars[rand.Intn(len(chars))]
	}

	referralCode := strings.ToUpper(p)  + string(randomPart)
	return referralCode
}