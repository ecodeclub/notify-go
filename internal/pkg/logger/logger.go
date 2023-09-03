package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	l     *zap.Logger // zap ensure that zap.Logger is safe for concurrent use
	level Level
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}
func (l *Logger) DPanic(msg string, fields ...Field) {
	l.l.DPanic(msg, fields...)
}
func (l *Logger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

// New create a new logger (not support log rotating).
func New(writer io.Writer, level Level, opts ...Option) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(writer),
		level,
	)
	logger := &Logger{
		l:     zap.New(core),
		level: level,
	}
	return logger
}

type LevelEnablerFunc func(lvl Level) bool

type RotateOptions struct {
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

type TeeOption struct {
	Filename string
	Ropts    RotateOptions
	Lef      LevelEnablerFunc
}

func NewTee(tops []TeeOption, opts ...Option) *Logger {
	var cores []zapcore.Core
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
	for _, top := range tops {
		top := top

		// 多输出的日志级别的固定写法？
		lv := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return top.Lef(lvl)
		})

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   top.Filename,
				MaxSize:    top.Ropts.MaxSize,
				MaxBackups: top.Ropts.MaxBackups,
				MaxAge:     top.Ropts.MaxAge,
				Compress:   top.Ropts.Compress,
			}),
			lv,
		)
		cores = append(cores, core)
	}

	logger := &Logger{
		l: zap.New(zapcore.NewTee(cores...), opts...),
	}
	return logger
}

//func Sync() error {
//	if std != nil {
//		return std.Sync()
//	}
//	return nil
//}

var (
	std   = New(os.Stderr, InfoLevel)
	multi = NewTee(func() []TeeOption {
		return []TeeOption{
			{
				Filename: "./log/access.log",
				Lef: func(lvl Level) bool {
					return lvl <= FatalLevel
				},
				Ropts: RotateOptions{
					MaxSize:    1,    // Mb
					MaxAge:     2,    // 最多保留2天
					MaxBackups: 3,    // 最多保留3个压缩文件
					Compress:   true, // 历史文件是否压缩
				},
			},
			{
				Filename: "./log/error.log",
				Lef: func(lvl Level) bool {
					return lvl > InfoLevel
				},
				Ropts: RotateOptions{
					MaxSize:    1,
					MaxAge:     2,
					MaxBackups: 3,
					Compress:   true,
				},
			},
		}
	}())
)

func Default() *Logger {
	return multi
}

// ResetDefault not safe for concurrent use
func ResetDefault(l *Logger) {
	std = l
	Info = std.Info
	Warn = std.Warn
	Error = std.Error
	DPanic = std.DPanic
	Panic = std.Panic
	Fatal = std.Fatal
	Debug = std.Debug
}
