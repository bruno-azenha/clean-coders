package processor

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestWriterHandlerFixture(t *testing.T) {
	gunit.Run(new(WriterHandlerFixture), t)
}

type WriterHandlerFixture struct {
	*gunit.Fixture

	handler *WriterHandler
	input   chan *Envelope
	buffer  *ReadWriteSpyBuffer
	writer  *csv.Writer
}

func (whf *WriterHandlerFixture) Setup() {
	whf.buffer = NewReadWriteSpyBuffer("")
	whf.input = make(chan *Envelope, 10)
	whf.handler = NewWriterHandler(whf.input, whf.buffer)
}

func (whf *WriterHandlerFixture) TestHeaderWritter() {
	close(whf.input)
	whf.handler.Handle()

}

func (whf *WriterHandlerFixture) TestOutputClosed() {
	close(whf.input)
	whf.handler.Handle()

	whf.So(whf.buffer.timesClosed, should.Equal, 1)
}

var recordMatchingHeader = AddressOutput{
	Status:        "Status",
	DeliveryLine1: "DeliveryLine1",
	LastLine:      "LastLine",
	City:          "City",
	State:         "State",
	ZIPCode:       "ZIPCode",
}

func (whf *WriterHandlerFixture) TestHeaderMatchesRecords() {
	whf.input <- &Envelope{Output: recordMatchingHeader}
	close(whf.input)

	whf.handler.Handle()
	whf.assertHeaderMatchesRecord()
}

func (whf *WriterHandlerFixture) assertHeaderMatchesRecord() {
	lines := whf.outputLines()
	header := lines[0]
	record := lines[1]

	whf.So(header, should.Equal, "Status,DeliveryLine1,LastLine,City,State,ZIPCode")
	whf.So(record, should.Equal, header)
}

func (whf *WriterHandlerFixture) TestAllEnvelopesWritten() {
	whf.sendEnvelopes(2)

	whf.handler.Handle()

	if lines := whf.outputLines(); whf.So(lines, should.HaveLength, 3) {
		whf.So(lines[1], should.Equal, "A1,B1,C1,D1,E1,F1")
		whf.So(lines[2], should.Equal, "A2,B2,C2,D2,E2,F2")
	}
}

func (whf *WriterHandlerFixture) sendEnvelopes(count int) {
	for i := 1; i < count+1; i++ {
		whf.input <- &Envelope{Output: createOutput(strconv.Itoa(i))}
	}
	close(whf.input)
}

func createOutput(index string) AddressOutput {
	return AddressOutput{
		Status:        "A" + index,
		DeliveryLine1: "B" + index,
		LastLine:      "C" + index,
		City:          "D" + index,
		State:         "E" + index,
		ZIPCode:       "F" + index,
	}
}

func (whf *WriterHandlerFixture) outputLines() []string {
	outputFile := strings.TrimSpace(whf.buffer.String())
	return strings.Split(outputFile, "\n")
}

///////////////////////////

type ReadWriteSpyBuffer struct {
	*bytes.Buffer
	timesClosed int
}

func NewReadWriteSpyBuffer(value string) *ReadWriteSpyBuffer {
	return &ReadWriteSpyBuffer{
		Buffer: bytes.NewBufferString(value),
	}
}

func (sb *ReadWriteSpyBuffer) Close() error {
	sb.timesClosed++
	return nil
}
