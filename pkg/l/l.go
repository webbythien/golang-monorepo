package l

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/k0kubun/pp/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/webbythien/monorepo/pkg/prettylog"
)

// ConsoleEncoderName ...
const (
	ConsoleEncoderName = "pp"
	LokiEncoderName    = "loki"
)

var _ = func() error {
	err := zap.RegisterEncoder(ConsoleEncoderName, func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return newEncoder(cfg), nil
	})
	if err != nil {
		slog.Error("failed to register custom console encoder", slog.Any("err", err))
	}
	return err
}()

var envPatterns []*regexp.Regexp

type config struct {
	LogLevel              string `yaml:"log_level" mapstructure:"log_level"`                               // Default log level
	DebugEnabler          string `yaml:"debug_enabler" mapstructure:"debug_enabler"`                       // Comma separated list of debug enabler
	ConsoleObjectMaxDepth int    `yaml:"console_object_max_depth" mapstructure:"console_object_max_depth"` // Max depth for console object
}

func loadConfig() *config {
	cfg := defaultConfig
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
	if debugEnabler := os.Getenv("DEBUG_ENABLER"); debugEnabler != "" {
		cfg.DebugEnabler = debugEnabler
	}
	if consoleObjectMaxDepth := os.Getenv("CONSOLE_OBJECT_MAX_DEPTH"); consoleObjectMaxDepth != "" {
		consoleObjectMaxDepthInt, err := strconv.Atoi(consoleObjectMaxDepth)
		if err != nil {
			slog.Error("failed to parse CONSOLE_OBJECT_MAX_DEPTH", slog.Any("err", err))
		} else {
			cfg.ConsoleObjectMaxDepth = consoleObjectMaxDepthInt
		}
	}
	return &cfg
}

func (cfg *config) getSlogLevel() slog.Level {
	switch cfg.LogLevel {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (cfg *config) getZapLevel() zapcore.Level {
	var lv zapcore.Level
	err := lv.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		panic(err)
	}
	return lv
}

var defaultConfig = config{
	LogLevel:              "INFO",
	DebugEnabler:          "",
	ConsoleObjectMaxDepth: 3,
}

var objPrinter = pp.New()

var initConfigOnce sync.Once

var sharedConfig *config

func initConfig() {
	cfg := loadConfig()
	// Set default max depth for prettylog
	pp.SetDefaultMaxDepth(cfg.ConsoleObjectMaxDepth)
	objPrinter.SetExportedOnly(true)
	if os.Getenv("LOG_BLIND_MODE") == "true" {
		objPrinter.SetColoringEnabled(false)
		DefaultConsoleEncoderConfig.EncodeLevel = EncodeLevel
		DefaultConsoleEncoderConfig.EncodeTime = PrettyNanoTimeEncoder
	}

	switch os.Getenv("LOG_ENCODER") {
	case LokiEncoderName:
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level:     cfg.getSlogLevel(),
			AddSource: false,
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if enableRedact {
					// Use lowercase for case-insensitive comparison
					if _, shouldRedact := redactKeysMap[strings.ToLower(a.Key)]; shouldRedact {
						a.Value = slog.StringValue("[REDACTED]")
					}
				}
				switch a.Key {
				case slog.LevelKey:
					a.Key = "level"
				case slog.TimeKey:
					a.Key = "timestamp"
				case slog.MessageKey:
					a.Key = "message"
				case slog.SourceKey:
					a.Key = "caller"
				}
				return a
			},
		})))
		slog.Info("Enable log at level", "LOG_LEVEL", cfg.getSlogLevel())
	case ConsoleEncoderName:
		fallthrough
	default:
		// Set default slog logger
		slog.SetDefault(slog.New(prettylog.NewHandler(&slog.HandlerOptions{
			Level:     cfg.getSlogLevel(),
			AddSource: false,
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if enableRedact {
					// Use lowercase for case-insensitive comparison
					if _, shouldRedact := redactKeysMap[strings.ToLower(a.Key)]; shouldRedact {
						a.Value = slog.StringValue("[REDACTED]")
					}
				}
				return a
			},
		})))
		// slog.Info("Enable log at level", "LOG_LEVEL", cfg.getSlogLevel())
	}

	enablers = make(map[string]zap.AtomicLevel)
	enablers[""] = zap.NewAtomicLevel()

	// Set default zap logger
	var lv zapcore.Level
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		lv = zapcore.DebugLevel
	case "INFO":
		lv = zapcore.InfoLevel
	case "WARN":
		lv = zapcore.WarnLevel
	case "ERROR":
		lv = zapcore.ErrorLevel
	default:
		lv = cfg.getZapLevel()
	}
	for _, enabler := range enablers {
		enabler.SetLevel(lv)
	}

	debugEnabler := os.Getenv("DEBUG_ENABLER")
	var errPattern string
	envPatterns, errPattern = initPatterns(debugEnabler)
	if errPattern != "" {
		slog.Error("Unable to parse LOG_LEVEL. Please set it to a proper value.", slog.String("invalid", errPattern))
		os.Exit(1)
	}

	sharedConfig = cfg
}

