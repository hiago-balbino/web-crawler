package main

import (
	"github.com/hiago-balbino/web-crawler/cmd"
	"github.com/hiago-balbino/web-crawler/internal/pkg/logger"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		logger.
			GetLogger().
			Fatal("error initializing the application", logger.FieldError(err))
	}
}
