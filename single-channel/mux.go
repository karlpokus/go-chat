package main

import (
	"io"
	"net"
)

type Storage map[net.Addr]net.Conn

type Mux struct {
	ops chan func(Storage)
}

func (m *Mux) Add(conn net.Conn) {
	m.ops <- func(s Storage) {
		s[conn.RemoteAddr()] = conn
	}
}

func (m *Mux) Remove(conn net.Conn) {
	m.ops <- func(s Storage) {
		delete(s, conn.RemoteAddr())
	}
}

func (m *Mux) BroadcastAll(msg string) error {
	m.ops <- func(s Storage) {
		for _, conn := range s {
			io.WriteString(conn, msg)
		}
	}
	return nil
}

func (m *Mux) BroadcastPeers(msg string, me net.Conn) {
	m.ops <- func(s Storage) {
		for addr, conn := range s {
			if addr != me.RemoteAddr() {
				io.WriteString(conn, msg)
			}
		}
	}
}

func (m *Mux) loop() {
	s := make(Storage)
	for op := range m.ops { // never close ops chan
		op(s)
	}
}

func NewMux() *Mux {
	return &Mux{
		ops: make(chan func(Storage)),
	}
}
