package middleware

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/color"
	"github.com/valyala/fasttemplate"
	"github.com/xulichen/halfway/pkg/log"
	"go.uber.org/zap"
)

type (
	// LoggerConfig defines the config for Logger middleware.
	LoggerConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echoMiddleware.Skipper

		// Output is a writer where logs in JSON format are written.
		// Optional. Default value os.Stdout.
		Output io.Writer

		template *fasttemplate.Template
		colorer  *color.Color
		pool     *sync.Pool

		Handler BodyDumpHandler
	}

	BodyDumpHandler func(echo.Context, time.Time, []byte, []byte)

	bodyDumpResponseWriter struct {
		io.Writer
		http.ResponseWriter
	}
)

var (
	// DefaultLoggerConfig is the default Logger middleware config.
	DefaultLoggerConfig = LoggerConfig{
		Skipper: func(context echo.Context) bool {
			if strings.Index(context.Request().URL.Path, "/swagger") != -1 {
				return true
			}
			return false
		},
		colorer: color.New(),
	}
)

// restapi监控日志
func Logger() echo.MiddlewareFunc {
	return LoggerWithConfig(LoggerConfig{Handler: func(context echo.Context, start time.Time, i []byte, i2 []byte) {
		path := context.Request().URL.Path
		if context.Request().URL.Query().Encode() != "" {
			path += "?" + context.Request().URL.Query().Encode()
		}

		l := log.GetLogger().GetZapLog().With(zap.String("time", time.Now().Format("2006-01-02 15:04:05")),
			zap.Int("status_code", context.Response().Status),
			zap.String("method", context.Request().Method),
			zap.String("path", path),
			zap.String("latency", time.Now().Sub(start).String()),
			zap.String("clientIP", context.RealIP()),
			zap.String("authorization", context.Request().Header.Get("Authorization")))

		if len(i) == 0 {
			i = []byte("{}")
		} else {
			i = bytes.Replace(i, []byte("\r\n"), []byte(""), -1)
			i = bytes.Replace(i, []byte(" "), []byte(""), -1)
		}
		if len(i2) == 0 {
			i2 = []byte("{}")
		}
		// apm链路追踪
		ctx := context.Request().Context()
		l = log.InjectCtx(l, ctx)
		l.Info("", zap.Reflect("request", json.RawMessage(i)), zap.Reflect("response", json.RawMessage(i2)))
	}})
}

// LoggerWithConfig returns a Logger middleware with config.
// See: `Logger()`.
func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultLoggerConfig.Skipper
	}
	if config.Output == nil {
		config.Output = DefaultLoggerConfig.Output
	}

	config.colorer = color.New()
	config.colorer.SetOutput(config.Output)
	config.pool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 256))
		},
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			start := time.Now()

			if config.Skipper(c) {
				return next(c)
			}

			// Request
			var reqBody []byte
			if c.Request().Body != nil { // Read
				reqBody, _ = ioutil.ReadAll(c.Request().Body)
			}
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

			// Response
			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			if err = next(c); err != nil {
				c.Error(err)
			}
			// Callback
			config.Handler(c, start, reqBody, resBody.Bytes())

			return
		}
	}
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
