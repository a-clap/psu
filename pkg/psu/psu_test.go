/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package psu_test

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"psu/pkg/psu"
	"testing"
	"time"
)

type PSUTestSuite struct {
	suite.Suite
	mock *ConnMock
}

type ConnMock struct {
	mock.Mock
}

func TestPSU(t *testing.T) {
	suite.Run(t, new(PSUTestSuite))
}

func (t *PSUTestSuite) SetupTest() {
	t.mock = new(ConnMock)
}

func (t *PSUTestSuite) psu() *psu.PSU {
	p, _ := psu.New(psu.WithConn(t.mock))
	t.Require().NotNil(p)
	return p
}

func (t *PSUTestSuite) Test_SetState() {

	firstWrite := []byte("OP1 1\r\n")
	secondWrite := []byte("OP1?\r\n")
	expectedReply := []byte("1\r\n")

	r := t.Require()
	openCall := t.mock.On("Open").Return(nil)
	setDeadline := t.mock.On("SetDeadline", mock.Anything).Return(nil).NotBefore(openCall)
	firstWriteCall := t.mock.On("Write", firstWrite).Return(len(firstWrite), nil).NotBefore(setDeadline)
	secondWriteCall := t.mock.On("Write", secondWrite).Return(len(secondWrite), nil).NotBefore(firstWriteCall)

	readCall := t.mock.On("Read", mock.Anything).Return(len(expectedReply), nil).Run(func(args mock.Arguments) {
		buffer := args.Get(0).([]byte)
		copy(buffer, expectedReply)
	}).NotBefore(setDeadline, secondWriteCall)
	t.mock.On("Close").Return(nil).NotBefore(readCall)

	p := t.psu()
	v, err := p.SetState(1, true)
	r.Equal(true, v)
	r.Nil(err)
}

func (t *PSUTestSuite) Test_GetState() {

	expectedWrite := []byte("OP1?\r\n")
	expectedReply := []byte("1\r\n")

	r := t.Require()
	openCall := t.mock.On("Open").Return(nil)
	setDeadline := t.mock.On("SetDeadline", mock.Anything).Return(nil).NotBefore(openCall)
	writeCall := t.mock.On("Write", expectedWrite).Return(len(expectedWrite), nil).NotBefore(setDeadline)
	readCall := t.mock.On("Read", mock.Anything).Return(len(expectedReply), nil).Run(func(args mock.Arguments) {
		buffer := args.Get(0).([]byte)
		copy(buffer, expectedReply)
	}).NotBefore(setDeadline, writeCall)
	t.mock.On("Close").Return(nil).NotBefore(readCall)

	p := t.psu()
	v, err := p.State(1)
	r.Equal(true, v)
	r.Nil(err)
}

func (t *PSUTestSuite) Test_SetCurrent() {

	expectedWrite := []byte("I1?\r\n")
	expectedReply := []byte("I1 7.45\r\n")

	r := t.Require()
	openCall := t.mock.On("Open").Return(nil)
	setDeadline := t.mock.On("SetDeadline", mock.Anything).Return(nil).NotBefore(openCall)
	writeCall := t.mock.On("Write", expectedWrite).Return(len(expectedWrite), nil).NotBefore(setDeadline)
	readCall := t.mock.On("Read", mock.Anything).Return(len(expectedReply), nil).Run(func(args mock.Arguments) {
		buffer := args.Get(0).([]byte)
		copy(buffer, expectedReply)
	}).NotBefore(setDeadline, writeCall)
	t.mock.On("Close").Return(nil).NotBefore(readCall)

	p := t.psu()
	v, err := p.SetCurrent(1)
	r.Equal("7.45", v)
	r.Nil(err)
}

func (t *PSUTestSuite) Test_SetVoltage() {

	expectedWrite := []byte("V1?\r\n")
	expectedReply := []byte("V1 27.45\r\n")

	r := t.Require()
	openCall := t.mock.On("Open").Return(nil)
	setDeadline := t.mock.On("SetDeadline", mock.Anything).Return(nil).NotBefore(openCall)
	writeCall := t.mock.On("Write", expectedWrite).Return(len(expectedWrite), nil).NotBefore(setDeadline)
	readCall := t.mock.On("Read", mock.Anything).Return(len(expectedReply), nil).Run(func(args mock.Arguments) {
		buffer := args.Get(0).([]byte)
		copy(buffer, expectedReply)
	}).NotBefore(setDeadline, writeCall)
	t.mock.On("Close").Return(nil).NotBefore(readCall)

	p := t.psu()
	v, err := p.SetVoltage(1)
	r.Equal("27.45", v)
	r.Nil(err)
}

func (t *PSUTestSuite) Test_ActualVoltage() {

	expectedWrite := []byte("V1O?\r\n")
	expectedReply := []byte("123.45V\r\n")

	r := t.Require()
	openCall := t.mock.On("Open").Return(nil)
	setDeadline := t.mock.On("SetDeadline", mock.Anything).Return(nil).NotBefore(openCall)
	writeCall := t.mock.On("Write", expectedWrite).Return(len(expectedWrite), nil).NotBefore(setDeadline)
	readCall := t.mock.On("Read", mock.Anything).Return(len(expectedReply), nil).Run(func(args mock.Arguments) {
		buffer := args.Get(0).([]byte)
		copy(buffer, expectedReply)
	}).NotBefore(setDeadline, writeCall)
	t.mock.On("Close").Return(nil).NotBefore(readCall)

	p := t.psu()
	v, err := p.ActualVoltage(1)
	r.Equal("123.45", v)
	r.Nil(err)
}

func (t *PSUTestSuite) Test_ActualCurrent() {

	expectedWrite := []byte("I1O?\r\n")
	expectedReply := []byte("123.45A\r\n")

	r := t.Require()
	openCall := t.mock.On("Open").Return(nil)
	setDeadline := t.mock.On("SetDeadline", mock.Anything).Return(nil).NotBefore(openCall)
	writeCall := t.mock.On("Write", expectedWrite).Return(len(expectedWrite), nil).NotBefore(setDeadline)
	readCall := t.mock.On("Read", mock.Anything).Return(len(expectedReply), nil).Run(func(args mock.Arguments) {
		buffer := args.Get(0).([]byte)
		copy(buffer, expectedReply)
	}).NotBefore(setDeadline, writeCall)
	t.mock.On("Close").Return(nil).NotBefore(readCall)

	p := t.psu()
	v, err := p.ActualCurrent(1)
	r.Equal("123.45", v)
	r.Nil(err)
}

func (t *PSUTestSuite) TestNew() {
	r := t.Require()
	{
		p, err := psu.New()
		r.Nil(p)
		r.NotNil(err)
		r.ErrorIs(psu.ErrNoConnInterface, err)
	}
	{
		p, err := psu.New(psu.WithConn(t.mock))
		r.NotNil(p)
		r.Nil(err)
	}
}

func (c *ConnMock) Open() error {
	return c.Called().Error(0)
}

func (c *ConnMock) SetDeadline(t time.Time) error {
	return c.Called(t).Error(0)
}

func (c *ConnMock) Read(p []byte) (n int, err error) {
	args := c.Called(p)
	return args.Int(0), args.Error(1)
}

func (c *ConnMock) Write(p []byte) (n int, err error) {
	args := c.Called(p)
	return args.Int(0), args.Error(1)
}

func (c *ConnMock) Close() error {
	return c.Called().Error(0)
}
