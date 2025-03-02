package logger

import (
	sdklog "cosmossdk.io/log"
	"github.com/0xPellNetwork/pelldvs/libs/log"
)

type DVSLogAdapter struct {
	sdkLogger sdklog.Logger
}

func NewDVSLogAdapter(sdkLogger sdklog.Logger) log.Logger {
	return &DVSLogAdapter{sdkLogger: sdkLogger}
}

func (a *DVSLogAdapter) Debug(msg string, keyVals ...any) {
	a.sdkLogger.Debug(msg, keyVals...)
}

func (a *DVSLogAdapter) Info(msg string, keyVals ...any) {
	a.sdkLogger.Info(msg, keyVals...)
}

func (a *DVSLogAdapter) Error(msg string, keyVals ...any) {
	a.sdkLogger.Error(msg, keyVals...)
}

func (a *DVSLogAdapter) With(keyVals ...any) log.Logger {
	return &DVSLogAdapter{sdkLogger: a.sdkLogger.With(keyVals...)}
}
