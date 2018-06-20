package api

import (
	"bytes"
	"github.com/charleswklau/pdfcpu/pkg/log"
	"github.com/charleswklau/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
	"time"
)

// WriteBuf generates a PDF file bytes for a given PDFContext.
func WriteBuf(ctx *pdfcpu.PDFContext) (*bytes.Buffer, error) {
	b, err := pdfcpu.WritePDFBuf(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Write failed.")
	}

	if ctx.StatsFileName != "" {
		err = pdfcpu.AppendStatsFile(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "Write stats failed.")
		}
	}

	return b, nil
}

// Merge some PDF files together and write the result to the buffer.
// This corresponds to concatenating these files in the order specified by filesIn.
// The first entry of filesIn serves as the destination xRefTable where all the remaining files gets merged into.
func MergeAsBuf(filesIn []string, config *pdfcpu.Configuration) (*bytes.Buffer, error) {

	ctxDest, _, _, err := readAndValidate(filesIn[0], config, time.Now())
	if err != nil {
		return nil, err
	}

	if ctxDest.XRefTable.Version() < pdfcpu.V15 {
		v, _ := pdfcpu.Version("1.5")
		ctxDest.XRefTable.RootVersion = &v
		log.Stats.Println("Ensure V1.5 for writing object & xref streams")
	}

	// Repeatedly merge files into fileDest's xref table.
	for _, f := range filesIn[1:] {
		err = appendTo(f, ctxDest)
		if err != nil {
			return nil, err
		}
	}

	err = pdfcpu.OptimizeXRefTable(ctxDest)
	if err != nil {
		return nil, err
	}

	err = pdfcpu.ValidateXRefTable(ctxDest.XRefTable)
	if err != nil {
		return nil, err
	}

	b, err := WriteBuf(ctxDest)
	if err != nil {
		return nil, err
	}

	log.Stats.Printf("XRefTable:\n%s\n", ctxDest)

	return b, nil
}
