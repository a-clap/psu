/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package psu

type Option func(*PSU) error

func WithConn(c Conn) Option {
	return func(psu *PSU) error {
		psu.conn = c
		return nil
	}
}

func WithSocketConn(host, port string) Option {
	return func(psu *PSU) error {
		s := &socket{
			addr: host + ":" + port,
			Conn: nil,
		}
		psu.conn = s
		return nil
	}
}
