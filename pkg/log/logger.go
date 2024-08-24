package log

import "go.uber.org/zap"

var (
	logger *zap.SugaredLogger
)

func L() *zap.SugaredLogger {
	if logger == nil {
		z, _ := zap.NewProduction()
		logger = z.Sugar()
	}

	return logger
}