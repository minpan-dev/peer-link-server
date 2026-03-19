package logger

import (
	"peer-link-server/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg *config.LogConfig) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	var zapCfg zap.Config
	if cfg.Format == "console" {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	zapCfg.DisableStacktrace = true
	return zapCfg.Build()
}

func MustNew(cfg *config.LogConfig) *zap.Logger {
	l, err := New(cfg)
	if err != nil {
		panic("init logger: " + err.Error())
	}
	return l
}
