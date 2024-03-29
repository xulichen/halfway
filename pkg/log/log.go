package log

import (
	"context"
	"fmt"
	"time"

	"github.com/xulichen/halfway/pkg/consts"
	"go.elastic.co/apm"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger = NewLogger(Development())

func GetLogger() *Loggers {
	return logger
}

type Loggers struct {
	zapLog *zap.Logger
}

// An Option configures a Logger.
type Option interface {
	apply(cfg *zap.Config)
}

type optionFunc func(cfg *zap.Config)

func (f optionFunc) apply(cfg *zap.Config) {
	f(cfg)
}

func Development() Option {
	return optionFunc(func(cfg *zap.Config) {
		cfg.OutputPaths = []string{consts.DefaultLogStdout}
		cfg.ErrorOutputPaths = []string{consts.DefaultLogStdout}
	})
}

func Production() Option {
	return optionFunc(func(cfg *zap.Config) {
		cfg.OutputPaths = []string{consts.DefaultLogPath}
		cfg.ErrorOutputPaths = []string{consts.DefaultLogPath}
	})
}

func NewLogger(opts ...Option) *Loggers {
	cfg := zap.NewProductionConfig()
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeTime = TimeEncoder
	cfg.EncoderConfig.EncodeLevel = LevelEncoder
	cfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	cfg.EncoderConfig.ConsoleSeparator = " "
	zapLog, _ := cfg.Build()
	zapLog = zapLog.WithOptions(zap.AddCallerSkip(1))
	l := new(Loggers)
	l.zapLog = zapLog
	return l
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	layout := "[2006-01-02 15:04:05]"
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}
	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}
	enc.AppendString(t.Format(layout))
}
func LevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("[%s]", l.CapitalString()))
}

func (l *Loggers) GetZapLog() *zap.Logger {
	return l.zapLog
}

func (l *Loggers) Print(i ...interface{}) {
	l.zapLog.Info(fmt.Sprint(i...))
}
func (l *Loggers) Printf(format string, args ...interface{}) {
	l.zapLog.Info(fmt.Sprintf(format, args...))
}

func (l *Loggers) Println(args ...interface{}) {
	l.zapLog.Info(fmt.Sprint(args...))
}

func (l *Loggers) Debug(i ...interface{}) {
	l.zapLog.Debug(fmt.Sprint(i...))
}
func (l *Loggers) Debugf(format string, args ...interface{}) {
	l.zapLog.Debug(fmt.Sprintf(format, args...))
}

func (l *Loggers) Info(i ...interface{}) {
	l.zapLog.Info(fmt.Sprint(i...))
}
func (l *Loggers) Infof(format string, args ...interface{}) {
	l.zapLog.Info(fmt.Sprintf(format, args...))
}

func (l *Loggers) Warn(i ...interface{}) {
	l.zapLog.Warn(fmt.Sprint(i...))
}
func (l *Loggers) Warnf(format string, args ...interface{}) {
	l.zapLog.Warn(fmt.Sprintf(format, args...))
}

func (l *Loggers) Error(i ...interface{}) {
	l.zapLog.Error(fmt.Sprint(i...))
}

func (l *Loggers) Errorf(format string, args ...interface{}) {
	l.zapLog.Error(fmt.Sprintf(format, args...))
}

func (l *Loggers) Fatal(i ...interface{}) {
	l.zapLog.Fatal(fmt.Sprint(i...))
}

func (l *Loggers) Fatalf(format string, args ...interface{}) {
	l.zapLog.Error(fmt.Sprintf(format, args...))
}

func (l *Loggers) Fatalln(i ...interface{}) {
	l.zapLog.Fatal(fmt.Sprint(i...))
}

func (l *Loggers) Panic(i ...interface{}) {
	l.zapLog.Panic(fmt.Sprint(i...))
}

func (l *Loggers) Panicf(format string, args ...interface{}) {
	l.zapLog.Panic(fmt.Sprintf(format, args...))
}

func (l *Loggers) Panicln(i ...interface{}) {
	l.zapLog.Panic(fmt.Sprint(i...))
}

func (l *Loggers) WithContext(ctx context.Context) *Loggers {
	cp := *l
	tx := apm.TransactionFromContext(ctx)
	traceContext := tx.TraceContext()
	span := apm.SpanFromContext(ctx)
	if span != nil {
		spanId := span.TraceContext().Span
		cp.zapLog = cp.zapLog.With(zap.String("trace.id", traceContext.Trace.String()),
			zap.String("transaction.id", traceContext.Trace.String()),
			zap.String("span.id", spanId.String()))
	} else {
		cp.zapLog = cp.zapLog.With(zap.String("trace.id", traceContext.Trace.String()),
			zap.String("transaction.id", traceContext.Trace.String()),
			zap.String("span.id", ""))
	}
	return &cp
}

// InjectCtx 日志注入上下文
func InjectCtx(l *zap.Logger, ctx context.Context) *zap.Logger {
	tx := apm.TransactionFromContext(ctx)
	traceContext := tx.TraceContext()
	span := apm.SpanFromContext(ctx)
	if span != nil {
		spanId := span.TraceContext().Span
		l = l.With(zap.String("trace.id", traceContext.Trace.String()),
			zap.String("transaction.id", traceContext.Trace.String()),
			zap.String("span.id", spanId.String()))
	} else {
		l = l.With(zap.String("trace.id", traceContext.Trace.String()),
			zap.String("transaction.id", traceContext.Trace.String()),
			zap.String("span.id", ""))
	}
	return l
}
