package server

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kzz45/neverdown/pkg/jwttoken"
	"github.com/kzz45/neverdown/pkg/zaplogger"
	"go.uber.org/zap"

	rbacv1 "github.com/kzz45/discovery/pkg/apis/rbac/v1"
	"github.com/kzz45/neverdown/pkg/execpod"

	"github.com/kzz45/neverdown/pkg/openx/aggregator/proto"
	"github.com/kzz45/neverdown/pkg/openx/aggregator/resources"
)

const (
	RouterAuthority      = "/authority"
	RouterApi            = "/api/:token"
	RouterPodSSH         = "/ssh/namespace/:namespace/pod/:pod/shell/:container/:command/token/:token"
	RouterPodLogStream   = "/log/stream/namespace/:namespace/pod/:pod/container/:container/sinceSeconds/:SinceSeconds/sinceTime/:SinceTime/token/:token"
	RouterPodLogDownload = "/log/download/namespace/:namespace/pod/:pod/container/:container/previous/:previous/sinceSeconds/:SinceSeconds/sinceTime/:SinceTime"
)

func (s *Server) registerRouters(router *gin.Engine) {
	router.POST(RouterAuthority, s.authority)
	router.GET(RouterApi, s.handler)
	router.GET(RouterPodSSH, s.ssh)
	router.GET(RouterPodLogStream, s.logStreaming)
	router.GET(RouterPodLogDownload, s.logDownload)
}

func (s *Server) authority(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	res, err := s.resources.Login(data)
	if err != nil {
		zaplogger.Sugar().Error(err)
	}
	c.Data(http.StatusOK, "multipart/form-data", res)
}

func (s *Server) handler(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	p, ok := s.members.ValidToken(token)
	if !ok {
		zaplogger.Sugar().Infof("Invalid token:%s", token)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	s.connections.Handler(c.Writer, c.Request, p)
}

func ValidToken(token string) (*jwttoken.Claims, error) {
	claims, err := jwttoken.Parse(token)
	if err != nil {
		zaplogger.Sugar().Error("ValidToken err",
			zap.String("token", token),
			zap.Error(err))
		return nil, err
	}
	return claims, nil
}

func setOptionWithSince(c *gin.Context, opt *execpod.ExecOptions) error {
	// check `sinceSeconds` and `sinceTime`
	sinceSec, err := strconv.Atoi(c.Param("SinceSeconds"))
	if err != nil {
		zaplogger.Sugar().Errorw("Convert sinceSeconds failed", "SinceSeconds", c.Param("SinceSeconds"), "err", err)
		return err
	}
	if sinceSec > 0 {
		a := int64(sinceSec)
		opt.SinceSeconds = &a
		opt.SinceTime = nil
	} else {
	}
	return nil
}

func (s *Server) ssh(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	claims, err := ValidToken(token)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	options := &execpod.ExecOptions{
		Namespace:     c.Param("namespace"),
		PodName:       c.Param("pod"),
		ContainerName: c.Param("container"),
		Follow:        true,
		Command:       []string{c.Param("command")},
	}
	bo, err := s.resources.ValidateAccess(
		claims.Username,
		options.Namespace,
		rbacv1.GroupVersionKind{
			Group:   resources.PodsGroupVersionKind.Group,
			Version: resources.PodsGroupVersionKind.Version,
			Kind:    resources.PodsGroupVersionKind.Kind,
		},
		string(proto.VerbPodsSSH),
	)
	if err != nil {
		zaplogger.Sugar().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !bo {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	s.podshandler.SSH(c.Writer, c.Request, options)
}

func (s *Server) logStreaming(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	claims, err := ValidToken(token)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	options := &execpod.ExecOptions{
		Namespace:     c.Param("namespace"),
		PodName:       c.Param("pod"),
		ContainerName: c.Param("container"),
		Follow:        true,
		Command:       []string{c.Param("command")},
	}
	bo, err := s.resources.ValidateAccess(
		claims.Username,
		options.Namespace,
		rbacv1.GroupVersionKind{
			Group:   resources.PodsGroupVersionKind.Group,
			Version: resources.PodsGroupVersionKind.Version,
			Kind:    resources.PodsGroupVersionKind.Kind,
		},
		string(proto.VerbPodsLogStreaming),
	)
	if err != nil {
		zaplogger.Sugar().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !bo {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err = setOptionWithSince(c, options); err != nil {
		c.Abort()
		return
	}
	s.podshandler.LogStreaming(c.Writer, c.Request, options)
}

func (s *Server) logDownload(c *gin.Context) {
	token := c.Request.Header.Get(jwttoken.TokenKey)
	if token == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	claims, err := ValidToken(token)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	pre, err := strconv.ParseBool(c.Param("previous"))
	if err != nil {
		zaplogger.Sugar().Error(err)
		c.Abort()
		return
	}
	options := &execpod.ExecOptions{
		Namespace:       c.Param("namespace"),
		PodName:         c.Param("pod"),
		ContainerName:   c.Param("container"),
		Follow:          false,
		UsePreviousLogs: pre,
	}
	bo, err := s.resources.ValidateAccess(
		claims.Username,
		options.Namespace,
		rbacv1.GroupVersionKind{
			Group:   resources.PodsGroupVersionKind.Group,
			Version: resources.PodsGroupVersionKind.Version,
			Kind:    resources.PodsGroupVersionKind.Kind,
		},
		string(proto.VerbPodsLogDownload),
	)
	if err != nil {
		zaplogger.Sugar().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !bo {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	reader, err := execpod.LogDownload(s.podshandler.ClientSet, options)
	if err != nil {
		zaplogger.Sugar().Error(err)
		c.Abort()
		return
	}
	defer func() {
		zaplogger.Sugar().Info("LogTransmit readCloser close")
		if err := reader.Close(); err != nil {
			zaplogger.Sugar().Error(err)
		}
	}()
	fileContentDisposition := fmt.Sprintf("attachment;filename=%s_%s_%s.log", options.Namespace, options.PodName, options.ContainerName)
	c.Header("Content-Type", "text/plain")
	c.Header("Content-Disposition", fileContentDisposition)
	if _, err = io.Copy(c.Writer, reader); err != nil {
		zaplogger.Sugar().Error(err)
	}
}
