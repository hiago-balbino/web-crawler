package main

import (
	"github.com/hiago-balbino/web-crawler/v2/cmd"
	"github.com/hiago-balbino/web-crawler/v2/internal/pkg/logger"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		logger.
			GetLogger().
			Fatal("error initializing the application", logger.FieldError(err))
	}
}
