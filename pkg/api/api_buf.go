package api

import (
	"bytes"
	"fmt"
	"github.com/charleswklau/pdfcpu/pkg/log"
	"github.com/charleswklau/pdfcpu/pkg/pdfcpu"
	"github.com/oxtoacart/bpool"
	"github.com/pkg/errors"
	"time"
)

var (
	pdfBufp = bpool.NewBufferPool(512)
)

// Write generates a PDF file for a given PDFContext.
func WriteBuf(ctx *pdfcpu.PDFContext, ob *bytes.Buffer) (*bytes.Buffer, error) {

	fmt.Printf("writing %s ...\n", ctx.Write.DirName+ctx.Write.FileName)
	//logInfoAPI.Printf("writing to %s..\n", fileName)

	b, err := pdfcpu.WritePDFBuf(ctx, ob)
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

func MergeAsBuf(filesIn []string, config *pdfcpu.Configuration, ob *bytes.Buffer) (*bytes.Buffer, error) {

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

	ctxDest.Write.Command = "Merge"

	ctxDest.Write.DirName = ""
	ctxDest.Write.FileName = ""

	_, err = WriteBuf(ctxDest, ob)
	if err != nil {
		return nil, err
	}

	log.Stats.Printf("XRefTable:\n%s\n", ctxDest)

	return ob, nil
}
