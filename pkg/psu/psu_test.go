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
