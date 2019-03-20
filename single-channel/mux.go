package main

import (
	"fmt"
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
		logConns(s)
	}
	logAction(conn, "joined the chat")
}

func (m *Mux) Remove(conn net.Conn, action string) {
	m.ops <- func(s Storage) {
		delete(s, conn.RemoteAddr())
		logConns(s)
	}
	logAction(conn, action)
}

func (m *Mux) BroadcastAll(msg string) error {
	m.ops <- func(s Storage) {
		for _, conn := range s {
			io.WriteString(conn, formatMsg("server", msg))
		}
	}
	return nil
}

func (m *Mux) BroadcastPeers(msg string, me net.Conn) {
	m.ops <- func(s Storage) {
		for addr, conn := range s {
			if addr != me.RemoteAddr() {
				io.WriteString(conn, formatMsg(me.RemoteAddr().String(), msg))
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

func formatMsg(src string, msg string) string {
	return fmt.Sprintf("%s: %s", src, msg)
}

func logConns(s Storage) {
	fmt.Printf("%d open connections\n", len(s))
}

func logAction(conn net.Conn, action string) {
	fmt.Printf("%s %s\n", conn.RemoteAddr().String(), action)
}
