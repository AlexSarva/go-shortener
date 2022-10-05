package compress

import (
	"bytes"
	"compress/gzip"
)

// GzipCompress compress data into gzip format
type GzipCompress struct {
	Data []byte
}

// NewGzipCompress initializer of GzipCompress struct
func NewGzipCompress(data []byte) *GzipCompress {
	return &GzipCompress{
		Data: data,
	}
}

// Compress GzipCompress compress slice of bites.
func (g GzipCompress) Compress() []byte {
	var b bytes.Buffer

	w, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		return nil
	}

	_, err = w.Write(g.Data)
	if err != nil {
		return nil
	}

	err = w.Close()
	if err != nil {
		return nil
	}

	return b.Bytes()
}
