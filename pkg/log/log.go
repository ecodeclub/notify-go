// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/pborman/uuid"
)

type (
	ContextLogKey struct{}
	LogIDKey      struct{}
)

func FromContext(ctx context.Context) *Logger {
	if l, ok := ctx.Value(ContextLogKey{}).(*Logger); ok {
		return l
	}
	return Default()
}

type Logger struct {
	*slog.Logger
}

func Default() *Logger {
	l := newLogger(os.Stdout, slog.LevelInfo)
	l.Logger = l.Logger.With("LOGID", uuid.NewUUID().String())
	return l
}

func New() *Logger {
	return newLogger(os.Stdout, slog.LevelInfo)
}

func newLogger(w io.Writer, level slog.Level) *Logger {
	l := &Logger{
		Logger: slog.New(
			slog.NewTextHandler(w, &slog.HandlerOptions{
				AddSource: false,
				Level:     level,
			}),
		),
	}
	return l
}

func (l *Logger) WithContext(ctx context.Context) context.Context {
	if _, ok := ctx.Value(ContextLogKey{}).(*Logger); ok {
		return ctx
	}
	return context.WithValue(ctx, ContextLogKey{}, l)
}

func (l *Logger) Auto(msg string, err error, args ...any) {
	if err != nil {
		l.Error(msg, append(args, "err", err.Error())...)
		return
	}
	l.Info(msg, append(args, "err", nil)...)
}

func (l *Logger) WithFields(args ...any) *Logger {
	return &Logger{Logger: l.Logger.With(args...)}
}

func (l *Logger) WithLogID(ctx context.Context) *Logger {
	var (
		logId string
		ok    bool
	)
	logId, ok = ctx.Value(LogIDKey{}).(string)
	if !ok {
		logId = uuid.NewUUID().String()
	}

	return &Logger{Logger: l.Logger.With("LOGID", logId)}
}
