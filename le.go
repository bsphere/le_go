package logentries

import (
	"fmt"
	"io"
	"net"
	"time"
)

type LE struct {
	conn  net.Conn
	token string
}

func Connect(token string) (*LE, error) {
	le := LE{
		token: token,
	}

	if err := le.reopenConnection(); err != nil {
		return nil, err
	}

	return &le, nil
}

func (le *LE) Close() error {
	if le.conn != nil {
		return le.conn.Close()
	}

	return nil
}

func (le *LE) reopenConnection() error {
	conn, err := net.Dial("tcp", "data.logentries.com:80")
	if err != nil {
		return err
	}

	le.conn = conn

	return nil
}

func (le *LE) isOpenConnection() bool {
	if le.conn == nil {
		return false
	}

	buf := make([]byte, 1)

	le.conn.SetReadDeadline(time.Now())

	if _, err := le.conn.Read(buf); err.(net.Error).Timeout() == true &&
		err != io.EOF {

		le.conn.SetReadDeadline(time.Time{})

		return true
	} else {
		le.conn.Close()

		return false
	}
}

func (le *LE) Println(msg string) {
	if !le.isOpenConnection() {
		if err := le.reopenConnection(); err != nil {
			return
		}
	}

	fmt.Fprintln(le.conn, le.token, msg)
}
