package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewSugarLogger() *zap.SugaredLogger {
	// Define custom encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // Capitalize the log level names
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC timestamp format
		EncodeDuration: zapcore.SecondsDurationEncoder, // Duration in seconds
		EncodeCaller:   zapcore.ShortCallerEncoder,     // Short caller (file and line)
	}
	
	// Define log file
	logFile, err := os.OpenFile("app.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}
	
	//Define multiwriter
	multiWriter := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(logFile),
		zapcore.AddSync(os.Stdout),
	)
	// Create a core logger with JSON encoding
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // Using JSON encoder
		multiWriter,
		zap.DebugLevel, // устанавливаем уровень логирования
	)

	// Показываем файл, где возникла ошибки и stack трэйс для уровня zap.ErrorLevel
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCallerSkip(1),
	)

	sugar := logger.Sugar()

	return sugar
}
