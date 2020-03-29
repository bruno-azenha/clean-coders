package processor

import (
	"strconv"
	"testing"

	"github.com/smartystreets/assertions/should"

	"github.com/smartystreets/gunit"
)

func TestReaderHandlerFixture(t *testing.T) {
	gunit.Run(new(ReaderHandlerFixture), t)
}

type ReaderHandlerFixture struct {
	*gunit.Fixture

	handler *ReaderHandler
	buffer  *ReadWriteSpyBuffer
	output  chan *Envelope
}

func (rhf *ReaderHandlerFixture) Setup() {
	rhf.buffer = NewReadWriteSpyBuffer("")
	rhf.output = make(chan *Envelope, 10)
	rhf.handler = NewReaderHandler(rhf.buffer, rhf.output)

	const header = "Street1,City,State,ZIPCode"
	rhf.writeLine(header)
}

func (rhf *ReaderHandlerFixture) TestAllCSVRecordsSentToOutput() {
	rhf.writeLine("A1,B1,C1,D1")
	rhf.writeLine("A2,B2,C2,D2")

	rhf.handler.Handle()

	rhf.assertRecordsSent()
	rhf.assertCleanup()
}

func (rhf *ReaderHandlerFixture) writeLine(line string) {
	rhf.buffer.WriteString(line + "\n")
}

func (rhf *ReaderHandlerFixture) assertRecordsSent() {
	rhf.So(<-rhf.output, should.Resemble, buildEnvelope(initialSequenceValue))
	rhf.So(<-rhf.output, should.Resemble, buildEnvelope(initialSequenceValue+1))
}

func (rhf *ReaderHandlerFixture) assertCleanup() {
	rhf.So(<-rhf.output, should.Resemble, &Envelope{Sequence: initialSequenceValue + 2, EOF: true})
	rhf.So(<-rhf.output, should.BeNil)
	rhf.So(rhf.buffer.timesClosed, should.Equal, 1)
}

func buildEnvelope(index int) *Envelope {
	suffix := strconv.Itoa(index + 1)
	return &Envelope{
		Sequence: index,
		Input: AddressInput{
			Street1: "A" + suffix,
			City:    "B" + suffix,
			State:   "C" + suffix,
			ZIPCode: "D" + suffix,
		}}
}

func (rhf *ReaderHandlerFixture) TestMalformedInputReturnsError() {
	malformedRecord := "A1"
	rhf.writeLine(malformedRecord)

	err := rhf.handler.Handle()

	if rhf.So(err, should.NotBeNil) {
		rhf.So(err.Error(), should.Resemble, "Malformed input")
	}
}
