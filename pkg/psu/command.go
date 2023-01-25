/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package psu

import (
	"errors"
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
	_ commander = (*actualCurrentType)(nil)
	_ commander = (*setCurrentType)(nil)
	_ commander = (*getStateType)(nil)
	_ commander = (*setStateType)(nil)
)

type actualVoltageType struct {
	section string
}

type setVoltageType struct {
	section string
}

type actualCurrentType struct {
	section string
}

type setCurrentType struct {
	section string
}

type getStateType struct {
	section string
}

type setStateType struct {
	section string
	value   bool
}

func (*setStateType) Parse(reply []string) (string, error) {
	panic("shouldn't be called")
}

func (*setStateType) WriteOnly() bool {
	return true
}

func (s *setStateType) Command() command {
	value := "1"
	if !s.value {
		value = "0"
	}
	return command("OP" + s.section + " " + value)
}

func (*getStateType) Parse(reply []string) (string, error) {
	if len(reply) != 1 {
		return "", ErrUnexpectedLen
	}
	return reply[0], nil
}

func (*getStateType) WriteOnly() bool {
	return false
}

func (g *getStateType) Command() command {
	return command("OP" + g.section + "?")
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
	return command("I" + s.section + "?")
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
	return command("I" + a.section + "O?")
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
	return command("V" + s.section + "?")
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
	return command("V" + a.section + "O?")
}
