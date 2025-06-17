// logger/global.go
package logger

import (
    "go.uber.org/zap"
)

// Global SugaredLogger instance
var (
    Sugar *zap.SugaredLogger
)

func InitLogger() {
    Sugar = NewSugarLogger()
}