package logger

import (
	"github.com/minhmannh2001/mongo-change-stream-processor/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// For mapping config logger to app logger levels
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// Logger methods interface
type Logger interface {
	InitLogger()
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	DPanicf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	WithFiled(field zapcore.Field) *zap.Logger
}

// Logger
type appLogger struct {
	cfg         *config.Configuration
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
}

// App Logger constructor
func NewApiLogger(cfg ...*config.Configuration) Logger {
	applg := &appLogger{}
	if len(cfg) == 0 {
		applg.DefaultInit()
	} else {
		applg.cfg = cfg[0]
		applg.InitLogger()
	}
	return applg
}
