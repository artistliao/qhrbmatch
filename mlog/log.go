package mlog

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

var (
	atomicLevel = zap.NewAtomicLevel()
	logger      *zap.Logger
)

type Level = zapcore.Level

const (
	DebugLevel int = iota - 1 // -1
	// InfoLevel is the default logging priority.
	InfoLevel // 0
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel // 1
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel // 2
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel // 3
	// PanicLevel logs a message, then panics.
	PanicLevel // 4
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel // 5
)

const (
	//ALL, TRACE, DEBUG int = -1, -1, -1
	TRACE = iota - 1
	INFO
	WARN
	ERROR
	FATAL
	DEBUG = TRACE
	ALL   = TRACE
)

func init() {
	logger = zap.NewExample()
}

// Detach
// @Description: 自定义context, 可实现将重要信息向下透传
type Detach struct {
	ctx context.Context
}

func (d Detach) Deadline() (time.Time, bool) {
	return time.Time{}, false
}
func (d Detach) Done() <-chan struct{} {
	return nil
}
func (d Detach) Err() error {
	return nil
}
func (d Detach) Value(key interface{}) interface{} {
	return d.ctx.Value(key)
}

func getCallerName(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	f := runtime.FuncForPC(pc)
	fname := f.Name()
	ind := strings.LastIndex(f.Name(), "/")
	if ind >= 0 {
		fname = fname[ind+1:]
	}

	return fname
	//fname = fmt.Sprintf("(%s) ", fname)
}
func Trace(v ...interface{}) {
	if !logger.Core().Enabled(zap.DebugLevel) {
		return
	}

	//logger.Debug(fmt.Sprint(v...))
	logger.Debug(fmt.Sprintf("%s", v...), zap.String("fname", getCallerName(2)))
}
func Tracef(format string, v ...interface{}) {
	if !logger.Core().Enabled(zap.DebugLevel) {
		return
	}
	/*
		pc, _, _, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		fname := f.Name()
		ind := strings.LastIndex(f.Name(), "/")
		if ind >= 0 {
			fname = fname[ind+1:]
		}
		fname = fmt.Sprintf("(%s) ", fname)
	*/
	logger.Debug(fmt.Sprintf("%s", v...), zap.String("fname", getCallerName(2)))
}
func CTracef(skipCaller int, format string, v ...interface{}) {
	if !logger.Core().Enabled(zap.DebugLevel) {
		return
	}
	/*
		pc, _, _, _ := runtime.Caller(skipCaller)
		f := runtime.FuncForPC(pc)
		fname := f.Name()
		ind := strings.LastIndex(f.Name(), "/")
		if ind >= 0 {
			fname = fname[ind+1:]
		}
		fname = fmt.Sprintf("(%s) ", fname)
	*/
	fname := fmt.Sprintf("%s", getCallerName(skipCaller+1))
	logger.Debug(fmt.Sprintf(format, v...), zap.String("fname", fname))
}

func Debug(v ...interface{}) {
	if !logger.Core().Enabled(zap.DebugLevel) {
		return
	}
	//logger.Debug(fmt.Sprint(v...))
	logger.Debug(fmt.Sprintf("%s", v...), zap.String("fname", getCallerName(2)))
}
func Debugf(format string, v ...interface{}) {
	if !logger.Core().Enabled(zap.DebugLevel) {
		return
	}
	//logger.Debug(fmt.Sprintf(format, v...))
	logger.Debug(fmt.Sprintf(format, v...), zap.String("fname", getCallerName(2)))
}
func Info(v ...interface{}) {
	logger.Info(fmt.Sprint(v...))
}
func Infof(format string, v ...interface{}) {
	logger.Info(fmt.Sprintf(format, v...))
}
func Warn(v ...interface{}) {
	logger.Warn(fmt.Sprint(v...))
}
func Warnf(format string, v ...interface{}) {
	logger.Warn(fmt.Sprintf(format, v...))
}
func Error(v ...interface{}) {
	logger.Error(fmt.Sprint(v...))
}
func Errorf(format string, v ...interface{}) {
	logger.Error(fmt.Sprintf(format, v...))
}
func Fatal(v ...interface{}) {
	logger.Panic(fmt.Sprint(v...))
}
func Fatalf(format string, v ...interface{}) {
	log.Panic(fmt.Sprintf(format, v...))
}

type Field = zap.Field

func ZDebug(msg string, fields ...Field) {
	logger.Debug(msg, append(fields, zap.String("fname", getCallerName(2)))...)
}

func ZInfo(msg string, fields ...Field) {
	logger.Info(msg, fields...)
}

func ZWarn(msg string, fields ...Field) {
	logger.Warn(msg, fields...)
}

