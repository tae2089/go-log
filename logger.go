package golog

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var log *zap.Logger

func init() {
	var err error
	config := zap.NewProductionConfig()
	encoderConfig := createZapEncoderConfig()
	config.EncoderConfig = encoderConfig

	log, err = config.Build(zap.AddCallerSkip(1),
		zap.WrapCore(setupInfoFile),
		zap.WrapCore(setupErrorFile))

	if err != nil {
		panic(err)
	}
}

func Info(msg string, fields ...zapcore.Field) {
	log.Info(msg, fields...)
}

func Debug(msg string, fields ...zapcore.Field) {
	log.Debug(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	log.Error(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	log.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	log.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zapcore.Field) {
	log.Panic(msg, fields...)
}
func createZapEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""
	return encoderConfig
}

func setupInfoFile(c zapcore.Core) zapcore.Core {
	cfg := createZapEncoderConfig()

	logFile := "info-%Y-%m-%d.log"

	rotator, err := rotatelogs.New(
		logFile,
		rotatelogs.WithMaxAge(60*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour*24))
	if err != nil {
		panic(err)
	}
	w := zapcore.AddSync(rotator)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		w,
		zap.InfoLevel)
	cores := zapcore.NewTee(c, core)
	return cores
}

func setupErrorFile(c zapcore.Core) zapcore.Core {
	cfg := createZapEncoderConfig()

	logFile := "error-%Y-%m-%d.log"

	rotator, err := rotatelogs.New(
		logFile,
		rotatelogs.WithMaxAge(60*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour*24))
	if err != nil {
		panic(err)
	}
	w := zapcore.AddSync(rotator)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		w,
		zap.ErrorLevel)
	cores := zapcore.NewTee(c, core)
	return cores
}
