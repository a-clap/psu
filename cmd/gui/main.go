/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"os"
	"path/filepath"
	"psu/pkg/psu"
	"time"
)

type config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Sections []int  `json:"sections"`
}

func main() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	file, err := os.ReadFile(filepath.Join(path, "config.json"))
	if err != nil {
		panic(err)
	}

	cfg := config{}
	if err := json.Unmarshal(file, &cfg); err != nil {
		panic(err)
	}

	p, err := psu.New(psu.WithSocketConn(cfg.Host, cfg.Port))
	if err != nil {
		panic(err)
	}

	v, err := psu.NewView(
		psu.ViewWithPSU(p),
		psu.ViewWithSections(cfg.Sections...),
	)

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
	v.Refresh()

	ctn := container.NewMax(v.Content())
	w := gui.NewWindow("CPX400DP")
	w.SetContent(ctn)

	w.Resize(fyne.NewSize(280, 160))
	w.ShowAndRun()

}
