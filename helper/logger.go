package helper

import (
	"os"
	"time"
	"io"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func InitLogger() {
	logDir := "logs"
	logFile := logDir + "/app.log"

	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	writer := io.MultiWriter(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		},
		file,
	)

	Log = zerolog.New(writer).
		With().
		Timestamp().
		Str("app", "uas_be").
		Logger()
}