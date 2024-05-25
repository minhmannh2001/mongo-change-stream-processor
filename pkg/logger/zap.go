package logger

import (
	"os"

	"github.com/minhmannh2001/mongo-change-stream-processor/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (l *appLogger) WithFiled(field zap.Field) *zap.Logger {
	return l.logger.With(field)
}

func (l *appLogger) getLoggerLevel(cfg *config.Configuration) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}
func (l *appLogger) DefaultInit() {
	logger, _ := zap.NewProduction()
	l.logger = logger
	l.sugarLogger = logger.Sugar()
	if err := l.sugarLogger.Sync(); err != nil {
		l.sugarLogger.Error(err)
	}
}

func (l *appLogger) InitLogger() {
	logLevel := l.getLoggerLevel(l.cfg)

	logWriter := zapcore.AddSync(os.Stdout)

	var encoderCfg zapcore.EncoderConfig
	if l.cfg.AppConfig.Mode == "Development" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	var encoder zapcore.Encoder
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	if l.cfg.Logger.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	l.logger = logger
	l.sugarLogger = logger.Sugar()
	if err := l.sugarLogger.Sync(); err != nil {
		l.sugarLogger.Error(err)
	}
}

// Logger methods
func (l *appLogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *appLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *appLogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *appLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *appLogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

func (l *appLogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

func (l *appLogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}
