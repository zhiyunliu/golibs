package xzip

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
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

func TestCompressFormats(t *testing.T) {
	want := []byte("hello, formats\x00\xff")

	for _, format := range []Format{FormatZIP, FormatZlib} {
		compressed, err := Compress(want, format)
		if err != nil {
			t.Fatalf("Compress(%d) error = %v", format, err)
		}

		switch format {
		case FormatZIP:
			archive, err := zip.NewReader(bytes.NewReader(compressed), int64(len(compressed)))
			if err != nil || len(archive.File) != 1 || archive.File[0].Name != "location" {
				t.Fatalf("invalid ZIP archive: files = %v, error = %v", archive, err)
			}
		case FormatZlib:
			reader, err := zlib.NewReader(bytes.NewReader(compressed))
			if err != nil {
				t.Fatalf("zlib.NewReader() error = %v", err)
			}
			decoded, readErr := io.ReadAll(reader)
			closeErr := reader.Close()
			if readErr != nil || closeErr != nil || !bytes.Equal(decoded, want) {
				t.Fatalf("zlib stream: decoded = %v, read error = %v, close error = %v", decoded, readErr, closeErr)
			}
		}

		got, err := Decompress(compressed, format)
		if err != nil || !bytes.Equal(got, want) {
			t.Fatalf("Decompress(%d) = %v, %v; want %v", format, got, err, want)
		}
	}
}

func TestCompressLevelZlib(t *testing.T) {
	want := []byte("compressed with level")
	compressed, err := CompressLevel(want, FormatZlib, BestCompression)
	if err != nil {
		t.Fatalf("CompressLevel() error = %v", err)
	}
	got, err := Decompress(compressed, FormatZlib)
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("Decompress() = %v, %v; want %v", got, err, want)
	}
	if _, err := CompressLevel(want, FormatZlib, 10); err == nil {
		t.Fatal("CompressLevel() accepted invalid level")
	}
}

func TestInvalidFormat(t *testing.T) {
	if _, err := Compress([]byte("hello"), Format(99)); err == nil {
		t.Fatal("Compress() accepted invalid format")
	}
	if _, err := Decompress([]byte("hello"), Format(99)); err == nil {
		t.Fatal("Decompress() accepted invalid format")
	}
	if _, err := Decompress([]byte("not a zlib stream"), FormatZlib); err == nil {
		t.Fatal("Decompress() accepted invalid zlib stream")
	}
}
