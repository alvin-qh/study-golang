package excel

import (
	"unsafe"

	"github.com/xuri/excelize/v2"
)

type Options excelize.Options

func optionsConv(opts []Options) []excelize.Options {
	ptr := unsafe.SliceData(opts)
	return unsafe.Slice((*excelize.Options)(ptr), len(opts))
}

type Excel struct {
	file *excelize.File
}

func New(filename string, opts ...Options) *Excel {
	f := excelize.NewFile(optionsConv(opts)...)
	f.SaveAs(filename)

	return &Excel{file: f}
}

func Open(filename string, opts ...Options) (*Excel, error) {
	f, err := excelize.OpenFile(filename, optionsConv(opts)...)
	if err != nil {
		return nil, err
	}
	return &Excel{file: f}, err
}

func (e *Excel) Save(filename string, opts ...Options) {
	e.file.Save(optionsConv(opts)...)
}