// Logger wraps zap.Logger
type Logger struct {
	*zap.Logger
	S *zap.SugaredLogger
}

// Shorthanded functions for logging.
var (
	Any        = zap.Any
	Bool       = zap.Bool
	Duration   = zap.Duration
	Float64    = zap.Float64
	Int        = zap.Int
	Int64      = zap.Int64
	Skip       = zap.Skip
	String     = zap.String
	Stringer   = zap.Stringer
	Time       = zap.Time
	Uint       = zap.Uint
	Uint32     = zap.Uint32
	Uint64     = zap.Uint64
	Uintptr    = zap.Uintptr
	ByteString = zap.ByteString
)

// DefaultConsoleEncoderConfig ...
var DefaultConsoleEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	EncodeTime:     ConsoleTimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   ShortColorCallerEncoder,
}

func ConsoleTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	layout := prettylog.LogTimeFormat
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}

	enc.AppendString(t.Format(layout))
}

var LokiEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "timestamp",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    EncodeLevel,
	EncodeTime:     RFC3339NanoTimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

var logLevelSeverity = map[zapcore.Level]string{
	zapcore.DebugLevel:  "DEBUG",
	zapcore.InfoLevel:   "INFO",
	zapcore.WarnLevel:   "WARNING",
	zapcore.ErrorLevel:  "ERROR",
	zapcore.DPanicLevel: "CRITICAL",
	zapcore.PanicLevel:  "ALERT",
	zapcore.FatalLevel:  "EMERGENCY",
}

// EncodeLevel maps the internal Zap log level to the appropriate Stackdriver
// level.
func EncodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(logLevelSeverity[l])
}

// RFC3339NanoTimeEncoder serializes a time.Time to an RFC3339Nano-formatted
// string with nanoseconds precision.
func RFC3339NanoTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.RFC3339Nano))
}

func PrettyNanoTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(prettylog.LogTimeFormat))
}

// Error wraps error for zap.Error.
func Error(err error) zapcore.Field {
	if err == nil {
		return Skip()
	}
	return String("error", err.Error())
}

// Interface ...
func Interface(key string, val interface{}) zapcore.Field {
	if val, ok := val.(fmt.Stringer); ok {
		return zap.Stringer(key, val)
	}
	return zap.Reflect(key, val)
}

// Stack ...
func Stack() zapcore.Field {
	return zap.Stack("stack")
}

// Int32 ...
func Int32(key string, val int32) zapcore.Field {
	return zap.Int(key, int(val))
}

// Object ...
var Object = zap.Any

type dd struct {
	v interface{}
}

func (d dd) String() string {
	return objPrinter.Sprint(d.v)
}

// Dump renders object for debugging
func Dump(v interface{}) fmt.Stringer {
	return dd{v}
}

const (
	CallerPrintLength  = 30
	MessagePrintLength = 60
)

// ShortColorCallerEncoder encodes caller information with sort path filename and enable color.
func ShortColorCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%-*s", CallerPrintLength, prettylog.ShortenCallerPath(caller.TrimmedPath())))
}

func NewWithName(name string, opts ...zap.Option) Logger {
	return newWithName(name, opts...)
}

