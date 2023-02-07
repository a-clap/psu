/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package psu_test

import (
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"psu/pkg/psu"
	"testing"
	"time"
)

type ViewTestSuite struct {
	suite.Suite
	mock *AccessMocker
}

type AccessMocker struct {
	mock.Mock
}

func TestView(t *testing.T) {
	suite.Run(t, new(ViewTestSuite))
}

func (t *ViewTestSuite) SetupTest() {
	t.mock = new(AccessMocker)
}
func (t *ViewTestSuite) TestRefresh() {
	args := []struct {
		name       string
		number     []int
		retSection []*psu.Section
		retError   []error
	}{
		{
			name:   "single section",
			number: []int{0},
			retSection: []*psu.Section{
				{
					State:         false,
					ActualVoltage: "1",
					SetVoltage:    "2",
					ActualCurrent: "3",
					SetCurrent:    "4",
				},
			},
			retError: []error{nil},
		},
		{
			name:   "single section 2",
			number: []int{1},
			retSection: []*psu.Section{
				{
					State:         false,
					ActualVoltage: "1",
					SetVoltage:    "2",
					ActualCurrent: "3",
					SetCurrent:    "4",
				},
			},
			retError: []error{nil},
		},
		{
			name:   "two sections",
			number: []int{1, 2},
			retSection: []*psu.Section{
				{
					State:         false,
					ActualVoltage: "1",
					SetVoltage:    "2",
					ActualCurrent: "3",
					SetCurrent:    "4",
				},
				{
					State:         false,
					ActualVoltage: "1",
					SetVoltage:    "2",
					ActualCurrent: "3",
					SetCurrent:    "4",
				},
			},
			retError: []error{nil, nil},
		},
	}
	for _, arg := range args {
		r := t.Require()
		t.mock = new(AccessMocker)
		r.Equal(len(arg.number), len(arg.retSection))
		r.Equal(len(arg.retSection), len(arg.retError))

		for i := 0; i < len(arg.number); i++ {
			t.mock.On("Section", arg.number[i]).Return(arg.retSection[i], arg.retError[i])
		}

		v, err := psu.NewView(
			psu.ViewWithAccess(t.mock),
			psu.ViewWithSections(arg.number...),
		)
		_ = test.NewApp()

		r.Nil(err)
		r.NotNil(v)

		v.Refresh()
		// Force scheduler
		<-time.After(10 * time.Millisecond)
		t.mock.AssertExpectations(t.T())

	}

}

func (t *ViewTestSuite) TestNew() {
	{
		// No interface
		v, err := psu.NewView()
		t.Nil(v)
		t.NotNil(err)
		t.ErrorIs(psu.ErrNoAccess, err)
	}
	{
		// No section
		v, err := psu.NewView(psu.ViewWithAccess(t.mock))
		t.Nil(v)
		t.NotNil(err)
		t.ErrorIs(psu.ErrNoSection, err)
	}
	{
		// All good
		v, err := psu.NewView(psu.ViewWithAccess(t.mock), psu.ViewWithSections(0))
		t.NotNil(v)
		t.Nil(err)
	}
}

func (a *AccessMocker) Section(section int) (*psu.Section, error) {
	args := a.Called(section)
	return args.Get(0).(*psu.Section), args.Error(1)
}

func (a *AccessMocker) SetState(section int, value bool) (bool, error) {
	args := a.Called(section, value)
	return args.Bool(0), args.Error(1)
}
