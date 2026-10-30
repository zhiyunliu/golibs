package xzip

import (
	"bytes"
	"compress/gzip"
	"io"
)

// GzipCompress 压缩
func GzipCompress(raw []byte) (result []byte, err error) {
	return GzipCompressLevel(raw, gzip.DefaultCompression)
}

// GzipCompressLevel 使用指定等级压缩。
func GzipCompressLevel(raw []byte, level Level) (result []byte, err error) {
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, int(level))
	if err != nil {
		return nil, err
	}
	if _, err := zw.Write(raw); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GzipDecompress(raw []byte) (result []byte, err error) {
	zr, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	return io.ReadAll(zr)
}
