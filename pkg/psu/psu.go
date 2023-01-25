/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package psu

import (
	"errors"
	"io"
	"time"
)

type Conn interface {
	Open() error
	SetDeadline(time.Time) error
	io.ReadWriteCloser
}

type PSU struct {
	conn Conn
}

var (
	ErrNoConnInterface = errors.New("lack of Conn interface")
)

func New(options ...Option) (*PSU, error) {
	p := new(PSU)
	for _, option := range options {
		if err := option(p); err != nil {
			return nil, err
		}
	}
	if err := p.verify(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *PSU) verify() error {
	if p.conn == nil {
		return ErrNoConnInterface
	}
	return nil
}
