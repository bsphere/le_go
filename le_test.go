package le_go

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestConnectOpensConnection(t *testing.T) {
	le, err := Connect("")
	if err != nil {
		t.Fatal(err)
	}

	defer le.Close()

	if le.conn == nil {
		t.Fail()
	}
}

func TestConnectSetsToken(t *testing.T) {
	le, err := Connect("myToken")
	if err != nil {
		t.Fatal(err)
	}

	defer le.Close()

	if le.token != "myToken" {
		t.Fail()
	}
}

func TestWriteReopensConnection(t *testing.T) {
	le, err := Connect("")
	if err != nil {
		t.Fatal(err)
	}

	oldConn := le.conn
	err = le.conn.Close()
	if err != nil {
		t.Fail()
	}

	written, err := le.Write([]byte("should reopen"))
	if written != 13 || err != nil {
		t.Error(written, err)
	}

	if le.conn == oldConn {
		t.Fail()
	}
}

func TestCloseClosesConnection(t *testing.T) {
	le, err := Connect("")
	if err != nil {
		t.Fatal(err)
	}

	err = le.Close()
	if err != nil {
		t.Fail()
	}

	if le.conn != nil {
		t.Fail()
	}

	err = le.Close()
	if err != ErrClosed {
		t.Fail()
	}

	written, err := le.Write([]byte("write after close"))
	if written != 0 || err != ErrClosed {
		t.Error(written, err)
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
	le, err := Connect("myToken")
	if err != nil {
		t.Fatal(err)
	}

	defer le.Close()

	// the test will fail to compile if Logger doesn't implement io.Writer
	func(w io.Writer) {}(le)
}

func TestReplaceNewline(t *testing.T) {
	le, err := Connect("myToken")
	if err != nil {
		t.Fatal(err)
	}

	defer le.Close()

	buf := le.makeBuf([]byte("1\n2\n3"))

	if strings.Count(string(buf), "\u2028") != 2 {
		t.Fail()
	}
}

func TestAddNewline(t *testing.T) {
	le, err := Connect("myToken")
	if err != nil {
		t.Fatal(err)
	}

	defer le.Close()

	buf := le.makeBuf([]byte("123"))

	if !strings.HasSuffix(string(buf), "\n") {
		t.Fail()
	}
}

func ExampleLogger() {
	le, err := Connect("XXXX-XXXX-XXXX-XXXX") // replace with token
	if err != nil {
		panic(err)
	}

	defer le.Close()

	le.Println("another test message")
}

func ExampleLogger_write() {
	le, err := Connect("XXXX-XXXX-XXXX-XXXX") // replace with token
	if err != nil {
		panic(err)
	}

	defer le.Close()

	fmt.Fprintln(le, "another test message")
}

func BenchmarkMakeBuf(b *testing.B) {
	le := Logger{token: "token"}

	for i := 0; i < b.N; i++ {
		le.makeBuf([]byte("test\nstring\n"))
	}
}

func BenchmarkMakeBufWithoutNewlineSuffix(b *testing.B) {
	le := Logger{token: "token"}

	for i := 0; i < b.N; i++ {
		le.makeBuf([]byte("test\nstring"))
	}
}

func BenchmarkMakeBufWithPrefix(b *testing.B) {
	le := Logger{token: "token"}
	le.SetPrefix("prefix")

	for i := 0; i < b.N; i++ {
		le.makeBuf([]byte("test\nstring\n"))
	}
}
