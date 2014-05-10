package logentries

import (
	"fmt"
	"io"
	"net"
	"time"
)

type LE struct {
	Conn  net.Conn
	Token string
}

func Connect(token string) (*LE, error) {
	le := LE{
		Token: token,
	}

	if err := le.reopenConnection(); err != nil {
		return nil, err
	}

	return &le, nil
}

func (le *LE) Close() error {
	if le.Conn != nil {
		return le.Conn.Close()
	}

	return nil
}

func (le *LE) reopenConnection() error {
	conn, err := net.Dial("tcp", "data.logentries.com:80")
	if err != nil {
		return err
	}

	le.Conn = conn

	return nil
}

func (le *LE) isOpenConnection() bool {
	if le.Conn == nil {
		return false
	}

	buf := make([]byte, 1)

	le.Conn.SetReadDeadline(time.Now())

	if _, err := le.Conn.Read(buf); err == io.EOF {
		le.Conn.Close()

		return false
	} else {
		le.Conn.SetReadDeadline(time.Time{})

		return true
	}
}

func (le *LE) Println(msg string) {
	if !le.isOpenConnection() {
		if err := le.reopenConnection(); err != nil {
			return
		}
	}

	fmt.Fprintln(le.Conn, le.Token, msg)
}
