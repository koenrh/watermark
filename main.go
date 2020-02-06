package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/validate"
)

type Document struct {
	inFileName     string
	targetFileName string
	watermark      *pdfcpu.Watermark
}

func NewDocument(textLines []string, inFileName string, includeDate bool) (*Document, error) {
	doc := &Document{}
	doc.inFileName = inFileName
	doc.targetFileName = strings.TrimSuffix(inFileName, filepath.Ext(inFileName)) + "_watermarked.pdf"

	if includeDate {
		textLines = append(textLines, time.Now().Format("January 2, 2006"))
	}

	config := pdfcpu.DefaultWatermarkConfig()
	config.Opacity = 0.45
	config.TextLines = textLines
	config.OnTop = true

	doc.watermark = config

	return doc, nil
}

func (d *Document) EmbedWatermark() error {
	context, err := pdfcpu.ReadFile(d.inFileName, pdfcpu.NewDefaultConfiguration())
	if err != nil {
		return err
	}

	if err := validate.XRefTable(context.XRefTable); err != nil {
		return err
	}

	if err := pdfcpu.AddWatermarks(context, nil, d.watermark); err != nil {
		return err
	}

	f, err := os.Create(d.targetFileName)
	if err != nil {
		return err
	}

	if err := api.WriteContext(context, f); err != nil {
		f.Close()
		return err
	}

	return f.Close()
}

type sliceFlags []string

func (i *sliceFlags) String() string {
	return strings.Join(*i, ", ")
}

func (i *sliceFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s [options] <path>:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var textLines sliceFlags

	flag.Var(&textLines, "text", "Watermark text lines.")
	includeDate := flag.Bool("date", false, "Whether to inlcude a date.")

	flag.Parse()

	if len(textLines) == 0 || flag.NArg() != 1 {
		Usage()
		os.Exit(1)
	}

	outFileName := flag.Arg(0)

	d, err := NewDocument(textLines, outFileName, *includeDate)

	if err != nil {
		fmt.Println("failed to create watermark")
		os.Exit(1)
	}

	err = d.EmbedWatermark()

	if err != nil {
		fmt.Println("failed to write watermark to PDF document")
		os.Exit(1)
	}
}
