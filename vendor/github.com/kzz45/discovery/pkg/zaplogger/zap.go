package zaplogger

import (
	"log"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	EnvZapEncoding  = "ENV_ZAP_ENCODING"
	EnvZapLevel     = "ENV_ZAP_LEVEL"
	EncodingJson    = "json"
	EncodingConsole = "console"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func init() {
	var err error
	logger, err = NewProductionConfig().Build()
	if err != nil {
		log.Panic(err)
	}
	sugar = logger.Sugar()
}

func getEncodingFromEnv() string {
	encoding := os.Getenv(EnvZapEncoding)
	switch encoding {
	case "", EncodingConsole:
		encoding = EncodingConsole
	case EncodingJson:
		encoding = EncodingJson
	default:
		encoding = EncodingConsole
	}
	return encoding
}

func getLevelFromEnv() zapcore.Level {
	lv := os.Getenv(EnvZapLevel)
	if len(lv) == 0 {
		return zapcore.InfoLevel
	}
	t, err := strconv.ParseInt(lv, 10, 8)
	if err != nil {
		panic(err)
	}
	a := zapcore.Level(int8(t))
	if a < zapcore.DebugLevel || a > zapcore.FatalLevel {
		panic("invalid ENV_ZAP_LEVEL")
	}
	return a
}

func Sync() {
	_ = logger.Sync()
}

func Sugar() *zap.SugaredLogger {
	return sugar
}

func Logger() *zap.Logger {
	return logger
}

func NewProductionConfig() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(getLevelFromEnv()),
		Development: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         getEncodingFromEnv(),
		EncoderConfig:    NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func NewProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		// Keys can be anything except the empty string.
		//TimeKey:        "T",
		//LevelKey:       "L",
		//NameKey:        "N",
		//CallerKey:      "C",
		//FunctionKey:    zapcore.OmitKey,
		//MessageKey:     "M",
		//StacktraceKey:  "S",
	}
}

func EpochTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//nanos := t.UnixNano()
	//sec := float64(nanos) / float64(time.Second)
	//enc.AppendFloat64(sec)
	enc.AppendString(t.Format("2006-01-02T15:04:05.000000Z"))
}
