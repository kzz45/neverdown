package env

import (
	"os"
	"strconv"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	"go.uber.org/zap"
)

const (
	WebsocketPrintWriteData   = "WEBSOCKET_PRINT_WRITE"
	WebsocketPrintReadData    = "WEBSOCKET_PRINT_READ"
	WebsocketKeepAliveTimeout = "WEBSOCKET_KEEPALIVE_TIMEOUT"
)

const (
	Show   = true
	Hidden = false
)

// GetWebsocketWrite returns the debug mode of the websocket conns
func GetWebsocketWrite() bool {
	a := os.Getenv(WebsocketPrintWriteData)
	if a == "" {
		zaplogger.Sugar().Infof("no env variable getting from '%s' in the container", WebsocketPrintWriteData)
		return Hidden
	}
	t, err := strconv.ParseBool(a)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return Hidden
	}
	return t
}

// GetWebsocketRead returns the debug mode of the websocket conns
func GetWebsocketRead() bool {
	a := os.Getenv(WebsocketPrintReadData)
	if a == "" {
		zaplogger.Sugar().Infof("no env variable getting from '%s' in the container", WebsocketPrintReadData)
		return Hidden
	}
	t, err := strconv.ParseBool(a)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return Hidden
	}
	return t
}

// GetWebsocketKeepaliveTimeout
func GetWebsocketKeepaliveTimeout(defaultValue int64) int64 {
	a := os.Getenv(WebsocketKeepAliveTimeout)
	if a == "" {
		return defaultValue
	}
	t, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		zaplogger.Sugar().Errorw("GetWebsocketKeepaliveTimeout",
			zap.Int64("default", defaultValue),
			zap.Error(err),
		)
		return defaultValue
	}
	return t
}
