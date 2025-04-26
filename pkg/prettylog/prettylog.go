package prettylog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/k0kubun/pp/v3"
)

const (
	LogTimeFormat = "06-01-02 15:04:05.000Z07" // RFC3339Milis

	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97

	// ppScale = pp.Black | pp.BackgroundWhite - black
)

func colorizer(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

type Handler struct {
	h        slog.Handler
	r        func([]string, slog.Attr) slog.Attr
	b        *bytes.Buffer
	m        *sync.Mutex
	writer   io.Writer
	colorize bool
	printer  *pp.PrettyPrinter
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), b: h.b, r: h.r, m: h.m, writer: h.writer, colorize: h.colorize, printer: h.printer}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name), b: h.b, r: h.r, m: h.m, writer: h.writer, colorize: h.colorize, printer: h.printer}
}

// TODO: revamp this later, don't use in production
func (h *Handler) computeAttrs(
	ctx context.Context,
	r slog.Record,
) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()
	if err := h.h.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}

	if _, ok := attrs[CallerKey]; !ok {
		_, file, line, _ := runtime.Caller(4)
		// get last 2 parts of the file path
		parts := strings.Split(file, "/")
		file = strings.Join(parts[len(parts)-2:], "/")
		attrs[CallerKey] = fmt.Sprintf("%s:%d", file, line)
	}
	return attrs, nil
}

const CallerKey = "caller"

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	colorize := func(code int, value string) string {
		return value
	}
	if h.colorize && os.Getenv("LOG_BLIND_MODE") != "true" {
		colorize = colorizer
	}

	var level string
	levelAttr := slog.Attr{
		Key:   slog.LevelKey,
		Value: slog.AnyValue(r.Level),
	}
	if h.r != nil {
		levelAttr = h.r([]string{}, levelAttr)
	}

	if !levelAttr.Equal(slog.Attr{}) {
		switch {
		case r.Level <= slog.LevelDebug:
			level = colorize(blue, "DBG")
		case r.Level <= slog.LevelInfo:
			level = colorize(green, "INF")
		case r.Level <= slog.LevelWarn:
			level = colorize(yellow, "WRN")
		case r.Level <= slog.LevelError:
			level = colorize(red, "ERR")
		default:
			level = colorize(red, "???")
		}
	}

	var timestamp string
	timeAttr := slog.Attr{
		Key:   slog.TimeKey,
		Value: slog.StringValue(r.Time.Format(LogTimeFormat)),
	}
	if h.r != nil {
		timeAttr = h.r([]string{}, timeAttr)
	}
	if !timeAttr.Equal(slog.Attr{}) {
		timestamp = colorize(lightGray, timeAttr.Value.String())
	}

	var msg string
	msgAttr := slog.Attr{
		Key:   slog.MessageKey,
		Value: slog.StringValue(r.Message),
	}
	if h.r != nil {
		msgAttr = h.r([]string{}, msgAttr)
	}
	if !msgAttr.Equal(slog.Attr{}) {
		msg = colorize(darkGray, msgAttr.Value.String())
	}

	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	const delim = "\t"
	out := strings.Builder{}
	if len(timestamp) > 0 {
		out.WriteString(timestamp)
		out.WriteString(delim)
	}
	if len(level) > 0 {
		out.WriteString(level)
		out.WriteString(delim)
	}
	if caller, ok := attrs[CallerKey]; ok {
		callerStr, _ := caller.(string)
		callerValue := ShortenCallerPath(callerStr)
		callerValue = colorize(lightGray, fmt.Sprintf("%-*s", CallerPrintLength, callerValue))
		out.WriteString(callerValue)
		out.WriteString(delim)
		delete(attrs, CallerKey)
	}
	if len(msg) > 0 {
		_, _ = fmt.Fprintf(&out, "%-*s", MessagePrintLength, msg)
		out.WriteString(delim)
		out.WriteString(delim)
	}
	for key, val := range attrs {
		out.WriteString(delim)
		out.WriteString(colorize(lightGray, key))
		out.WriteRune('=')
		_, _ = h.printer.Fprint(&out, val)
	}

	_, err = io.WriteString(h.writer, out.String()+"\n")
	if err != nil {
		return err
	}

	return nil
}

func suppressDefaults(
	next func([]string, slog.Attr) slog.Attr,
) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

var printer = configPrinter()

func Print(args ...interface{}) (n int, err error) {
	if printer == nil {
		return pp.Print(args...)
	} else {
		return printer.Print(args...)
	}
}

func SetPrinter(altPrinter *pp.PrettyPrinter) {
	sync.OnceFunc(func() {
		printer = altPrinter
	})
}

func init() {
	SetPrinter(configPrinter())
}

func configPrinter() *pp.PrettyPrinter {
	printer := pp.New()
	printer.SetColorScheme(pp.ColorScheme{
		Bool:            pp.Cyan,
		Integer:         pp.Blue,
		Float:           pp.Magenta,
		String:          pp.Red,
		StringQuotation: pp.Red,
		EscapedChar:     pp.Magenta,
		FieldName:       pp.Yellow,
		PointerAdress:   pp.Blue,
		Nil:             pp.Cyan,
		Time:            pp.Blue,
		StructName:      pp.Green,
		ObjectLength:    pp.Blue,
	})
	printer.SetExportedOnly(true)
	if os.Getenv("LOG_BLIND_MODE") == "true" {
		printer.SetColoringEnabled(false)
	}
	return printer
}
func NewHandlerWithOptions(handlerOptions *slog.HandlerOptions, options ...Option) *Handler {
	if handlerOptions == nil {
		handlerOptions = &slog.HandlerOptions{}
	}

	buf := &bytes.Buffer{}

	printer := configPrinter()

	handler := &Handler{
		b: buf,
		h: slog.NewJSONHandler(buf, &slog.HandlerOptions{
			Level:       handlerOptions.Level,
			AddSource:   handlerOptions.AddSource,
			ReplaceAttr: suppressDefaults(handlerOptions.ReplaceAttr),
		}),
		r:       handlerOptions.ReplaceAttr,
		m:       &sync.Mutex{},
		printer: printer,
	}

	for _, opt := range options {
		opt(handler)
	}

	return handler
}

func NewHandler(opts *slog.HandlerOptions) *Handler {
	return NewHandlerWithOptions(opts, WithDestinationWriter(os.Stdout), WithColor())
}

type Option func(h *Handler)

func WithDestinationWriter(writer io.Writer) Option {
	return func(h *Handler) {
		h.writer = writer
	}
}

func WithColor() Option {
	return func(h *Handler) {
		h.colorize = true
	}
}

// ShortenCallerPath shortens the caller path if the callerPath is longer than CallerPrintLength.
// keeping the last CallerPrintLength-3 characters.
// masking the first as 3 asterisks
func ShortenCallerPath(path string) string {
	if len(path) <= CallerPrintLength {
		return path
	}
	return "***" + path[len(path)-CallerPrintLength+3:]
}

const (
	CallerPrintLength  = 30
	MessagePrintLength = 60
)
