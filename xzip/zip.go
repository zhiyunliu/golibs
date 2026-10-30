package xzip

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"errors"
	"io"
)

type Level int

const (
	NoCompression      Level = flate.NoCompression
	BestSpeed          Level = flate.BestSpeed
	BestCompression    Level = flate.BestCompression
	DefaultCompression Level = flate.DefaultCompression
	HuffmanOnly        Level = flate.HuffmanOnly
)

// ZipCompress 把裸数据打包成含单个 "location" 条目的 zip（模拟分发格式）。
func ZipCompress(raw []byte) (result []byte, err error) {
	return ZipCompressLevel(raw, flate.DefaultCompression)
}

// ZipCompressLevel 使用指定等级把裸数据打包成含单个 "location" 条目的 zip。
func ZipCompressLevel(raw []byte, level Level) (result []byte, err error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	zw.RegisterCompressor(zip.Deflate, func(writer io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(writer, int(level))
	})
	w, err := zw.Create("location")
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(raw); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ZipDecompress 把 zip 解包。
func ZipDecompress(raw []byte) (result []byte, err error) {
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return nil, err
	}

	for _, file := range zr.File {
		if file.FileInfo().IsDir() {
			continue
		}

		reader, err := file.Open()
		if err != nil {
			return nil, err
		}
		result, readErr := io.ReadAll(reader)
		closeErr := reader.Close()
		if readErr != nil {
			return nil, readErr
		}
		if closeErr != nil {
			return nil, closeErr
		}
		return result, nil
	}

	return nil, errors.New("zip archive contains no files")
}
