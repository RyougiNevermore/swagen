package zlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.SugaredLogger

func New(verbose bool, version string) {
	var consoleEncoderConfig zapcore.EncoderConfig
	if verbose {
		consoleEncoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		consoleEncoderConfig = zap.NewProductionEncoderConfig()
	}
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	core := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	}))

	tee := zapcore.NewTee(core)

	zLog := zap.New(tee)
	zLog = zLog.With(zap.Namespace("swagen"), zap.String("version", version))

	zLog = zLog.WithOptions(zap.AddCaller())

	logger = zLog.Sugar()

}

func Log() *zap.SugaredLogger {
	return logger
}

func Sync() {
	err := logger.Sync()
	if err != nil {
		panic(err)
	}
}
