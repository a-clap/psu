/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package psu

import (
	"go.uber.org/zap/zapcore"
	"time"
)

type Option func(*PSU) error

func WithConn(c Conn) Option {
	return func(psu *PSU) error {
		psu.conn = c
		return nil
	}
}

func WithSocketConn(host, port string) Option {
	return func(psu *PSU) error {
		s := &socket{
			addr: host + ":" + port,
			Conn: nil,
		}
		psu.conn = s
		return nil
	}
}
func WithReadWriteDeadline(t time.Duration) Option {
	return func(psu *PSU) error {
		psu.deadline = t
		return nil
	}
}

func WithLogger(l Logger) Option {
	return func(*PSU) error {
		log = l
		return nil
	}
}

func WithLogLevel(logLvl string) Option {
	return func(*PSU) error {
		lvl, err := zapcore.ParseLevel(logLvl)
		if err != nil {
			return err
		}
		log = NewDefaultZap(lvl)
		return nil
	}
}
