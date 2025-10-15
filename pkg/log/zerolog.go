package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
	"time"
)

func newZerolog() zerolog.Logger {
	// 自定义控制台格式
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
		FormatTimestamp: func(i interface{}) string {
			return time.Now().Format("2006-01-02 15:04:05,000")
		},
		NoColor: true,
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			"request_id",
			zerolog.CallerFieldName,
			"user",
			zerolog.MessageFieldName,
		},
		FieldsExclude: []string{"user", "request_id"},
		FormatCaller: func(i interface{}) string {
			var c string
			if cc, ok := i.(string); ok {
				c = cc
			}
			if len(c) > 0 {
				if cwd, err := os.Getwd(); err == nil {
					if rel, err := filepath.Rel(cwd, c); err == nil {
						c = rel
					}
				}
			}
			return c + ";"
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%s;", i)
		},
	}

	// 创建logger
	logger := zerolog.New(consoleWriter).
		With().
		Timestamp().
		CallerWithSkipFrameCount(3).
		Logger()
	return logger
}
