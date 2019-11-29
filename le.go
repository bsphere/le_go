// Package le_go provides a Golang client library for logging to
// logentries.com over a TCP connection.
//
// it uses an access token for sending log events.
package le_go

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

// Logger represents a Logentries logger,
// it holds the open TCP connection, access token, prefix and flags.
//
// all Logger operations are thread safe and blocking,
// log operations can be invoked in a non-blocking way by calling them from
// a goroutine.
type Logger struct {
	conn              net.Conn
	flag              int
	mu                sync.Mutex
	writeLock         sync.Mutex
	prefix            string
	host              string
	token             string
	buf               []byte
	lastRefreshAt     time.Time
	writeTimeout      time.Duration
	_testWaitForWrite *sync.WaitGroup
}

const lineSep = "\n"
const maxLogLength int = 65000 //add 535 chars of headroom for the filename, timestamp and header
var defaultWriteTimeout = 10 * time.Second

// Connect creates a new Logger instance and opens a TCP connection to
// logentries.com,
// The token can be generated at logentries.com by adding a new log,
// choosing manual configuration and token based TCP connection.
func Connect(host, token string) (*Logger, error) {
	logger := Logger{
		host:          host,
		token:         token,
		lastRefreshAt: time.Now(),
		writeTimeout:  defaultWriteTimeout,
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
	conn, err := tls.Dial("tcp", logger.host, &tls.Config{})
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

	if time.Now().After(logger.lastRefreshAt.Add(15 * time.Minute)) {
		logger.conn.Close()
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

// Fatal is same as Print() but calls to os.Exit(1)
func (logger *Logger) Fatal(v ...interface{}) {
	logger.Output(3, fmt.Sprint(v...), handleFatalActions)
}

// Fatalf is same as Printf() but calls to os.Exit(1)
func (logger *Logger) Fatalf(format string, v ...interface{}) {
	logger.Output(3, fmt.Sprintf(format, v...), handleFatalActions)
}

// Fatalln is same as Println() but calls to os.Exit(1)
func (logger *Logger) Fatalln(v ...interface{}) {
	logger.Output(3, fmt.Sprintln(v...), handleFatalActions)
}

// Flags returns the logger flags
func (logger *Logger) Flags() int {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	return logger.flag
}

// Taken with slight modification from src/log/log.go
// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 3.
// Output does the actual writing to the TCP connection
func (l *Logger) Output(calldepth int, s string, doAsync func()) {
	defer func() {
		if re := recover(); re != nil {
			fmt.Printf("Panicked in logger.output %v\n", re)
			debug.PrintStack()
			panic(re)
		}
	}()
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(log.Lshortfile|log.Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}

	// Replace all but the trailing newline with `\u2028`
	count := strings.Count(s, lineSep)
	s = strings.Replace(s, lineSep, "\u2028", count-1)

	go func() {
		l.writeToLogEntries(s, file, now, line)
		doAsync()
	}()
}

// Panic is same as Print() but calls to panic
func (logger *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	logger.Output(3, s, handlePanicActions(s))
}

// Panicf is same as Printf() but calls to panic
func (logger *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logger.Output(3, s, handlePanicActions(s))
}

// Panicln is same as Println() but calls to panic
func (logger *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	logger.Output(3, s, handlePanicActions(s))
}

// Prefix returns the logger prefix
func (logger *Logger) Prefix() string {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	return logger.prefix
}

// Print logs a message
func (logger *Logger) Print(v ...interface{}) {
	logger.Output(3, fmt.Sprint(v...), handlePrintActions)
}

// Printf logs a formatted message
func (logger *Logger) Printf(format string, v ...interface{}) {
	logger.Output(3, fmt.Sprintf(format, v...), handlePrintActions)
}

// Println logs a message with a linebreak
func (logger *Logger) Println(v ...interface{}) {
	logger.Output(3, fmt.Sprintln(v...), handlePrintActions)
}

// SetFlags sets the logger flags
func (logger *Logger) SetFlags(flag int) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.flag = flag
}

// SetPrefix sets the logger prefix
func (logger *Logger) SetPrefix(prefix string) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.prefix = prefix
}

// Write writes a bytes array to the Logentries TCP connection,
// it adds the access token and prefix and also replaces
// line breaks with the unicode \u2028 character
func (logger *Logger) Write(p []byte) (n int, err error) {
	if err := logger.ensureOpenConnection(); err != nil {
		return 0, err
	}
	return logger.conn.Write(p)
}

// Taken wholesale from src/log/log.go
// formatHeader writes log header to buf in following order:
//   * l.prefix (if it's not blank),
//   * date and/or time (if corresponding flags are provided),
//   * file and line number (if corresponding flags are provided).
func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	*buf = append(*buf, l.prefix...)
	if l.flag&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
		if l.flag&log.LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&log.Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(log.Ltime|log.Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&log.Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flag&(log.Lshortfile|log.Llongfile) != 0 {
		if l.flag&log.Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

// Taken wholesale from src/log/log.go
// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l *Logger) writeToLogEntries(s, file string, now time.Time, line int) {
	l.writeLock.Lock()
	defer l.writeLock.Unlock()

	var i, n int
	var err error

	for i = 0; i < len(s); i = i + n {
		end := i + maxLogLength - 2
		if end > len(s) {
			end = len(s)
		}
		l.buf = l.buf[:0]
		l.buf = append(l.buf, (l.token + " ")...)
		l.formatHeader(&l.buf, now, file, line)
		l.buf = append(l.buf, s[i:end]...)
		if len(s) == 0 || s[len(s)-1] != '\n' {
			l.buf = append(l.buf, '\n')
		}
		err = l.conn.SetWriteDeadline(time.Now().Add(l.writeTimeout))
		if err != nil {
			log.Printf("le_go: Error setting write deadline: %s", err.Error())
			log.Printf("Wanted to log: %s", s)
			return
		}
		n, err = l.Write(l.buf)
		if err != nil {
			log.Printf("Error in write call: %s", err.Error())
			log.Printf("Wanted to log: %s", s)
		}

		if l._testWaitForWrite != nil {
			l._testWaitForWrite.Done()
		}
	}
}

func handleFatalActions() {
	os.Exit(1)
}

func handlePanicActions(s string) func() {
	return func() {
		panic(s)
	}
}

func handlePrintActions() {
	return
}
