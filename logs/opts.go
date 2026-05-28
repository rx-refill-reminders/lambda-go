package logs

import (
	"io"
	"log/slog"
)

type LoggerStyle int

const (
	StyleJSON LoggerStyle = iota
	StyleText
)

func (ls LoggerStyle) String() string {
	switch ls {
	case StyleJSON:
		return "json"
	case StyleText:
		return "text"
	default:
		return "unknown"
	}
}

type LoggerOpts struct {
	// Where the logs go
	Out io.Writer

	// Defaults to json
	Style LoggerStyle

	// Defaults to info
	Level slog.Level
}
