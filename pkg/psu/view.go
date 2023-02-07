/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package psu

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"time"
)

type View struct {
	psu            Access
	sectionNumbers []int
	sections       []*viewSection
	trigger, close chan struct{}
	ticker         *time.Ticker
	refreshButton  *widget.Button
}

type viewSection struct {
	section int
	psu     Access
	number  *widget.Label
	voltage *widget.Label
	current *widget.Label
	enable  *widget.Button
}

type Access interface {
	Section(section int) (*Section, error)
	SetState(section int, value bool) (bool, error)
}

var (
	ErrNoAccess  = errors.New("no Access interface")
	ErrNoSection = errors.New("no section to handle")
)

func NewView(opts ...ViewOption) (*View, error) {
	v := &View{
		psu:           nil,
		sections:      nil,
		trigger:       make(chan struct{}),
		close:         make(chan struct{}),
		ticker:        time.NewTicker(1 * time.Hour),
		refreshButton: widget.NewButtonWithIcon("", theme.MediaReplayIcon(), nil),
	}
	v.ticker.Stop()
	v.refreshButton.OnTapped = func() {
		v.Refresh()
	}
	for _, opt := range opts {
		if err := opt(v); err != nil {
			return nil, err
		}
	}

	if err := v.verify(); err != nil {
		return nil, err
	}

	v.sections = make([]*viewSection, len(v.sectionNumbers))
	for i, sec := range v.sectionNumbers {
		v.sections[i] = newViewSection(sec, v.psu)
	}

	go v.backgroundRefresh()

	return v, nil
}

func (v *View) Content() fyne.CanvasObject {
	title := container.NewHBox(layout.NewSpacer(), widget.NewLabel("CPX400"), layout.NewSpacer(), v.refreshButton)
	sections := len(v.sections)

	number := container.NewGridWithColumns(sections)
	enable := container.NewGridWithColumns(sections)
	voltage := container.NewGridWithColumns(sections)
	current := container.NewGridWithColumns(sections)
	for _, section := range v.sections {
		number.Add(section.number)
		enable.Add(section.enable)
		voltage.Add(section.voltage)
		current.Add(section.current)
	}
	return container.NewGridWithRows(5,
		title,
		number,
		enable,
		voltage,
		current)

}

func (v *View) Refresh() {
	v.trigger <- struct{}{}
}

func (v *View) Close() {
	close(v.close)
}

func (v *View) BackgroundRefresh(t time.Duration) {
	v.ticker.Reset(t)
}
func (v *View) StopBackgroundRefresh() {
	v.ticker.Stop()
}

func (v *View) backgroundRefresh() {
	running := true
	for running {
		select {
		case <-v.close:
			running = false
		case <-v.trigger:
			v.refresh()
		case <-v.ticker.C:
			v.refresh()
		}
	}
}

func (v *View) refresh() {
	for _, section := range v.sections {
		section.refresh()
	}
}

func newViewSection(number int, access Access) *viewSection {
	section := strconv.FormatInt(int64(number), 32)
	v := &viewSection{
		section: number,
		psu:     access,
		number:  widget.NewLabelWithStyle(section, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		voltage: widget.NewLabelWithStyle("0/8", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		current: widget.NewLabelWithStyle("0/13", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		enable:  widget.NewButton("", func() {}),
	}
	v.enable.Importance = widget.HighImportance
	return v
}

func (vs *viewSection) refresh() {
	data, err := vs.psu.Section(vs.section)
	if err != nil {
		const errText = "err"
		vs.voltage.SetText(errText)
		vs.current.SetText(errText)
		return
	}
	text := data.ActualVoltage + " / " + data.SetVoltage + " V DC"
	vs.voltage.SetText(text)

	text = data.ActualCurrent + " / " + data.SetCurrent + " A"
	vs.current.SetText(text)

	vs.enable.OnTapped = func() {
		_, _ = vs.psu.SetState(vs.section, !data.State)
		vs.refresh()
	}

	if data.State {
		vs.enable.Icon = theme.MediaStopIcon()
		vs.enable.SetText("OFF")
	} else {
		vs.enable.Icon = theme.MediaStopIcon()
		vs.enable.SetText("ON")
	}

}

func (v *View) verify() error {
	if v.psu == nil {
		return ErrNoAccess
	}
	if len(v.sectionNumbers) == 0 {
		return ErrNoSection
	}
	return nil
}
