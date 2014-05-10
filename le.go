package logentries

import (
	"fmt"
	"net"
)

type LE struct {
	Conn  net.Conn
	Token string
}

func Connect(token string) (*LE, error) {
	le := LE{
		Token: token,
	}

	conn, err := net.Dial("tcp", "data.logentries.com:80")
	if err != nil {
		return nil, err
	}

	le.Conn = conn

	return &le, nil
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

	fmt.Fprintln(le.Conn, le.Token, msg)
}
