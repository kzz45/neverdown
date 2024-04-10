package handler

import (
	"fmt"
	"time"

	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	"github.com/kzz45/neverdown/pkg/authx/http/proto"
	"github.com/kzz45/neverdown/pkg/authx/validation"

	"github.com/kzz45/neverdown/pkg/jwttoken"
	"github.com/kzz45/neverdown/pkg/zaplogger"
)

func (h *Handler) login(data []byte) (res []byte, err error) {
	meta := &rbacv1.AccountMeta{}
	if err = meta.Unmarshal(data); err != nil {
		return nil, err
	}
	if err = validation.Account(meta.Username); err != nil {
		return nil, err
	}
	isRoot := isAdmin(meta.Username, meta.Password)
	zaplogger.Sugar().Infof("isRoot: %v", isRoot)

	if !isRoot {
		bo := h.adminRbacServiceAccount.Validate(meta.Username, meta.Password)
		if !bo {
			return nil, fmt.Errorf("invalid username or password")
		}
	}
	token, err := jwttoken.Generate(meta.Username, isRoot)
	if err != nil {
		return nil, err
	}
	ctx := &proto.Context{
		Token:    token,
		IsAdmin:  isRoot,
		ExpireAt: int32(time.Now().Unix() + jwttoken.GetTokenExpirationFromEnv()),
	}
	return ctx.Marshal()
}

const (
	RootUsername = "admin"
	RootPassword = "Uh4zIyIQqrDB0lrdYA2jaOtF1DBobqla"
)

func isAdmin(username, password string) bool {
	return username == RootUsername && password == RootPassword
}

func (h *Handler) validate(token string) (*jwttoken.Claims, error) {
	if token == "" {
		return nil, fmt.Errorf("nil token")
	}
	return jwttoken.Parse(token)
}
