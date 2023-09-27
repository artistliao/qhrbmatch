package mlog

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"testing"
)

/**
* @ Author: tly
* @ Date: 2023/2/22 15:58
* @ Desc:
 */

func TestInitLogger(t *testing.T) {
	type User struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	InitLogger(&Params{
		Path: "./access.log",
	})

	Info("before set level info .....")
	atomicLevel.SetLevel(-1)
	Debug("debug info", zap.String("password", "123456"))
	Info("user", zap.Any("user", User{Name: "tly", Password: "123456"}))
	Info("info log")
}

func TestLoggerCtx(t *testing.T) {
	InitLogger(&Params{
		Path: "./access.log",
	})
	ctx := context.Background()
	InfoCtx(ctx, "before ctx  info .....")
	tp := trace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	spanCtx, span := otel.Tracer("test").Start(ctx, "test")
	InfoCtx(spanCtx, "after ctx info .....")
	span.End()
}
