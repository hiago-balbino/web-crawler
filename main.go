package main

import (
	"github.com/hiago-balbino/web-crawler/cmd"
	"github.com/hiago-balbino/web-crawler/internal/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		logger.
			GetLogger().
			Fatal("error initializing the application", zapcore.Field{Type: zapcore.StringType, String: err.Error()})
	}
}
