package xzip

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"testing"
)

func TestGzipCompressAndDecompress(t *testing.T) {
	testCases := [][]byte{
		nil,
		{},
		[]byte("hello, gzip\x00\xff"),
	}

	for _, want := range testCases {
		compressed, err := GzipCompress(want)
		if err != nil {
			t.Fatalf("GzipCompress() error = %v", err)
		}
		got, err := GzipDecompress(compressed)
		if err != nil {
			t.Fatalf("GzipDecompress() error = %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("GzipDecompress() = %v, want %v", got, want)
		}
	}
}

func TestGzipDecompressInvalidData(t *testing.T) {
	if _, err := GzipDecompress([]byte("not gzip data")); err == nil {
		t.Fatal("GzipDecompress() error = nil, want non-nil")
	}
}

func TestGzipCompressLevel(t *testing.T) {
	want := []byte("hello, gzip level")
	levels := []int{
		gzip.NoCompression,
		gzip.BestSpeed,
		gzip.DefaultCompression,
		gzip.BestCompression,
		gzip.HuffmanOnly,
	}

	for _, level := range levels {
		compressed, err := GzipCompressLevel(want, Level(level))
		if err != nil {
			t.Fatalf("GzipCompressLevel(%d) error = %v", level, err)
		}
		got, err := GzipDecompress(compressed)
		if err != nil {
			t.Fatalf("GzipDecompress() error = %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("GzipDecompress() = %v, want %v", got, want)
		}
	}

	if _, err := GzipCompressLevel(want, 10); err == nil {
		t.Fatal("GzipCompressLevel() error = nil, want non-nil")
	}
}

func TestZipCompressAndDecompress(t *testing.T) {
	testCases := [][]byte{
		nil,
		{},
		[]byte("hello, zip\x00\xff"),
	}

	for _, want := range testCases {
		compressed, err := ZipCompress(want)
		if err != nil {
			t.Fatalf("ZipCompress() error = %v", err)
		}
		got, err := ZipDecompress(compressed)
		if err != nil {
			t.Fatalf("ZipDecompress() error = %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("ZipDecompress() = %v, want %v", got, want)
		}
	}
}

func TestZipDecompressInvalidData(t *testing.T) {
	if _, err := ZipDecompress([]byte("not zip data")); err == nil {
		t.Fatal("ZipDecompress() error = nil, want non-nil")
	}
}

func TestZipCompressLevel(t *testing.T) {
	want := []byte("hello, zip level")
	levels := []int{
		flate.NoCompression,
		flate.BestSpeed,
		flate.DefaultCompression,
		flate.BestCompression,
		flate.HuffmanOnly,
	}

	for _, level := range levels {
		compressed, err := ZipCompressLevel(want, Level(level))
		if err != nil {
			t.Fatalf("ZipCompressLevel(%d) error = %v", level, err)
		}
		got, err := ZipDecompress(compressed)
		if err != nil {
			t.Fatalf("ZipDecompress() error = %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("ZipDecompress() = %v, want %v", got, want)
		}
	}

	if _, err := ZipCompressLevel(want, 10); err == nil {
		t.Fatal("ZipCompressLevel() error = nil, want non-nil")
	}
}

func TestZipDecompressWithoutFiles(t *testing.T) {
	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)
	if _, err := writer.Create("directory/"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	if _, err := ZipDecompress(buf.Bytes()); err == nil {
		t.Fatal("ZipDecompress() error = nil, want non-nil")
	}
}
