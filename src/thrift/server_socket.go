/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package thrift

import (
	"bmautil/netutil"
	"logger"
	"net"
	"time"
)

type TServerSocket struct {
	listener      net.Listener
	addr          string
	clientTimeout time.Duration
	interrupted   bool

	// ip limit
	WhiteList func() []string
	BlackList func() []string
}

func NewTServerSocket(listenAddr string) (*TServerSocket, error) {
	return NewTServerSocketTimeout(listenAddr, 0)
}

func NewTServerSocketTimeout(listenAddr string, clientTimeout time.Duration) (*TServerSocket, error) {
	addr := listenAddr
	return &TServerSocket{addr: addr, clientTimeout: clientTimeout}, nil
}

func (p *TServerSocket) Listen() error {
	if p.IsListening() {
		return nil
	}
	l, err := net.Listen("tcp", p.addr)
	if err != nil {
		return err
	}
	p.listener = l
	return nil
}

func (p *TServerSocket) Accept() (TTransport, error) {
	if p.interrupted {
		return nil, errTransportInterrupted
	}
	if p.listener == nil {
		return nil, NewTTransportException(NOT_OPEN, "No underlying server socket")
	}
	for {
		conn, err := p.listener.Accept()
		if err != nil {
			return nil, NewTTransportExceptionFromError(err)
		}
		addr := conn.RemoteAddr().String()
		var wl, bl []string
		if p.WhiteList != nil {
			wl = p.WhiteList()
		}
		if p.BlackList != nil {
			bl = p.BlackList()
		}
		if ok, msg := netutil.IpAccept(addr, wl, bl, true); !ok {
			logger.Warn(tag, "unaccept(%s) address %s", msg, addr)
			conn.Close()
			continue
		}
		return NewTSocketFromConnTimeout(conn, p.clientTimeout), nil
	}
}

// Checks whether the socket is listening.
func (p *TServerSocket) IsListening() bool {
	return p.listener != nil
}

// Connects the socket, creating a new socket object if necessary.
func (p *TServerSocket) Open() error {
	if p.IsListening() {
		return NewTTransportException(ALREADY_OPEN, "Server socket already open")
	}
	if l, err := net.Listen("tcp", p.addr); err != nil {
		return err
	} else {
		p.listener = l
	}
	return nil
}

func (p *TServerSocket) Close() error {
	defer func() {
		p.listener = nil
	}()
	if p.IsListening() {
		return p.listener.Close()
	}
	return nil
}

func (p *TServerSocket) Interrupt() error {
	p.interrupted = true
	return nil
}
