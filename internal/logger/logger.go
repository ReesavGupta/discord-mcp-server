package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Level    string
	Format   string
	FilePath string
}

func NewLogger(config Config) (*logrus.Logger, error) {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set formatter
	if config.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Set output
	if config.FilePath != "" {
		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(config.FilePath), 0755); err != nil {
			return nil, err
		}

		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}

		// Write to both file and stderr
		logger.SetOutput(io.MultiWriter(file, os.Stderr))
	} else {
		logger.SetOutput(os.Stderr)
	}

	return logger, nil
}
