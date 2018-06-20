package api

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// it is testing the MergeToBuf function
func TestMergeToeBuf(t *testing.T) {
	test1 := filepath.Join(os.Getenv("GOPATH"), "src/github.com/charleswklau/pdfcpu/pkg/api/test", "test1.pdf")
	test2 := filepath.Join(os.Getenv("GOPATH"), "src/github.com/charleswklau/pdfcpu/pkg/api/test", "test2.pdf")

	b, err := MergeToBuf([]string{test1, test2}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err = testWriteIntoDisk(b); err != nil {
		t.Fatal(err)
	}
}

// it is testing the MergeFileToBuf function
func TestMergeFilesToBuf(t *testing.T) {
	test1 := filepath.Join(os.Getenv("GOPATH"), "src/github.com/charleswklau/pdfcpu/pkg/api/test", "test1.pdf")
	test2 := filepath.Join(os.Getenv("GOPATH"), "src/github.com/charleswklau/pdfcpu/pkg/api/test", "test2.pdf")

	file1, err := os.OpenFile(test1, os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer file1.Close()

	file2, err := os.OpenFile(test2, os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer file2.Close()

	b, err := MergeFileToBuf([]*os.File{file1, file2}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err = testWriteIntoDisk(b); err != nil {
		t.Fatal(err)
	}
}

// helper function write the buffer into the disk
func testWriteIntoDisk(buf *bytes.Buffer) error {

	outputFile := filepath.Join(os.Getenv("GOPATH"), "src/github.com/charleswklau/pdfcpu/pkg/api/test", "test-out.pdf")

	// remove file when it is already exist
	f, err := os.Stat(outputFile)
	if f != nil {
		if err := os.Remove(outputFile); err != nil {
			return err
		}
	}

	// write file
	fp, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	// close later and ignore errors on close
	defer fp.Close()

	if _, err := fp.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}
