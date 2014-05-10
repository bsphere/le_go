package logentries

import (
	"testing"
)

func TestConnectOpensConnection(t *testing.T) {
	le, err := Connect("")

	if err != nil {
		t.Error(err)
	}

	if le.conn == nil {
		t.Fail()
	}

	if le.isOpenConnection() == false {
		t.Fail()
	}
}

func TestConnectSetsToken(t *testing.T) {
	le, _ := Connect("myToken")

	if le.token != "myToken" {
		t.Fail()
	}
}

func TestCloseClosesConnection(t *testing.T) {
	le, _ := Connect("")

	le.Close()

	if le.isOpenConnection() == true {
		t.Fail()
	}
}

func TestReopenConnectionOpensConnection(t *testing.T) {
	le, _ := Connect("")

	le.Close()

	le.reopenConnection()

	if le.isOpenConnection() == false {
		t.Fail()
	}
}

func TestEnsureOpenConnectionDoesNothingOnOpenConnection(t *testing.T) {
	le, _ := Connect("")

	old := &le.conn

	le.reopenConnection()

	if old != &le.conn {
		t.Fail()
	}
}

func TestEnsureOpenConnectionCreatesNewConnection(t *testing.T) {
	le, _ := Connect("")

	le.Close()

	le.reopenConnection()

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
