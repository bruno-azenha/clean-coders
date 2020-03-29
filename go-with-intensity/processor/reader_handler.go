package processor

import (
	"encoding/csv"
	"io"
)

type ReaderHandler struct {
	reader *csv.Reader
	closer io.Closer
	output *Envelope
}

func NewReaderHandler(reader io.ReadCloser, output chan *Envelope) *ReaderHandler {
	return &ReaderHandler{
		reader: csv.Reader(reader),
		closer: reader,
		output: output,
	}
}

func (rh *ReaderHandler) Handle() {

}
