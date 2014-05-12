logentries
=============

Golang library for logentries.com

It is compatible with http://golang.org/pkg/log/#Logger

godoc - http://godoc.org/github.com/bsphere/le_go


Usage
-----
Add a new manual TCP token log at logentries.com and copy the token

Installation: `go get github.com/bsphere/le_go`


```go
import (
	"github.com/bsphere/le_go"

func main() {
	le, err := logentries.Connect("XXXX-XXXX-XXXX-XXXX") // replace with token
	if err != nil {
		panic(err)
	}

	defer le.Close()

	le.Println("another test message")
}
```
