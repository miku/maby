maby
====

[MAB](http://www.dnb.de/DE/Standardisierung/Formate/MAB/mab_node.html) library for Go.

```go
package main

import (
    "fmt"
    "io"
    "log"
    "os"

    "github.com/miku/maby"
)

func main() {
    r := maby.NewReader(os.Stdin)
    r.StripCollation = true
    for {
        rec, err := r.ReadRecord()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }
        if f, ok := rec.FieldByKey("331"); ok {
            fmt.Println(f.Value)
        }
    }
}
```

```shell
$ go run main.go < some.mab
Die Comédia dell Arte und Deutschland
Der Schauspieler
Taschenbuch Elektrotechnik
Systeme der Elektroenergietechnik
Bauelemente und Bausteine der Informationstechnik
Bauelemente und Bausteine der Informationstechnik
Grundlagen der Informationstechnik
```
