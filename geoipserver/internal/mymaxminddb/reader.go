package mymaxminddb

import (
	"fmt"
	"net"
	"time"

	"github.com/oschwald/maxminddb-golang"
)

const timestampFormat = time.RFC3339

type Reader struct {
	maxminddb.Reader
	Metadata Metadata
}

type Metadata struct {
	filename string
	*maxminddb.Metadata
}

func (m *Metadata) FormatVersion() string {
	return fmt.Sprintf("%d.%d", m.BinaryFormatMajorVersion, m.BinaryFormatMinorVersion)
}

func (m *Metadata) BuildTimestamp() string {
	return time.Unix(int64(m.BuildEpoch), 0).Format(timestampFormat)
}

func Open(filename string) (*Reader, error) {
	db, err := maxminddb.Open(filename)
	if err != nil {
		return nil, err
	}
	dbWrapper := &Reader{*db, Metadata{filename, &db.Metadata}}
	return dbWrapper, err
}

func (r *Reader) GetRecord(ip net.IP) (interface{}, error) {
	var record interface{}
	err := r.Lookup(ip, &record)
	return record, err
}
