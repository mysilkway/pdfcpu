package api

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteBuf(t *testing.T) {
	test1 := filepath.Join(os.Getenv("GOPATH"), "src/github.com/charleswklau/pdfcpu/pkg/api/test", "test1.pdf")
	test2 := filepath.Join(os.Getenv("GOPATH"), "src/github.com/charleswklau/pdfcpu/pkg/api/test", "test2.pdf")

	b, err := MergeAsBuf([]string{test1, test2}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err = testWriteIntoDisk(b); err != nil {
		t.Fatal(err)
	}
}

func testWriteIntoDisk(buf *bytes.Buffer) error {

	outputFile := filepath.Join(os.Getenv("GOPATH"), "src/github.com/charleswklau/pdfcpu/pkg/api/test", "test-out.pdf")

	f, err := os.Stat(outputFile)
	if f != nil {
		if err := os.Remove(outputFile); err != nil{
			return err
		}
	}

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