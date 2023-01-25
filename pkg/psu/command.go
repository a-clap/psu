/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package psu

import (
	"errors"
	"strconv"
	"strings"
)

type command string

type commander interface {
	Parse(reply []string) (string, error)
	WriteOnly() bool
	Command() command
}

var (
	ErrUnexpectedLen = errors.New("unexpected reply length")
)

var (
	_ commander = (*actualVoltageType)(nil)
)

type actualVoltageType struct {
	section int
}

func (*actualVoltageType) Parse(reply []string) (string, error) {
	if len(reply) != 1 {
		return "", ErrUnexpectedLen
	}
	return strings.TrimSuffix(reply[0], "V"), nil
}

func (*actualVoltageType) WriteOnly() bool {
	return false
}

func (a *actualVoltageType) Command() command {
	return command("V" + strconv.FormatInt(int64(a.section), 10) + "O?")
}
