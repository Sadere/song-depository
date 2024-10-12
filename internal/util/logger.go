package util

import "go.uber.org/zap"

// Returns new instance of zap sugar logger
func NewZapLogger(level string) (*zap.SugaredLogger, error) {
	// parse string log level to zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	// create new logger config
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	// build logger from config
	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return zl.Sugar(), nil
}
