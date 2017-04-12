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

// Version of application.
const Version = "0.1.0"

var (
	stripCollection = flag.Bool("strip", false, "strip collation bytes (88, 89)")
	version         = flag.Bool("version", false, "show version")
)

func main() {
	flag.Parse()
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
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
