package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log   *zap.Logger
	sugar *zap.SugaredLogger
)

func InitLogger(logLevel string) {
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:    "time",
			LevelKey:   "level",
			NameKey:    "logger",
			MessageKey: "message",
			StacktraceKey: func() string {
				if logLevel == "debug" {
					return "stacktrace"
				}
				return ""
			}(),
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    map[string]interface{}{"app": "git-sync"},
	}

	var err error
	log, err = cfg.Build()
	if err != nil {
		panic(err)
	}
	sugar = log.Sugar()
	sugar.Debug("Logger initialized")
}

func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

func Info(args ...interface{}) {
	sugar.Info(args...)
}

func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

func Error(args ...interface{}) {
	sugar.Error(args...)
}

func Fatal(args ...interface{}) {
	sugar.Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	sugar.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	sugar.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	sugar.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	sugar.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	sugar.Fatalf(format, args...)
}
