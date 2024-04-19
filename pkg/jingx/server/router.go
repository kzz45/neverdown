package server

import (
	"net/http"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	"github.com/gin-gonic/gin"
)

const (
	RouterAuthority = "/jingx/authz"
	RouterApi       = "/jingx/api/:token"
)

func (s *Server) registerRouters(router *gin.Engine) {
	router.POST(RouterAuthority, s.authority)
	router.GET(RouterApi, s.handler)
}

func (s *Server) handler(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	p, ok := s.members.ValidToken(token)
	if !ok {
		zaplogger.Sugar().Infof("Invalid token:%s", token)
		return
	}
	s.connections.Handler(c.Writer, c.Request, p)
}

func (s *Server) authority(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	res, err := s.members.Resource().Login(data)
	if err != nil {
		zaplogger.Sugar().Error(err)
	}
	c.Data(http.StatusOK, "multipart/form-data", res)
}
