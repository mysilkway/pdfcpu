package pdfcpu

import (
	"github.com/mysilkway/pdfcpu/pkg/log"
	"github.com/pkg/errors"
	"os"
)

// ParseFileToPDFContext parses the File and generates a PDFContext, an in-memory representation containing a cross reference table.
func ParseFileToPDFContext(file *os.File, config *Configuration) (*PDFContext, error) {
	ctx, err := NewPDFContext("", file, config)
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
