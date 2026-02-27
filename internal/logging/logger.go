package logging

import (
	"fmt"

	"github.com/wb-go/wbf/logger"
)

type URLShortenerLogger struct {
	logger.Logger
}

var AppLogger, _ *URLShortenerLogger

func NewURLShortenerLogger() (*URLShortenerLogger, error) {
	Logger, err := logger.InitLogger(logger.ZapEngine, "URLShortener", "")
	if err != nil {
		fmt.Errorf("failed to initialize logger: %v", err)
		return nil, err
	}
	return &URLShortenerLogger{
		Logger: Logger,
	}, nil
}
