package server

import (
	"net/http"

	"github.com/kzz45/neverdown/pkg/jwttoken"
	"k8s.io/klog/v2"

	"github.com/gin-gonic/gin"
)

const (
	RouterApi    = "/authx/authz"
	RouterHealth = "/authx/healthz"
)

func (s *Server) registerRouters(router *gin.Engine) {
	router.GET(RouterHealth, s.healthz)
	router.POST(RouterApi, s.authz)
}

func (s *Server) healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "pong",
	})
}

func (s *Server) authz(c *gin.Context) {
	token := c.Request.Header.Get(jwttoken.TokenKey)
	data, err := c.GetRawData()
	if err != nil {
		klog.Error(err)
		return
	}
	// klog.Infof("token: %s, data: %v", token, string(data))
	res, err := s.handler.Handle(token, data)
	if err != nil {
		klog.Error(err)
	}
	c.Data(http.StatusOK, "multipart/form-data", res)
}
