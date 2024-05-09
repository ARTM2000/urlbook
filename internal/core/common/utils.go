package common

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/artm2000/urlbook/pkg"
)

func GetUTCCurrentMillis() uint64 {
	return uint64(time.Now().UTC().UnixMilli())
}

func GetValidator() pkg.SuperValidator {
	return pkg.NewSuperValidator()
}

func GetRandomUrlShortPhrase() string {
	phraseLength := 7
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	finalPhrase := ""

	for i := 0; i < phraseLength; i++ {
		c := chars[rand.Intn(len(chars))]
		finalPhrase = fmt.Sprintf("%s%s", finalPhrase, string(c))
	}

	return finalPhrase
}
