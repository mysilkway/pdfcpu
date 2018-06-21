package api

import (
	"bytes"
	"github.com/mysilkway/pdfcpu/pkg/log"
	"github.com/mysilkway/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
	"os"
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
func MergeToBuf(filesIn []string, config *pdfcpu.Configuration) (*bytes.Buffer, error) {

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

// Merge some PDF files together and write the result to the buffer.
// This corresponds to concatenating these files in the order specified by file array.
// The first entry of files serves as the destination xRefTable where all the remaining files gets merged into.
func MergeFileToBuf(files []*os.File, config *pdfcpu.Configuration) (*bytes.Buffer, error) {

	ctxDest, _, _, err := readFileAndValidate(files[0], config, time.Now())
	if err != nil {
		return nil, err
	}

	if ctxDest.XRefTable.Version() < pdfcpu.V15 {
		v, _ := pdfcpu.Version("1.5")
		ctxDest.XRefTable.RootVersion = &v
		log.Stats.Println("Ensure V1.5 for writing object & xref streams")
	}

	// Repeatedly merge files into fileDest's xref table.
	for _, f := range files[1:] {
		err = appendFileTo(f, ctxDest)
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

func readFileAndValidate(f *os.File, config *pdfcpu.Configuration, from1 time.Time) (ctx *pdfcpu.PDFContext, dur1, dur2 float64, err error) {

	ctx, err = ReadFile(f, config)
	if err != nil {
		return nil, 0, 0, err
	}
	dur1 = time.Since(from1).Seconds()

	from2 := time.Now()
	//fmt.Printf("validating %s ...\n", fileIn)
	//logInfoAPI.Printf("validating %s..\n", fileIn)
	err = pdfcpu.ValidateXRefTable(ctx.XRefTable)
	if err != nil {
		return nil, 0, 0, err
	}
	dur2 = time.Since(from2).Seconds()

	return ctx, dur1, dur2, nil
}

// ReadFile reads in a PDF file and builds an internal structure holding its cross reference table aka the PDFContext.
func ReadFile(f *os.File, config *pdfcpu.Configuration) (*pdfcpu.PDFContext, error) {

	//logInfoAPI.Printf("reading %s..\n", fileIn)
	ctx, err := pdfcpu.ParseFileToPDFContext(f, config)
	if err != nil {
		return nil, errors.Wrap(err, "Read failed.")
	}

	return ctx, nil
}

// appendFileTo appends file to ctxDest's page tree.
func appendFileTo(f *os.File, ctxDest *pdfcpu.PDFContext) error {

	log.Stats.Printf("appendTo: appending %s to %s\n", f, ctxDest.Read.FileName)

	// Build a PDFContext for fileIn.
	ctxSource, _, _, err := readFileAndValidate(f, ctxDest.Configuration, time.Now())
	if err != nil {
		return err
	}

	// Merge the source context into the dest context.
	return pdfcpu.MergeXRefTables(ctxSource, ctxDest)
}
