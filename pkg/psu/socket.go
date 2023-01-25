/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package psu

import (
	"net"
)

type socket struct {
	addr string
	net.Conn
}

func (s *socket) Open() (err error) {
	s.Conn, err = net.Dial("tcp", s.addr)
	return err
}
