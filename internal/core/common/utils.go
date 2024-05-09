package common

import (
	"time"

	"github.com/artm2000/urlbook/pkg"
)

func GetUTCCurrentMillis() uint64 {
	return uint64(time.Now().UTC().UnixMilli())
}

func GetValidator() pkg.SuperValidator {
	return pkg.NewSuperValidator()
}
