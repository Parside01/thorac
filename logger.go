package main

import "go.uber.org/zap"

func NewLogger(config *Config) *zap.Logger {
	logger := zap.NewExample().WithOptions(zap.AddCallerSkip(1))
	return logger.With(
		zap.String("component", "raft"),
		zap.String("node_id", config.Id),
		zap.String("advertise_address", config.Advertise),
	)
}
