package random

import (
	"github.com/thanhpk/randstr"
)

func GenRandomString(len int) string {
	return randstr.String(len)
}
