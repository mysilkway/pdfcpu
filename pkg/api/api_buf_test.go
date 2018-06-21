package api

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	pkgDir string
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	pkgDir = path.Dir(filename)
}

// it is testing the MergeToBuf function
func TestMergeToBuf(t *testing.T) {
	test1 := filepath.Join(pkgDir, "testdata/test1.pdf")
	test2 := filepath.Join(pkgDir, "testdata/test2.pdf")

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
	test1 := filepath.Join(pkgDir, "testdata/test1.pdf")
	test2 := filepath.Join(pkgDir, "testdata/test2.pdf")

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
	outputFile := filepath.Join(pkgDir, "testdata/test-out.pdf")

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
