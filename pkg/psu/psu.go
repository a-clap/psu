/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package psu

import (
	"errors"
	"io"
	"strings"
	"time"
)

type Conn interface {
	Open() error
	SetDeadline(time.Time) error
	io.ReadWriteCloser
}

type PSU struct {
	conn     Conn
	deadline time.Duration
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

func (p *PSU) ActualVoltage(section int) (string, error) {
	av := &actualVoltageType{section: section}
	reply, err := p.communicate(av)
	if err != nil {
		return "", err
	}
	return reply[av.Command()], nil
}

func (p *PSU) SetVoltage(section int) (string, error) {
	av := &setVoltageType{section: section}
	reply, err := p.communicate(av)
	if err != nil {
		return "", err
	}
	return reply[av.Command()], nil
}

func (p *PSU) communicate(cmds ...commander) (map[command]string, error) {
	reply := make(map[command]string)

	log.Debug("Connecting ...")
	if err := p.conn.Open(); err != nil {
		log.Error("Failed to connect: ", err)
		return nil, err
	}
	defer func() {
		log.Debug("Disconnecting...")
		if err := p.conn.Close(); err != nil {
			log.Error("Failed to disconnect: ", err)
		}
	}()

	for _, cmd := range cmds {
		p.setDeadline()
		writeCmd := cmd.Command()
		log.Debug("Writing to Conn: ", writeCmd)
		if _, err := p.conn.Write([]byte(writeCmd + "\r\n")); err != nil {
			log.Error("error on Write: ", err)
			return reply, err
		}
		if cmd.WriteOnly() {
			continue
		}
		// CPX usually respond within few bytes
		readBuffer := make([]byte, 64)
		p.setDeadline()
		size, err := p.conn.Read(readBuffer)
		if err != nil {
			log.Error("error on Read: ", err)
			return nil, err
		}
		data := strings.TrimSuffix(string(readBuffer[:size]), "\r\n")
		log.Debug("received data: ", data)
		cmdReply, err := cmd.Parse(strings.Split(data, " "))
		if err != nil {
			log.Errorf("error: %s, on parsing cmd %s\n", err, writeCmd)
			continue
		}
		reply[writeCmd] = cmdReply
	}

	return reply, nil
}

func (p *PSU) setDeadline() {
	if err := p.conn.SetDeadline(time.Now().Add(p.deadline)); err != nil {
		log.Error("Error on setting deadline: ", err)
	}
}
