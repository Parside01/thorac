package core

import "go.uber.org/zap"

var L *zap.Logger

func InitLogger() error {
	L, _ = zap.NewProduction(zap.AddCaller())
	return nil
}
