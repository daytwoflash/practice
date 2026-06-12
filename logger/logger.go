package logger

import (
	"log/slog"
	"os"
)

func Init(mode string) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var handler slog.Handler
	switch mode {
	case "release":
		// 生产环境：Info级，JSON格式便于采集，标准输出（可改为日志文件）
		opts.Level = slog.LevelInfo
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		// 开发环境：debug级，文本格式输出便于终端读取
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	// 用 handler 创建 logger，并设此 logger 为全局默认
	slog.SetDefault(slog.New(handler))
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
