package xfile

import (
	"path/filepath"
	"testing"
)

func TestGetNameWithoutExt(t *testing.T) {
	// Testing filename with extension
	filename1 := "example.txt"
	result1 := GetNameWithoutExt(filename1)
	if result1 != "example" {
		t.Errorf("Expected 'example', but got %s", result1)
	}

	// Testing filename without extension
	filename2 := "document"
	result2 := GetNameWithoutExt(filename2)
	if result2 != "document" {
		t.Errorf("Expected 'document', but got %s", result2)
	}

	// Testing filename starting with dot
	filename3 := ".hiddenfile"
	result3 := GetNameWithoutExt(filename3)
	if result3 != "" {
		t.Errorf("Expected '.hiddenfile', but got %s", result3)
	}

	// Testing filename without path
	filename4 := "folder/file"
	result4 := GetNameWithoutExt(filename4)
	if result4 != "file" {
		t.Errorf("Expected 'file', but got %s", result4)
	}

	// Testing filename with multiple dots
	filename5 := "archive.tar.gz"
	result5 := GetNameWithoutExt(filename5)
	if result5 != "archive.tar" {
		t.Errorf("Expected 'archive.tar', but got %s", result5)
	}
}

func TestGetAbs(t *testing.T) {
	// Testing with an existing file path
	existingPath := "testdata/example.txt"
	result, err := GetAbs(existingPath)
	expectedAbsPath, _ := filepath.Abs(existingPath)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if result != expectedAbsPath {
		t.Errorf("Expected absolute path %s, but got %s", expectedAbsPath, result)
	}

	// Testing with a non-existing file path
	nonExistingPath := "nonexistent/file.txt"
	abspath, err := GetAbs(nonExistingPath)
	if err != nil {
		t.Error("Expected an error for non-existing path, but no error was returned")
	}
	if len(abspath) <= len(nonExistingPath) {
		t.Error("Expected an error for non-existing path, but no error was returned")
	}
}

func TestExists(t *testing.T) {
	// Testing with an existing file path
	existingPath := "file.go"
	exists := Exists(existingPath)
	if !exists {
		t.Error("Expected true, but got false")
	}

	noneExistingPath := "nonexistent/file.go"
	noneExists := Exists(noneExistingPath)
	if noneExists {
		t.Error("Expected false, but got true")
	}
}
