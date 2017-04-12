package maby

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

const (
	// GS is the group separator.
	GS = 0x1D
	// RS is the record separator.
	RS = 0x1E
	// US is the unit separator.
	US = 0x1F
	// LeaderSize is the length of the header.
	LeaderSize = 24
	// FieldnameSize is the length of a field name.
	FieldnameSize = 4
	// IdentifierTag is the name of the identifier field.
	IdentifierTag = "001"
	// TagLength is the length of a tag.
	TagLength = 4
)

// Record is a MAB record.
type Record struct {
	Leader     string  `json:"leader"`
	Type       string  `json:"type"`
	Identifier string  `json:"id"`
	Fields     []Field `json:"fields"`
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

// Field is a single field.
type Field struct {
	Key   string `json:"k"`
	Value string `json:"v"`
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

// readRecord parses a single record.
func (r *Reader) readRecord(p []byte) (*Record, error) {
	var record Record
	if len(p) < LeaderSize {
		return nil, fmt.Errorf("missing or short header: %d", len(p))
	}
	record.Leader = string(p[:LeaderSize])
	record.Type = string(record.Leader[LeaderSize-1])

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
			field := Field{Key: name, Value: string(v)}
			if name == IdentifierTag {
				record.Identifier = string(v)
			}
			record.Fields = append(record.Fields, field)
		}
	}
	return &record, nil
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
