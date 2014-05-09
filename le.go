package logentries

import (
	"fmt"
	"net"
)

type LE struct {
	Conn  net.Conn
	Token string
}

func (le *LE) Init() error {
	conn, err := net.Dial("tcp", "data.logentries.com:80")
	if err != nil {
		return err
	}

	le.Conn = conn

	return nil
}

func (le *LE) Close() error {
	if le.Conn != nil {
		return le.Conn.Close()
	}

	return nil
}

func (le *LE) Println(msg string) {
	if le.Conn == nil {
		return
	}
}
