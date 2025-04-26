package l

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _bufPool = buffer.NewPool()

var _consolePool = sync.Pool{New: func() interface{} {
	return &encoder{}
}}

var (
	enableRedact      = os.Getenv("LOG_REDACT_DISABLED") != "true"
	defaultRedactKeys = GenerateKeywordVariations()
)

// Map for O(1) lookups of keys to redact
var redactKeysMap = make(map[string]struct{})

func init() {
	listKeyNeedToRedact := os.Getenv("LOG_REDACT_KEYS")

	// If custom keys are provided, use them
	if listKeyNeedToRedact != "" {
		keys := strings.Split(listKeyNeedToRedact, ",")
		for _, key := range keys {
			trimmedKey := strings.TrimSpace(key)
			if trimmedKey != "" {
				// Use lowercase for case-insensitive comparison
				redactKeysMap[strings.ToLower(trimmedKey)] = struct{}{}
			}
		}
	}

	// If no custom keys or empty list, use defaults
	if len(redactKeysMap) == 0 {
		for _, key := range defaultRedactKeys {
			redactKeysMap[strings.ToLower(key)] = struct{}{}
		}
	}
}

func getEncoder() *encoder {
	if ret, ok := _consolePool.Get().(*encoder); ok {
		return ret
	}

	// Create a new encoder if none is available in the pool
	return &encoder{
		objectEncoder: &objectEncoder{},
		lineEnding:    tokenLineEnding,
	}
}

type encoder struct {
	*zapcore.EncoderConfig
	*objectEncoder
	lineEnding string
}

func newEncoder(cfg zapcore.EncoderConfig) *encoder {
	lineEnding := tokenLineEnding
	if len(cfg.LineEnding) > 0 {
		lineEnding = cfg.LineEnding
	}
	return &encoder{
		EncoderConfig: &cfg,
		objectEncoder: getObjectEncoder(&cfg),
		lineEnding:    lineEnding,
	}
}

func (enc *encoder) Clone() zapcore.Encoder {
	clone := getEncoder()
	clone.EncoderConfig = enc.EncoderConfig
	clone.lineEnding = enc.lineEnding
	clone.objectEncoder = getObjectEncoder(enc.EncoderConfig)
	return clone
}

func (enc *encoder) EncodeMessage(message string, aEnc zapcore.PrimitiveArrayEncoder) {
	aEnc.AppendString(fmt.Sprintf("%-*s", MessagePrintLength, message))
}

func (enc *encoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// We need to redact fields before we add them to the buffer, so we can't use
	if enableRedact {
		fields = enc.redact(fields)
	}

	newLine := _bufPool.Get()

	newLine = enc.encodeEntry(newLine, ent)
	// Add any structured context.
	newLine = enc.encodeContext(newLine, fields)

	newLine.AppendString(enc.lineEnding)

	newLine = enc.encodeStacktrace(newLine, ent)

	return newLine, nil
}

func (enc *encoder) encodeStacktrace(buf *buffer.Buffer, ent zapcore.Entry) *buffer.Buffer {
	// If there's no traceback key, honor that; this allows users to force single-line output.
	if len(ent.Stack) > 0 && len(enc.StacktraceKey) > 0 {
		buf.AppendString(ent.Stack)
		buf.AppendString(enc.lineEnding)
	}
	return buf
}

func (enc *encoder) encodeEntry(buf *buffer.Buffer, ent zapcore.Entry) *buffer.Buffer {
	sliceEnc := getNextSliceEncoder()
	defer putNextSliceEncoder(sliceEnc)

	if len(enc.TimeKey) > 0 && enc.EncodeTime != nil {
		enc.EncodeTime(ent.Time, sliceEnc)
	}
	if len(enc.LevelKey) > 0 && enc.EncodeLevel != nil {
		enc.EncodeLevel(ent.Level, sliceEnc)
	}
	if ent.Caller.Defined && len(enc.CallerKey) > 0 && enc.EncodeCaller != nil {
		enc.EncodeCaller(ent.Caller, sliceEnc)
	}
	if len(enc.MessageKey) > 0 {
		enc.EncodeMessage(ent.Message, sliceEnc)
	}

	return sliceEnc.flush(buf)
}

func (enc *encoder) encodeContext(buf *buffer.Buffer, extra []zapcore.Field) *buffer.Buffer {
	objEnc := enc.objectEncoder
	defer func() {
		putObjectEncoder(objEnc)
		enc.objectEncoder = getObjectEncoder(enc.EncoderConfig)
	}()

	for i := range extra {
		extra[i].AddTo(objEnc)
	}

	return objEnc.flush(buf)
}

func (enc *encoder) redact(fields []zapcore.Field) []zapcore.Field {
	if !enableRedact || len(redactKeysMap) == 0 {
		return fields
	}

	result := make([]zapcore.Field, len(fields))
	for i, field := range fields {
		// Use lowercase for case-insensitive comparison
		_, shouldRedact := redactKeysMap[strings.ToLower(field.Key)]
		if shouldRedact {
			// Create new redacted field with same key but redacted value
			result[i] = zapcore.Field{
				Key:    field.Key,
				Type:   zapcore.StringType,
				String: "[REDACTED]",
			}
		} else {
			result[i] = field
		}
	}
	return result
}
