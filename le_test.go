package le_go

import (
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestConnectOpensConnection(t *testing.T) {
	le, err := Connect("data.logentries.com:443", "")
	if err != nil {
		t.Fatal(err)
	}

	defer le.Close()

	if le.conn == nil {
		t.Fail()
	}

	if le.isOpenConnection() == false {
		t.Fail()
	}
}

func TestConnectSetsToken(t *testing.T) {
	le, err := Connect("data.logentries.com:443", "myToken")
	if err != nil {
		t.Fatal(err)
	}

	defer le.Close()

	if le.token != "myToken" {
		t.Fail()
	}
}

func TestCloseClosesConnection(t *testing.T) {
	le, err := Connect("data.logentries.com:443", "")
	if err != nil {
		t.Fatal(err)
	}

	le.Close()

	if le.isOpenConnection() == true {
		t.Fail()
	}
}

func TestOpenConnectionOpensConnection(t *testing.T) {
	le, err := Connect("data.logentries.com:443", "")
	if err != nil {
		t.Fatal(err)
	}

	defer le.Close()

	le.openConnection()

	if le.isOpenConnection() == false {
		t.Fail()
	}
}

func TestEnsureOpenConnectionDoesNothingOnOpenConnection(t *testing.T) {
	le, err := Connect("data.logentries.com:443", "")
	if err != nil {
		t.Fatal(err)
	}

	defer le.Close()
	old := &le.conn

	le.openConnection()

	if old != &le.conn {
		t.Fail()
	}
}

func TestEnsureOpenConnectionCreatesNewConnection(t *testing.T) {
	le, err := Connect("data.logentries.com:443", "")
	if err != nil {
		t.Fatal(err)
	}

	defer le.Close()

	le.openConnection()

	if le.isOpenConnection() == false {
		t.Fail()
	}
}

func TestFlagsReturnsFlag(t *testing.T) {
	le := Logger{flag: 2}

	if le.Flags() != 2 {
		t.Fail()
	}
}

func TestSetFlagsSetsFlag(t *testing.T) {
	le := Logger{flag: 2}

	le.SetFlags(1)

	if le.flag != 1 {
		t.Fail()
	}
}

func TestPrefixReturnsPrefix(t *testing.T) {
	le := Logger{prefix: "myPrefix"}

	if le.Prefix() != "myPrefix" {
		t.Fail()
	}
}

func TestSetPrefixSetsPrefix(t *testing.T) {
	le := Logger{prefix: "myPrefix"}

	le.SetPrefix("myNewPrefix")

	if le.prefix != "myNewPrefix" {
		t.Fail()
	}
}

func TestLoggerImplementsWriterInterface(t *testing.T) {
	le, err := Connect("data.logentries.com:443", "myToken")
	if err != nil {
		t.Fatal(err)
	}

	defer le.Close()

	// the test will fail to compile if Logger doesn't implement io.Writer
	func(w io.Writer) {}(le)
}

func TestReplaceNewline(t *testing.T) {
	le, err := Connect("data.logentries.com:443", "myToken")
	if err != nil {
		t.Fatal(err)
	}

	defer le.Close()

	le.Println("1\n2\n3")

	<-time.After(10 * time.Millisecond)
	if strings.Count(string(le.buf), "\u2028") != 2 {
		t.Fail()
	}
}

func TestAddNewline(t *testing.T) {
	le, err := Connect("data.logentries.com:443", "myToken")
	if err != nil {
		t.Fatal(err)
	}

	defer le.Close()

	le.Print("123")

	<-time.After(10 * time.Millisecond)

	if !strings.HasSuffix(string(le.buf), "\n") {
		t.Fail()
	}

	le.Printf("%s", "123")

	<-time.After(10 * time.Millisecond)

	if !strings.HasSuffix(string(le.buf), "\n") {
		t.Fail()
	}
}

func TestCanSendMoreThan64k(t *testing.T) {
	le, err := Connect("data.logentries.com:443", "myToken")
	if err != nil {
		t.Fatal(err)
	}

	defer le.Close()

	longBytes := make([]byte, 140000)
	for i := 0; i < 140000; i++ {
		longBytes[i] = 'a'
	}
	longString := string(longBytes)
	// Fake the connection so we can hear about it
	fakeConn := fakeConnection{}
	le.conn = &fakeConn
	le.Print(longString)

	<-time.After(10 * time.Millisecond)
	if fakeConn.WriteCalls < 2 {
		t.Fail()
	}
}

type fakeConnection struct {
	WriteCalls int
}

func (f *fakeConnection) Write(b []byte) (int, error) {
	f.WriteCalls++
	return len(b), nil
}

func (*fakeConnection) Read(b []byte) (int, error) {
	return len(b), &fakeError{}
}

func (*fakeConnection) SetReadDeadline(time.Time) error { return nil }

func (*fakeConnection) Close() error                       { return nil }
func (*fakeConnection) LocalAddr() net.Addr                { return &fakeAddr{} }
func (*fakeConnection) RemoteAddr() net.Addr               { return &fakeAddr{} }
func (*fakeConnection) SetDeadline(t time.Time) error      { return nil }
func (*fakeConnection) SetWriteDeadline(t time.Time) error { return nil }

type fakeError struct{}

func (*fakeError) Error() string {
	return "fake network error"
}

func (*fakeError) Timeout() bool {
	return true
}

func (*fakeError) Temporary() bool {
	return true
}

type fakeAddr struct{}

func (f *fakeAddr) Network() string { return "" }
func (f *fakeAddr) String() string  { return "" }

func ExampleLogger() {
	le, err := Connect("data.logentries.com:443", "XXXX-XXXX-XXXX-XXXX") // replace with token
	if err != nil {
		panic(err)
	}

	defer le.Close()

	le.Println("another test message")
}

func ExampleLogger_write() {
	le, err := Connect("data.logentries.com:443", "XXXX-XXXX-XXXX-XXXX") // replace with token
	if err != nil {
		panic(err)
	}

	defer le.Close()

	fmt.Fprintln(le, "another test message")
}

// func BenchmarkMakeBuf(b *testing.B) {
// 	le := Logger{token: "token"}

// 	for i := 0; i < b.N; i++ {
// 		le.makeBuf([]byte("test\nstring\n"))
// 	}
// }

// func BenchmarkMakeBufWithoutNewlineSuffix(b *testing.B) {
// 	le := Logger{token: "token"}

// 	for i := 0; i < b.N; i++ {
// 		le.makeBuf([]byte("test\nstring"))
// 	}
// }

// func BenchmarkMakeBufWithPrefix(b *testing.B) {
// 	le := Logger{token: "token"}
// 	le.SetPrefix("prefix")

// 	for i := 0; i < b.N; i++ {
// 		le.makeBuf([]byte("test\nstring\n"))
// 	}
// }
