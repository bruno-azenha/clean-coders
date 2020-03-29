package processor

import (
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
}

func (rhf *ReaderHandlerFixture) Setup() {

}

func (rhf *ReaderHandlerFixture) TestCsvRecordSentInEnvelope() {
	buffer := NewReadWriteSpyBuffer("Street1,City,State,ZIPCode")
	output := make(chan *Envelope, 10)
	reader := NewReaderHandler(buffer, output)

	reader.Handle()

	rhf.So(<-output, should.Resemble)
}
