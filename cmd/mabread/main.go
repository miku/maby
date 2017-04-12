package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/miku/maby"
)

var stripCollection = flag.Bool("strip", false, "strip collation bytes (88, 89)")

func main() {
	flag.Parse()
	r := maby.NewReader(os.Stdin)
	r.StripCollation = *stripCollection
	for {
		rec, err := r.ReadRecord()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		b, err := json.Marshal(rec)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}
}