func ZError(msg string, fields ...Field) {
	logger.Error(msg, fields...)
}
func ZDPanic(msg string, fields ...Field) {
	logger.DPanic(msg, fields...)
}
func ZPanic(msg string, fields ...Field) {
	logger.Panic(msg, fields...)
}
func ZFatal(msg string, fields ...Field) {
	logger.Fatal(msg, fields...)
}

func DebugCtx(ctx context.Context, msg string, fields ...Field) {
	spanCtx := trace.SpanContextFromContext(ctx)
	logger.Debug(msg, append(fields, zap.String("traceId", spanCtx.TraceID().String()), zap.String("fname", getCallerName(2)))...)
}

func InfoCtx(ctx context.Context, msg string, fields ...Field) {
	spanCtx := trace.SpanContextFromContext(ctx)
	logger.Info(msg, append(fields, zap.String("traceId", spanCtx.TraceID().String()))...)
}

func WarnCtx(ctx context.Context, msg string, fields ...Field) {
	spanCtx := trace.SpanContextFromContext(ctx)
	logger.Warn(msg, append(fields, zap.String("traceId", spanCtx.TraceID().String()))...)
}

func ErrorCtx(ctx context.Context, msg string, fields ...Field) {
	spanCtx := trace.SpanContextFromContext(ctx)
	logger.Error(msg, append(fields, zap.String("traceId", spanCtx.TraceID().String()))...)
}
func DPanicCtx(ctx context.Context, msg string, fields ...Field) {
	spanCtx := trace.SpanContextFromContext(ctx)
	logger.DPanic(msg, append(fields, zap.String("traceId", spanCtx.TraceID().String()))...)
}
func PanicCtx(ctx context.Context, msg string, fields ...Field) {
	spanCtx := trace.SpanContextFromContext(ctx)
	logger.Panic(msg, append(fields, zap.String("traceId", spanCtx.TraceID().String()))...)
}
func FatalCtx(ctx context.Context, msg string, fields ...Field) {
	spanCtx := trace.SpanContextFromContext(ctx)
	logger.Fatal(msg, append(fields, zap.String("traceId", spanCtx.TraceID().String()))...)
}

func TracefCtx(ctx context.Context, format string, v ...interface{}) {
	spanCtx := trace.SpanContextFromContext(ctx)
	logger.Debug(fmt.Sprintf(format, v...), zap.String("traceId", spanCtx.TraceID().String()), zap.String("fname", getCallerName(2)))
}
func DebugfCtx(ctx context.Context, format string, v ...interface{}) {
	spanCtx := trace.SpanContextFromContext(ctx)
	logger.Debug(fmt.Sprintf(format, v...), zap.String("traceId", spanCtx.TraceID().String()), zap.String("fname", getCallerName(2)))
}
func ErrorfCtx(ctx context.Context, format string, v ...interface{}) {
	spanCtx := trace.SpanContextFromContext(ctx)
	logger.Error(fmt.Sprintf(format, v...), zap.String("traceId", spanCtx.TraceID().String()))
}

func WarnfCtx(ctx context.Context, format string, v ...interface{}) {
	spanCtx := trace.SpanContextFromContext(ctx)
	logger.Warn(fmt.Sprintf(format, v...), zap.String("traceId", spanCtx.TraceID().String()))
}

func InfofCtx(ctx context.Context, format string, v ...interface{}) {
	spanCtx := trace.SpanContextFromContext(ctx)
	logger.Info(fmt.Sprintf(format, v...), zap.String("traceId", spanCtx.TraceID().String()))
}

func SetLevel(level int) {
	atomicLevel.SetLevel(zapcore.Level(level))
}

type Params struct {
	Path       string //	路径
	MaxSize    int    //	MB
	MaxBackups int    //	备份个数
	MaxAge     int    //	保存时间,天
	Level      int
}

func InitLogger(params *Params) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("2006-01-02T15:04:05.000Z0700"))
	}
	if params.MaxBackups <= 0 {
		params.MaxBackups = 3
	}
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   params.Path,
		MaxSize:    params.MaxSize,
		MaxBackups: params.MaxBackups,
		MaxAge:     params.MaxAge,
		Compress:   true,
	})
	SetLevel(params.Level)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.NewMultiWriteSyncer(fileWriter, zapcore.AddSync(os.Stdout)),
		atomicLevel,
	)
	SetDefaultLogger(zap.New(core, zap.WithCaller(true), zap.AddCallerSkip(1)))
}

func SetDefaultLogger(lg *zap.Logger) {
	logger = lg
}

func SetDebugLogger(t *testing.T) {
	SetDefaultLogger(zaptest.NewLogger(t))
}

func GetLevel() int {
	return int(atomicLevel.Level())
}

func GetLogger() *zap.Logger {
	return logger
}
