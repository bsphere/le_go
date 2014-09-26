// Package le_go provides a Golang client library for logging to
// logentries.com over a TCP connection.
//
// it uses an access token for sending log events.
package le_go

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// Logger represents a Logentries logger,
// it holds the open TCP connection, access token, prefix and flags
type Logger struct {
	conn   net.Conn
	flag   int
	mu     sync.Mutex
	prefix string
	token  string
	buf    []byte
}

const lineSep = "\n"

// Connte creates a new Logger instance and opens a TCP connection to
// logentries.com,
// The token can be generated at logentries.com by adding a new log,
// choosing manual configuration and token based TCP connection.
func Connect(token string) (*Logger, error) {
	logger := Logger{
		token: token,
	}

	if err := logger.openConnection(); err != nil {
		return nil, err
	}

	return &logger, nil
}

// Close closes the TCP connection to logentries.com
func (logger *Logger) Close() error {
	if logger.conn != nil {
		return logger.conn.Close()
	}

	return nil
}

// Opens a TCP connection to logentries.com
func (logger *Logger) openConnection() error {
	conn, err := net.Dial("tcp", "data.logentries.com:80")
	if err != nil {
		return err
	}

	logger.conn = conn

	return nil
}

// It returns if the TCP connection to logentries.com is open
func (logger *Logger) isOpenConnection() bool {
	if logger.conn == nil {
		return false
	}

	buf := make([]byte, 1)

	logger.conn.SetReadDeadline(time.Now())

	_, err := logger.conn.Read(buf)

	switch err.(type) {
	case net.Error:
		if err.(net.Error).Timeout() == true {
			logger.conn.SetReadDeadline(time.Time{})

			return true
		}
	}

	return false
}

// It ensures that the TCP connection to logentries.com is open.
// If the connection is closed, a new one is opened.
func (logger *Logger) ensureOpenConnection() error {
	if !logger.isOpenConnection() {
		if err := logger.openConnection(); err != nil {
			return err
		}
	}

	return nil
}

func (logger *Logger) Fatal(v ...interface{}) {
	logger.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func (logger *Logger) Fatalf(format string, v ...interface{}) {
	logger.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (logger *Logger) Fatalln(v ...interface{}) {
	logger.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

func (logger *Logger) Flags() int {
	return logger.flag
}

func (logger *Logger) Output(calldepth int, s string) error {
	_, err := logger.Write([]byte(s))

	return err
}

func (logger *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	logger.Output(2, s)
	panic(s)
}

func (logger *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logger.Output(2, s)
	panic(s)
}

func (logger *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	logger.Output(2, s)
	panic(s)
}

func (logger *Logger) Prefix() string {
	return logger.prefix
}

func (logger *Logger) Print(v ...interface{}) {
	logger.Output(2, fmt.Sprint(v...))
}

func (logger *Logger) Printf(format string, v ...interface{}) {
	logger.Output(2, fmt.Sprintf(format, v...))
}

func (logger *Logger) Println(v ...interface{}) {
	logger.Output(2, fmt.Sprintln(v...))
}

func (logger *Logger) SetFlags(flag int) {
	logger.flag = flag
}

func (logger *Logger) SetPrefix(prefix string) {
	logger.prefix = prefix
}

// Write writes a bytes array to the Logentries TCP connection,
// it adds the access token and prefix and also replaces
// line breaks with the unicode \u2028 charachter
func (logger *Logger) Write(p []byte) (n int, err error) {
	if err := logger.ensureOpenConnection(); err != nil {
		return 0, err
	}

	logger.mu.Lock()
	defer logger.mu.Unlock()

	count := strings.Count(string(p), lineSep)
	p = []byte(strings.Replace(string(p), lineSep, "\u2028", count-1))

	logger.buf = logger.buf[:0]
	logger.buf = append(logger.buf, (logger.token + " ")...)
	logger.buf = append(logger.buf, (logger.prefix + " ")...)
	logger.buf = append(logger.buf, p...)

	return logger.conn.Write(logger.buf)
}
