package mlog

import (
	"github.com/ThreeDotsLabs/watermill"
	"go.uber.org/zap"
)

type WatermillLog struct {
	backend *zap.Logger
	fields  watermill.LogFields
}

func InitWatermillLog() watermill.LoggerAdapter {
	return &WatermillLog{backend: logger}
}

func (l *WatermillLog) Error(msg string, err error, fields watermill.LogFields) {
	fields = l.fields.Add(fields)
	fs := make([]zap.Field, 0, len(fields)+1)
	fs = append(fs, zap.Error(err))
	for k, v := range fields {
		fs = append(fs, zap.Any(k, v))
	}
	l.backend.Error(msg, fs...)
}

func (l *WatermillLog) Info(msg string, fields watermill.LogFields) {
	fields = l.fields.Add(fields)
	fs := make([]zap.Field, 0, len(fields)+1)
	for k, v := range fields {
		fs = append(fs, zap.Any(k, v))
	}
	l.backend.Info(msg, fs...)
}

func (l *WatermillLog) Debug(msg string, fields watermill.LogFields) {
	fields = l.fields.Add(fields)
	fs := make([]zap.Field, 0, len(fields)+1)
	for k, v := range fields {
		fs = append(fs, zap.Any(k, v))
	}
	l.backend.Debug(msg, fs...)
}

func (l *WatermillLog) Trace(msg string, fields watermill.LogFields) {
	fields = l.fields.Add(fields)
	fs := make([]zap.Field, 0, len(fields)+1)
	for k, v := range fields {
		fs = append(fs, zap.Any(k, v))
	}
	l.backend.Debug(msg, fs...)
}

func (l *WatermillLog) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &WatermillLog{
		backend: l.backend,
		fields:  l.fields.Add(fields),
	}
}
