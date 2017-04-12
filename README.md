maby
====

[MAB](http://www.dnb.de/DE/Standardisierung/Formate/MAB/mab_node.html) library for Go.

Show titles
-----------

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
$ go run cmd/mabread/main.go < fixtures/a.mab
Die Comédia dell Arte und Deutschland
Der Schauspieler
Taschenbuch Elektrotechnik
Systeme der Elektroenergietechnik
Bauelemente und Bausteine der Informationstechnik
Bauelemente und Bausteine der Informationstechnik
Grundlagen der Informationstechnik
```

Count record
------------

```go
package main

import (
    "log"
    "os"

    "github.com/miku/maby"
)

func main() {
    r := maby.NewReader(os.Stdin)
    records, err := r.ReadRecords()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("%d records", len(records))
}
```

```shell
$ go run cmd/mabread/main.go < fixtures/a.mab
2017/04/12 16:34:50 122064 records
```
