maby
====

[MAB](http://www.dnb.de/DE/Standardisierung/Formate/MAB/mab_node.html) library for Go.

Jsonify
-------

```shell
$ mabread < fixtures/a.mab | head -1 | jq .
{
  "leader": "00633nM2.01200024      h",
  "type": "h",
  "id": "HF0000005",
  "fields": [
    {
      "k": "001",
      "v": "HF0000005"
    },
    {
      "k": "002a",
      "v": "20000517"
    },
    {
      "k": "003",
      "v": "20010113"
    },
    {
      "k": "004",
      "v": "20161129"
    },
    {
      "k": "030",
      "v": "az1ddzzzc||37"
    },
    {
      "k": "050",
      "v": "a|||||||||||||"
    },
    {
      "k": "051",
      "v": "s|||||||"
    },
    {
      "k": "100b",
      "v": "Brauneck, Manfred$$9HF0000002$$b[Hrsg.]"
    },
    {
      "k": "331",
      "v": "Theaterlexikon"
    },
    {
      "k": "335",
      "v": "Begriffe und Epochen, Bühnen und Ensembles"
    },
    {
      "k": "359",
      "v": "Hrsg.: Manfred Brauneck ; Gerard Schneiling"
    },
    {
      "k": "403",
      "v": "3., vollst. überarb. u. erw. Ausg."
    },
    {
      "k": "410",
      "v": "Reinbek bei Hamburg"
    },
    {
      "k": "412",
      "v": "Rowohlt"
    },
    {
      "k": "425",
      "v": "1992"
    },
    {
      "k": "433",
      "v": "1137 S."
    },
    {
      "k": "434",
      "v": "Ill."
    },
    {
      "k": "451",
      "v": "Rowohlts Enzyklopädie ; 465"
    },
    {
      "k": "453m",
      "v": "HF0000006"
    },
    {
      "k": "454c",
      "v": "Rowohlts Enzyklopädie"
    },
    {
      "k": "455",
      "v": "465"
    },
    {
      "k": "501",
      "v": "Literaturverz. S. 1133-1137"
    },
    {
      "k": "540a",
      "v": "3-499-55465-8"
    },
    {
      "k": "710",
      "v": "Theaterwissenschaft/M$$9HF0000122"
    },
    {
      "k": "710",
      "v": "Lexika/U$$9HF0000110"
    },
    {
      "k": "104b",
      "v": "Schneiling, Gerard$$9HF0000003$$b[Hrsg.]"
    }
  ]
}
```

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
