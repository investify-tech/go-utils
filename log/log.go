package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

var logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

// LogInfo logs a message with the severity INFO.
func LogInfo(message string, args ...interface{}) {
	logger.Info().Msg(fmt.Sprintf(message, args...))
}

// LogWarn logs a message with the severity WARN.
func LogWarn(message string, args ...interface{}) {
	logger.Warn().Msg(fmt.Sprintf(message, args...))
}

// LogError logs a message with the severity ERROR.
func LogError(err error, message string, args ...interface{}) {
	logger.Err(err).Msg(fmt.Sprintf(message, args...))
}

// LogFatalAndQuit logs a message with the severity Fatal and quits the program execution.
func LogFatalAndQuit(err error, message string) {
	logger.Fatal().Msg(fmt.Sprintf(message+" - Error: %v", err))
}
