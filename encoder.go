package zapcloudlogging

import (
	"strconv"
	"time"

	"go.uber.org/zap/zapcore"
)

var logLevelSeverity = map[zapcore.Level]string{
	zapcore.DebugLevel:  "DEBUG",
	zapcore.InfoLevel:   "INFO",
	zapcore.WarnLevel:   "WARNING",
	zapcore.ErrorLevel:  "ERROR",
	zapcore.DPanicLevel: "CRITICAL",
	zapcore.PanicLevel:  "ALERT",
	zapcore.FatalLevel:  "EMERGENCY",
}

// severityEncoder is an encoder for severity.
//
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logseverity
func severityEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(logLevelSeverity[l])
}

type sourceLocation struct {
	File     string
	Line     int
	Function string
}

func (l sourceLocation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("file", l.File)
	enc.AddString("line", strconv.Itoa(l.Line))
	enc.AddString("function", l.Function)
	return nil
}

// sourceLocationEncoder is a encoder for SourceLocation.
//
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logentrysourcelocation
func sourceLocationEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if aenc, ok := enc.(zapcore.ArrayEncoder); ok {
		aenc.AppendObject(sourceLocation{
			File:     caller.File,
			Line:     caller.Line,
			Function: caller.Function,
		})
	} else {
		enc.AppendString(caller.TrimmedPath())
	}
}

type timestamp struct {
	Seconds int64
	Nanos   int
}

func (t timestamp) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("seconds", t.Seconds)
	enc.AddInt("nanos", t.Nanos)
	return nil
}

// timestampEncoder is a encoder for timestamp.
//
// https://cloud.google.com/logging/docs/agent/logging/configuration#timestamp-processing
func timestampEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	if aenc, ok := enc.(zapcore.ArrayEncoder); ok {
		aenc.AppendObject(timestamp{
			Seconds: t.Unix(),
			Nanos:   t.Nanosecond(),
		})
	} else {
		zapcore.RFC3339NanoTimeEncoder(t, enc)
	}
}

// https://cloud.google.com/logging/docs/structured-logging
var encoderConfig = zapcore.EncoderConfig{
	MessageKey:     "message",
	LevelKey:       "severity",
	TimeKey:        "timestamp",
	NameKey:        "logger",
	CallerKey:      "logging.googleapis.com/sourceLocation",
	FunctionKey:    zapcore.OmitKey,
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    severityEncoder,
	EncodeTime:     timestampEncoder,
	EncodeDuration: zapcore.MillisDurationEncoder,
	EncodeCaller:   sourceLocationEncoder,
}

// NewProductionEncoderConfig returns a zapcore.EncoderConfig for production environments.
func NewProductionEncoderConfig() zapcore.EncoderConfig {
	return encoderConfig
}

// NewDevelopmentEncoderConfig returns a zapcore.EncoderConfig for development environments.
func NewDevelopmentEncoderConfig() zapcore.EncoderConfig {
	return encoderConfig
}
