package pdfcpu

import (
	"bytes"
	"github.com/mysilkway/pdfcpu/pkg/log"
	"github.com/pkg/errors"
	"os"
)

// ParseFileToContext parses the File and generates a Context, an in-memory representation containing a cross reference table.
func ParseFileToContext(file *os.File, config *Configuration) (*Context, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, fileInfo.Size())

	// read file content to buffer
	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	ctx, err := NewContext(bytes.NewReader(buffer), fileInfo.Name(), fileInfo.Size(), config)
	if err != nil {
		return nil, err
	}

	if ctx.Reader15 {
		log.Info.Println("PDF Version 1.5 conforming reader")
	} else {
		log.Info.Println("PDF Version 1.4 conforming reader - no object streams or xrefstreams allowed")
	}

	// Populate xRefTable.
	err = readXRefTable(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "xRefTable failed")
	}

	// Make all objects explicitly available (load into memory) in corresponding xRefTable entries.
	// Also decode any involved object streams.
	err = dereferenceXRefTable(ctx, config)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}
