/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package psu

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is simple interface to notify user about serious problems (or just to debug)
type Logger interface {
	Errorf(format string, args ...interface{})
	Error(args ...interface{})

	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
}

func NewDefaultZap(level zapcore.Level) *zap.SugaredLogger {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return log.Sugar()
}

func NewNop() *zap.SugaredLogger {
	return zap.NewNop().Sugar()
}

var log Logger = NewDefaultZap(zapcore.DebugLevel)
