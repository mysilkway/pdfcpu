package pdfcpu

import (
	"bufio"
	"bytes"
	"github.com/charleswklau/pdfcpu/pkg/log"
)

func WritePDFBuf(ctx *PDFContext) (*bytes.Buffer, error) {

	b := bytes.NewBuffer([]byte{})

	ctx.Write.Writer = bufio.NewWriter(b)

	err := handleEncryption(ctx)
	if err != nil {
		return nil, err
	}

	// Since we support PDF Collections (since V1.7) for file attachments
	// we need to always generate V1.7 PDF files.
	err = writeHeader(ctx.Write, V17)
	if err != nil {
		return nil, err
	}

	log.Debug.Printf("offset after writeHeader: %d\n", ctx.Write.Offset)

	// Write root object(aka the document catalog) and page tree.
	err = writeRootObject(ctx)
	if err != nil {
		return nil, err
	}

	log.Debug.Printf("offset after writeRootObject: %d\n", ctx.Write.Offset)

	// Write document information dictionary.
	err = writeDocumentInfoDict(ctx)
	if err != nil {
		return nil, err
	}

	log.Debug.Printf("offset after writeInfoObject: %d\n", ctx.Write.Offset)

	// Write offspec additional streams as declared in pdf trailer.
	if ctx.AdditionalStreams != nil {
		_, _, err = writeDeepObject(ctx, ctx.AdditionalStreams)
		if err != nil {
			return nil, err
		}
	}

	err = writeEncryptDict(ctx)
	if err != nil {
		return nil, err
	}

	// Mark redundant objects as free.
	// eg. duplicate resources, compressed objects, linearization dicts..
	deleteRedundantObjects(ctx)

	err = writeXRef(ctx)
	if err != nil {
		return nil, err
	}

	// Write pdf trailer.
	_, err = writeTrailer(ctx.Write)
	if err != nil {
		return nil, err
	}

	err = setFileSizeOfWrittenFileBuf(ctx.Write)
	if err != nil {
		return nil, err
	}

	if ctx.Read != nil {
		ctx.Write.BinaryImageSize = ctx.Read.BinaryImageSize
		ctx.Write.BinaryFontSize = ctx.Read.BinaryFontSize
		logWriteStats(ctx)
	}

	return b, nil
}

func setFileSizeOfWrittenFileBuf(w *WriteContext) error {

	// Get file info for file just written but flush first to get correct file size.

	err := w.Flush()
	if err != nil {
		return err
	}

	w.FileSize = int64(w.Size())

	return nil
}
