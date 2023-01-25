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
	_ commander = (*setVoltageType)(nil)
)

type actualVoltageType struct {
	section int
}

type setVoltageType struct {
	section int
}

type actualCurrentType struct {
	section int
}

type setCurrentType struct {
	section int
}

func (*setCurrentType) Parse(reply []string) (string, error) {
	if len(reply) != 2 {
		return "", ErrUnexpectedLen
	}
	return reply[1], nil
}

func (*setCurrentType) WriteOnly() bool {
	return false
}

func (s *setCurrentType) Command() command {
	return command("I" + strconv.FormatInt(int64(s.section), 10) + "?")
}

func (*actualCurrentType) Parse(reply []string) (string, error) {
	if len(reply) != 1 {
		return "", ErrUnexpectedLen
	}
	return strings.TrimSuffix(reply[0], "A"), nil
}

func (*actualCurrentType) WriteOnly() bool {
	return false
}

func (a *actualCurrentType) Command() command {
	return command("I" + strconv.FormatInt(int64(a.section), 10) + "O?")
}

func (*setVoltageType) Parse(reply []string) (string, error) {
	if len(reply) != 2 {
		return "", ErrUnexpectedLen
	}
	return reply[1], nil
}

func (*setVoltageType) WriteOnly() bool {
	return false
}

func (s *setVoltageType) Command() command {
	return command("V" + strconv.FormatInt(int64(s.section), 10) + "?")
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
