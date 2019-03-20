package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

var mux = NewMux()

func handler(conn net.Conn) {
	defer conn.Close()

	mux.Add(conn)
	mux.BroadcastAll(fmt.Sprintf("%s joined the chat\n", conn.RemoteAddr().String()))

	for {
		conn.SetDeadline(time.Now().Add(15 * time.Second))
		var buf [128]byte
		n, err := conn.Read(buf[:])
		if err, ok := err.(net.Error); ok && err.Timeout() {
			mux.Remove(conn, "timed out")
			return
		}
		if err == io.EOF {
			mux.Remove(conn, "disconnected")
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		mux.BroadcastPeers(string(buf[:n]), conn)
	}
}

func main() {
	go mux.loop()

	l, err := net.Listen("tcp", "localhost:13990")
	if err != nil {
		panic(err)
	}
	fmt.Println("listening")
	for {
		conn, _ := l.Accept() // blocking
		go handler(conn)
	}
}
