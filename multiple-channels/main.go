package main

import (
	"io"
	"log"
	"net"
)

type Mux struct {
	add     chan net.Conn
	remove  chan net.Addr
	sendMsg chan string
}

func (m *Mux) Add(conn net.Conn) {
	m.add <- conn
}

func (m *Mux) Remove(addr net.Addr) {
	m.remove <- addr
}

func (m *Mux) SendMsg(msg string) error {
	m.sendMsg <- msg
	return nil
}

func (m *Mux) loop() {
	conns := make(map[net.Addr]net.Conn)
	for {
		select {
		case conn := <-m.add:
			conns[conn.RemoteAddr()] = conn
		case addr := <-m.remove:
			delete(conns, addr)
		case msg := <-m.sendMsg:
			for _, conn := range conns {
				io.WriteString(conn, msg)
			}
		}
	}
}

func main() {
	mux := &Mux{
		add:     make(chan net.Conn),
		remove:  make(chan net.Addr),
		sendMsg: make(chan string),
	}
	go mux.loop()

	l, err := net.Listen("tcp", "localhost:13990")
	if err != nil {
		panic(err)
	}
	log.Println("listening")
	for {
		conn, _ := l.Accept() // blocking
		mux.Add(conn)
		mux.SendMsg("welcome")
	}
}
