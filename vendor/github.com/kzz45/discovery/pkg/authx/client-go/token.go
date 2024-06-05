package authority

import (
	"fmt"

	"github.com/kzz45/discovery/pkg/jwttoken"
)

func Token(token string) (*jwttoken.Claims, error) {
	if token == "" {
		return nil, fmt.Errorf("nil token")
	}
	return jwttoken.Parse(token)
}
