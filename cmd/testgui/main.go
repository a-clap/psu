/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package main

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"math/rand"
	"psu/pkg/psu"
	"strconv"
	"time"
)

type access struct {
	sections []*psu.Section
}

func newAcccess(sections ...int) *access {
	a := &access{}
	a.sections = make([]*psu.Section, len(sections))
	for i, section := range sections {
		section += 1
		a.sections[i] = &psu.Section{
			State:         false,
			ActualVoltage: strconv.FormatInt(int64(section), 32),
			SetVoltage:    strconv.FormatInt(int64(section*100), 32),
			ActualCurrent: strconv.FormatInt(int64(section), 32),
			SetCurrent:    strconv.FormatInt(int64(section*100), 32),
		}
	}
	return a
}

func (a *access) Section(section int) (*psu.Section, error) {
	if section >= len(a.sections) {
		return nil, errors.New("no such section")
	}
	const min = 0.0
	const max = 30.0

	setVoltage := min + rand.Float64()*(max-min)
	actualVoltage := min + rand.Float64()*(setVoltage-min)

	setCurrent := min + rand.Float64()*(max-min)
	actualCurrent := min + rand.Float64()*(setCurrent-min)

	fmt := func(value float64) string {
		return strconv.FormatFloat(value, 'f', 2, 32)
	}

	a.sections[section].ActualVoltage = fmt(actualVoltage)
	a.sections[section].SetVoltage = fmt(setVoltage)
	a.sections[section].ActualCurrent = fmt(actualCurrent)
	a.sections[section].SetCurrent = fmt(setCurrent)

	return a.sections[section], nil

}

func (a access) SetState(section int, value bool) (bool, error) {
	if section >= len(a.sections) {
		return false, errors.New("no such section")
	}

	a.sections[section].State = value
	return a.sections[section].State, nil

}

var (
	_ psu.Access = (*access)(nil)
)

func main() {
	a := newAcccess(0, 1)
	v, err := psu.NewView(
		psu.ViewWithAccess(a),
		psu.ViewWithSections(0, 1))
	if err != nil {
		panic(err)
	}

	gui := app.New()
	gui.Settings().SetTheme(theme.DarkTheme())

	gui.Lifecycle().SetOnEnteredForeground(func() {
		v.BackgroundRefresh(1 * time.Second)
	})

	gui.Lifecycle().SetOnExitedForeground(func() {
		v.StopBackgroundRefresh()
	})

	// Update data before screen is started
	go v.Refresh()

	ctn := container.NewMax(v.Content())
	w := gui.NewWindow("CPX400DP")
	w.SetContent(ctn)

	w.Resize(fyne.NewSize(280, 160))
	w.ShowAndRun()

}
