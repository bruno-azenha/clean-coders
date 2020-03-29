package processor

import (
	"encoding/csv"
	"io"
)

type WriterHandler struct {
	input  chan *Envelope
	writer *csv.Writer
	closer io.Closer
}

func NewWriterHandler(input chan *Envelope, output io.WriteCloser) *WriterHandler {
	whf := &WriterHandler{
		input:  input,
		writer: csv.NewWriter(output),
		closer: output,
	}
	whf.writer.Write([]string{"Status", "DeliveryLine1", "LastLine", "City", "State", "ZIPCode"})
	return whf
}

func (wh *WriterHandler) Handle() {

	for envelope := range wh.input {
		wh.writeAddressOutput(envelope.Output)
	}

	wh.writer.Flush()
	wh.closer.Close()
}

func (wh *WriterHandler) writeAddressOutput(output AddressOutput) {
	wh.writer.Write([]string{
		output.Status,
		output.DeliveryLine1,
		output.LastLine,
		output.City,
		output.State,
		output.ZIPCode,
	})
}

func (wh *WriterHandler) writeValues(values ...string) {
	wh.writer.Write(values)
}
