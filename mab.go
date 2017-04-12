package maby

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	GS            = 0x1D
	RS            = 0x1E
	US            = 0x1F
	LeaderSize    = 24
	FieldnameSize = 4
	IdentifierTag = "001"
)

// Record is a MAB record.
type Record struct {
	Leader          string  `json:"leader"`
	Length          int     `json:"len"`
	Status          string  `json:"status"`
	Version         string  `json:"version"`
	IndicatorLength int     `json:"ilen"`
	TagLength       int     `json:"tlen"`
	Offset          int     `json:"offset"`
	Type            string  `json:"type"`
	Identifier      string  `json:"id"`
	Fields          []Field `json:"fields"`
}

// Field is a single field.
type Field struct {
	Key       string `json:"k"`
	Value     string `json:"v"`
	Subfields []struct {
		Key   string `json:"k"`
		Value string `json:"v"`
	} `json:"s,omitempty"`
}

// Reader reads MAB records.
type Reader struct {
	StripCollation bool
	r              *bufio.Reader
	readerDone     bool
}

// NewReader creates a new record reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

// ReadRecord reads the next record. Returns io.EOF if there are no more records.
func (r *Reader) ReadRecord() (*Record, error) {
	if r.readerDone {
		return nil, io.EOF
	}
	b, err := r.r.ReadBytes(GS)
	if err == io.EOF {
		if len(b) == 0 {
			return nil, err
		}
		r.readerDone = true
		return r.readRecord(b)
	}
	if err != nil {
		return nil, err
	}
	return r.readRecord(b)
}

// ReadRecords returns all records at once.
func (r *Reader) ReadRecords() (records []*Record, err error) {
	for {
		rec, err := r.ReadRecord()
		if err == io.EOF {
			return records, nil
		}
		if err != nil {
			return records, err
		}
		records = append(records, rec)
	}
}

// readRecord parses a single record.
func (r *Reader) readRecord(p []byte) (record *Record, err error) {
	if len(p) < LeaderSize {
		return nil, fmt.Errorf("missing or short header: %d", len(p))
	}

	record = new(Record)

	// Setup meta fields.
	record.Leader = string(p[:LeaderSize])

	if record.Length, err = strconv.Atoi(record.Leader[0:5]); err != nil {
		return nil, err
	}
	record.Status = string(record.Leader[5])
	record.Version = record.Leader[6:10]
	if record.IndicatorLength, err = strconv.Atoi(string(record.Leader[10])); err != nil {
		return nil, err
	}
	if record.TagLength, err = strconv.Atoi(string(record.Leader[11])); err != nil {
		return nil, err
	}
	if record.Offset, err = strconv.Atoi(record.Leader[12:17]); err != nil {
		return nil, err
	}
	record.Type = string(record.Leader[LeaderSize-1])

	// Fields.
	for _, part := range bytes.Split(p[LeaderSize:], []byte{RS}) {
		if len(part) < FieldnameSize {
			continue
		}
		name := string(bytes.TrimSpace(part[0:FieldnameSize]))
		for _, v := range bytes.Split(part[FieldnameSize:], []byte{US}) {
			if r.StripCollation {
				v = bytes.Replace(v, []byte{0x88}, []byte{}, -1)
				v = bytes.Replace(v, []byte{0x89}, []byte{}, -1)
			}
			if record.Identifier == "" && name == IdentifierTag {
				record.Identifier = string(v)
			}

			pp := strings.Split(string(v), "$$")

			switch {
			case len(pp) == 1:
				record.Fields = append(record.Fields, Field{
					Key:   name,
					Value: string(v),
				})
			case len(pp) > 1:
				field := Field{Key: name, Value: pp[0]}
				for _, sub := range pp[1:] {
					if len(sub) < 2 {
						continue
					}
					s := struct {
						Key   string `json:"k"`
						Value string `json:"v"`
					}{
						string(sub[0]),
						sub[1:],
					}
					field.Subfields = append(field.Subfields, s)
				}
				record.Fields = append(record.Fields, field)
			}
		}
	}
	return record, nil
}

// FieldsByKey returns all fields matching a number of given keys.
func (r *Record) FieldsByKey(keys ...string) (fields []Field) {
	for _, f := range r.Fields {
		for _, k := range keys {
			if k == f.Key {
				fields = append(fields, f)
			}
		}
	}
	return
}

// FieldByKey returns the first field matching a given key and a bool, if that
// field was found.
func (r *Record) FieldByKey(key string) (field Field, found bool) {
	for _, f := range r.Fields {
		if key == f.Key {
			return f, true
		}
	}
	return Field{}, false
}
