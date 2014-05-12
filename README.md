logentries
=============

Golang library for logentries.com

It is compatible with http://golang.org/pkg/log/#Logger

[![GoDoc](https://godoc.org/github.com/bsphere/le_go?status.png)](https://godoc.org/github.com/bsphere/le_go)

Usage
-----
Add a new manual TCP token log at [logentries.com](https://logentries.com/quick-start/) and copy the [token](https://logentries.com/doc/input-token/).

Installation: `go get github.com/bsphere/le_go`

**Note:** The Logger is blocking, it can be easily run in a goroutine by calling `go le.Println(...)`

```go
import (
	"github.com/bsphere/le_go"

func main() {
	le, err := le_go.Connect("XXXX-XXXX-XXXX-XXXX") // replace with token
	if err != nil {
		panic(err)
	}

	defer le.Close()

	le.Println("another test message")
}
```
