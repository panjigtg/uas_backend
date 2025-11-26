package helper

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func InitLogger() {
	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	Log = zerolog.New(writer).
		With().
		Timestamp().
		Str("app", "uas-prestasi").
		Logger()
}