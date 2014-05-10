logentries
=============

Golang library for logentries.com

It is compatible with http://golang.org/pkg/log/#Logger

godoc - http://godoc.org/github.com/bsphere/logentries


Usage
-----
Add a new manual TCP token log at logentries.com and copy the token

Installation: `go get github.com/bsphere/logentries`


```go
import (
	"github.com/bsphere/logentries"
)

func main() {
	le, err := logentries.Connect("XXXX-XXXX-XXXX-XXXX") // replace with token
	if err != nil {
		panic(err)
	}

	defer le.Close()

	le.Println("another test message")
}
```
