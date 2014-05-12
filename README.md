logentries
=============

Golang library for logentries.com

It is compatible with http://golang.org/pkg/log/#Logger

godoc - http://godoc.org/github.com/bsphere/logentries


Usage
-----
Add a new manual TCP token log at [logentries.com](https://logentries.com/quick-start/) and copy the [token](https://logentries.com/doc/input-token/).

Installation: `go get github.com/bsphere/logentries`

**Note:** The Logger is asynchronous, make sure your application does not terminate immediately otherwise no log message will be sent.

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
