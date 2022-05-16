package compress

import (
	"bytes"
	"compress/gzip"
)

type GzipCompress struct {
	Data []byte
}

func NewGzipCompress(data []byte) *GzipCompress {
	return &GzipCompress{
		Data: data,
	}
}

// GzipCompress сжимает слайс байт.
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
