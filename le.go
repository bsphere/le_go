package logentries

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Logger struct {
	conn   net.Conn
	flags  int
	prefix string
	token  string
}

func Connect(token string) (*Logger, error) {
	Logger := Logger{
		token: token,
	}

	if err := Logger.reopenConnection(); err != nil {
		return nil, err
	}

	return &Logger, nil
}

func (Logger *Logger) Close() error {
	if Logger.conn != nil {
		return Logger.conn.Close()
	}

	return nil
}

func (Logger *Logger) reopenConnection() error {
	conn, err := net.Dial("tcp", "data.logentries.com:80")
	if err != nil {
		return err
	}

	Logger.conn = conn

	return nil
}

func (Logger *Logger) isOpenConnection() bool {
	if Logger.conn == nil {
		return false
	}

	buf := make([]byte, 1)

	Logger.conn.SetReadDeadline(time.Now())

	if _, err := Logger.conn.Read(buf); err.(net.Error).Timeout() == true &&
		err != io.EOF {

		Logger.conn.SetReadDeadline(time.Time{})

		return true
	} else {
		Logger.conn.Close()

		return false
	}
}

func (Logger *Logger) ensureOpenConnection() error {
	if !Logger.isOpenConnection() {
		if err := Logger.reopenConnection(); err != nil {
			return err
		}
	}

	return nil
}

func (Logger *Logger) Fatal(v ...interface{}) {
	if err := Logger.ensureOpenConnection(); err != nil {
		return
	}

	fmt.Fprint(Logger.conn, Logger.token, Logger.prefix, v)
	os.Exit(1)
}

func (Logger *Logger) Fatalf(v ...interface{}) {
	if err := Logger.ensureOpenConnection(); err != nil {
		return
	}

	fmt.Fprintf(Logger.conn, Logger.token, Logger.prefix, v)
	os.Exit(1)
}

func (Logger *Logger) Fatalln(v ...interface{}) {
	if err := Logger.ensureOpenConnection(); err != nil {
		return
	}

	fmt.Fprintln(Logger.conn, Logger.token, Logger.prefix, v)
	os.Exit(1)
}

func (Logger *Logger) Flags() int {
	return Logger.flags
}

func (Logger *Logger) Output(calldepth int, s string) error {
	return nil
}

func (Logger *Logger) Panic(v ...interface{}) {
	if err := Logger.ensureOpenConnection(); err != nil {
		return
	}

	fmt.Fprint(Logger.conn, Logger.token, Logger.prefix, fmt.Sprint(v...))
	panic("")
}

func (Logger *Logger) Panicf(format string, v ...interface{}) {
	if err := Logger.ensureOpenConnection(); err != nil {
		return
	}

	fmt.Fprintf(Logger.conn, Logger.token, Logger.prefix, fmt.Sprintf(format, v...))
	panic("")
}

func (Logger *Logger) Panicln(v ...interface{}) {
	if err := Logger.ensureOpenConnection(); err != nil {
		return
	}

	fmt.Fprintln(Logger.conn, Logger.token, Logger.prefix, fmt.Sprintln(v...))
	panic("")
}

func (Logger *Logger) Prefix() string {
	return Logger.prefix
}

func (Logger *Logger) Print(v ...interface{}) {
	if err := Logger.ensureOpenConnection(); err != nil {
		return
	}

	fmt.Fprint(Logger.conn, Logger.token, Logger.prefix, fmt.Sprint(v...))
}

func (Logger *Logger) Printf(format string, v ...interface{}) {
	if err := Logger.ensureOpenConnection(); err != nil {
		return
	}

	fmt.Fprintf(Logger.conn, Logger.token, Logger.prefix, fmt.Sprintf(format, v...))
}

func (Logger *Logger) Println(v ...interface{}) {
	if err := Logger.ensureOpenConnection(); err != nil {
		return
	}

	fmt.Fprintln(Logger.conn, Logger.token, Logger.prefix, fmt.Sprintln(v...))
}

func (Logger *Logger) SetFlags(flag int) {
	Logger.flags = flag
}

func (Logger *Logger) SetPrefix(prefix string) {
	Logger.prefix = prefix
}
