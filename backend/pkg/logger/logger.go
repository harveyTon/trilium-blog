package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger zerolog.Logger

func Init(level string) {
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal().Err(err).Msg("Unable to create log directory")
	}

	currentLogFile := filepath.Join(logDir, "current.log")

	if info, err := os.Stat(currentLogFile); err == nil {
		lastModified := info.ModTime()
		if lastModified.Day() != time.Now().Day() {
			newName := filepath.Join(logDir, lastModified.Format("2006-01-02")+".log")
			if err := os.Rename(currentLogFile, newName); err != nil {
				log.Error().Err(err).Msg("Failed to rename old log file")
			}
		}
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   currentLogFile,
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
	}

	multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, lumberjackLogger)

	zerolog.SetGlobalLevel(parseLevel(level))
	Logger = zerolog.New(multi).With().Timestamp().Caller().Logger()

	log.Logger = Logger
}

func parseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

// GinLogger returns Gin's logging middleware
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)

		if len(c.Errors) > 0 {
			errors := make([]error, len(c.Errors))
			for i, err := range c.Errors {
				errors[i] = err
			}
			Logger.Error().Errs("Gin errors", errors).Str("path", path).Str("query", query).Send()
		} else {
			Logger.Info().
				Str("ip", c.ClientIP()).
				Str("method", c.Request.Method).
				Str("path", path).
				Int("status", c.Writer.Status()).
				Dur("latency", latency).
				Str("query", query).
				Msg("HTTP request")
		}
	}
}

func Debug(msg string) {
	Logger.Debug().Msg(msg)
}

func Info(msg string) {
	Logger.Info().Msg(msg)
}

func Warn(msg string) {
	Logger.Warn().Msg(msg)
}

func Error(msg string, err error) {
	Logger.Error().Err(err).Msg(msg)
}

func Fatal(msg string, err error) {
	Logger.Fatal().Err(err).Msg(msg)
}

func Debugf(format string, v ...interface{}) {
	Logger.Debug().Msg(fmt.Sprintf(format, v...))
}

func Infof(format string, v ...interface{}) {
	Logger.Info().Msg(fmt.Sprintf(format, v...))
}

func Warnf(format string, v ...interface{}) {
	Logger.Warn().Msg(fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {
	Logger.Error().Msg(fmt.Sprintf(format, v...))
}

func Fatalf(format string, v ...interface{}) {
	Logger.Fatal().Msg(fmt.Sprintf(format, v...))
}
