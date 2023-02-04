/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package psu

type ViewOption func(*View) error

func ViewWithPSU(p *PSU) ViewOption {
	return func(view *View) error {
		return ViewWithAccess(p)(view)
	}
}

func ViewWithAccess(a Access) ViewOption {
	return func(view *View) error {
		view.psu = a
		return nil
	}
}

func ViewWithSections(sections ...int) ViewOption {
	return func(view *View) error {
		view.sectionNumbers = append(view.sectionNumbers, sections...)
		return nil
	}
}