func newWithName(name string, opts ...zap.Option) Logger {
	if name == "" {
		_, filename, _, _ := runtime.Caller(1)
		name = filepath.Dir(truncFilename(filename))
	}

	var enabler zap.AtomicLevel
	if e, ok := enablers[name]; ok {
		enabler = e
	} else {
		enabler = zap.NewAtomicLevel()
		enablers[name] = enabler
	}

	setLogLevelFromEnv(name, enabler)

	var loggerConfig zap.Config
	switch os.Getenv("LOG_ENCODER") {
	case LokiEncoderName:
		loggerConfig = zap.Config{
			Level:       enabler,
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "json",
			EncoderConfig:    LokiEncoderConfig,
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
	case ConsoleEncoderName:
		fallthrough
	default:
		loggerConfig = zap.Config{
			Level:            enabler,
			Development:      false,
			Encoding:         ConsoleEncoderName, // "json",
			EncoderConfig:    DefaultConsoleEncoderConfig,
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}
	stacktraceLevel := zap.NewAtomicLevelAt(zapcore.PanicLevel)

	opts = append(opts, zap.AddStacktrace(stacktraceLevel))
	logger, err := loggerConfig.Build(opts...)
	if err != nil {
		panic(err)
	}
	return Logger{
		Logger: logger,
		S:      logger.WithOptions(zap.AddCallerSkip(1)).Sugar(),
	}
}

// New returns new zap.Logger
func New(opts ...zap.Option) Logger {
	initConfigOnce.Do(initConfig)
	return newWithName("", opts...)
}

func (logger Logger) AddCallerSkip(skip int) Logger {
	zl := logger.Logger.WithOptions(zap.AddCallerSkip(skip))
	return Logger{
		Logger: zl,
		S:      zl.Sugar(),
	}
}

// ServeHTTP supports logging level with an HTTP request.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	type payload struct {
		Name  string `json:"name"`
		Level string `json:"level"`
	}

	enc := json.NewEncoder(w)

	switch r.Method {
	case http.MethodGet:
		var payloads []payload
		for k, e := range enablers {
			lvl := e.Level()
			payloads = append(payloads, payload{
				Name:  k,
				Level: lvl.String(),
			})
		}
		err := enc.Encode(payloads)
		if err != nil {
			panic(err)
		}

	case http.MethodPut:
		var req payload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = enc.Encode(errorResponse{
				Error: fmt.Sprintf("Request body must be valid JSON: %v", err),
			})
			if err != nil {
				panic(err)
			}
			return
		}

		if req.Level == "" {
			w.WriteHeader(http.StatusBadRequest)
			err := enc.Encode(errorResponse{
				Error: errLevelNil.Error(),
			})
			if err != nil {
				panic(err)
			}
			return
		}

		var lv zapcore.Level
		err := lv.UnmarshalText([]byte(req.Level))
		if err != nil {
			panic(err)
		}

		if req.Name == "" {
			for _, enabler := range enablers {
				enabler.SetLevel(lv)
			}
		} else {
			enabler, ok := enablers[req.Name]
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				err := enc.Encode(errorResponse{
					Error: errEnablerNotFound.Error(),
				})
				if err != nil {
					panic(err)
				}
				return
			}

			enabler.SetLevel(lv)
		}

		err = enc.Encode(req)
		if err != nil {
			panic(err)
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		err := enc.Encode(errorResponse{
			Error: "Only GET and PUT are supported.",
		})
		if err != nil {
			panic(err)
		}
	}
}

var (
	errEnablerNotFound = errors.New("enabler not found")
	errLevelNil        = errors.New("must specify a logging level")

	enablers = make(map[string]zap.AtomicLevel)
)

func truncFilename(filename string) string {
	// index := strings.Index(filename, prefix)
	// return filename[index+len(prefix):]
	return filename
}

func initPatterns(envLog string) ([]*regexp.Regexp, string) {
	patterns := strings.Split(envLog, ",")
	result := make([]*regexp.Regexp, len(patterns))
	for i, p := range patterns {
		r, err := parsePattern(p)
		if err != nil {
			return nil, p
		}

		result[i] = r
	}
	return result, ""
}

func parsePattern(p string) (*regexp.Regexp, error) {
	p = strings.ReplaceAll(strings.Trim(p, " "), "*", ".*")
	return regexp.Compile(p)
}

func setLogLevelFromEnv(name string, enabler zap.AtomicLevel) {
	var zapLevel zapcore.Level = zap.DebugLevel
	if sharedConfig != nil {
		zapLevel = sharedConfig.getZapLevel()
	}
	for _, r := range envPatterns {
		if r.MatchString(name) {
			enabler.SetLevel(zapLevel)
		}
	}
}
