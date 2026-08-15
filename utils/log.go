package utils

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// this will write to a json file using the slog json handler and also prints to console using pretty handler
type Loggy struct {
	logDir      string
	jsonHandler slog.Handler
	Level       slog.Level

	attrs  []slog.Attr
	groups []string
}

func SetupLogger(LogDir string) {
	// NOTE: use the models.EnvVars to get the log level and set it here
	myLogger := Loggy{
		logDir: "logs",
		Level:  slog.LevelDebug,
	}
	// create a file inisde the logDir with the name of the current date and time in the format of YYYY-MM-DD_HH-MM-SS.json
	handler, err := createJsonHandler(myLogger.logDir, myLogger.Level)
	if err != nil {
		slog.Error("failed to create log file", "err", err)
		panic(err)
	}
	myLogger.jsonHandler = handler

	logger := slog.New(myLogger)
	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

func createJsonHandler(logDir string, level slog.Level) (*slog.JSONHandler, error) {
	fileName := filepath.Join(
		logDir,
		"logs_"+time.Now().Format("2006-01-02_15-04-05")+".json",
	)
	// create the logDir if it doesn't exist
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			slog.Error("failed to create log directory", "err", err)
			return nil, err
		}
	}

	f, err := os.Create(fileName)
	if err != nil {
		slog.Error("", "err", err)
		return nil, err
	}
	jsonHandler := slog.NewJSONHandler(f, &slog.HandlerOptions{
		AddSource: true,
		Level:     level, // NOTE: use the models.EnvVars to get the log level and set it here
	})
	return jsonHandler, nil
}

// Enabled implements slog.Handler.
func (l Loggy) Enabled(ctx context.Context, level slog.Level) bool {
	return l.jsonHandler.Enabled(ctx, level)
}

func (l Loggy) Handle(ctx context.Context, r slog.Record) error {
	msg := FormatLogEntryWithSource(r)

	for _, attr := range l.attrs {
		msg += fmt.Sprintf(
			"%v: %v\t",
			ColorByLevel(attr.Key, r.Level, false, false),
			attr.Value,
		)
	}

	fmt.Println(msg)

	return l.jsonHandler.Handle(ctx, r)
}

func (l Loggy) WithAttrs(attrs []slog.Attr) slog.Handler {
	l.jsonHandler = l.jsonHandler.WithAttrs(attrs)

	// Copy so handlers don't accidentally share/mutate slices.
	l.attrs = append(
		append([]slog.Attr{}, l.attrs...),
		attrs...,
	)

	return l
}

func (l Loggy) WithGroup(name string) slog.Handler {
	if name == "" {
		return l
	}

	l.jsonHandler = l.jsonHandler.WithGroup(name)

	l.groups = append(
		append([]string{}, l.groups...),
		name,
	)

	return l
}
func ColorByLevel(text string, rec slog.Level, isBold, isBG bool) string {
	if isBold {
		text = MakeBold(text)
	}
	if isBG {
		switch rec {
		case slog.LevelDebug:
			return ChangeColor(text, 46)
		case slog.LevelInfo:
			return ChangeColor(text, 42)
		case slog.LevelWarn:
			return ChangeColor(text, 43)
		case slog.LevelError:
			return ChangeColor(text, 41)
		default:
			return rec.String()
		}
	}
	switch rec {
	case slog.LevelDebug:
		return ChangeColor(text, 36)
	case slog.LevelInfo:
		return ChangeColor(text, 32)
	case slog.LevelWarn:
		return ChangeColor(text, 33)
	case slog.LevelError:
		return ChangeColor(text, 31)
	default:
		return rec.String()
	}
}

func ChangeColor(str string, color int) string {
	return fmt.Sprintf("\033[%vm%v\033[0m", color, str)
}

func MakeBold(str string) string {
	return fmt.Sprintf("\033[1m%v\033[0m", str)
}

func FormatLogEntryWithSource(rec slog.Record) string {
	// Extract the source file and line number
	wd, _ := os.Getwd()
	source := ""
	if rec.PC != 0 {
		fn := runtime.FuncForPC(rec.PC)
		if fn != nil {
			file, line := fn.FileLine(rec.PC)
			source = fmt.Sprintf("%s:%d", strings.TrimPrefix(file, wd+"/"), line)
		}
	}

	// Build the log entry
	result := fmt.Sprintf(
		"[%v]\t%v\t%v\t%v\t",
		rec.Time.Format("2006-01-02 15:04:05"),
		ColorByLevel(rec.Level.String(), rec.Level, true, false),
		source, // Include source info
		rec.Message,
	)
	rec.Attrs(func(a slog.Attr) bool {
		result += fmt.Sprintf("%v: %v\t", ColorByLevel(a.Key, rec.Level, false, false), a.Value)
		return true
	})
	return result
}
