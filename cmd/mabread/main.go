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
	// r.StripCollation = true
	// for {
	// 	rec, err := r.ReadRecord()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	if f, ok := rec.FieldByKey("331"); ok {
	// 		fmt.Println(f.Value)
	// 	}
	// }
}
