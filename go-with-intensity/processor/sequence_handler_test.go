package processor

import (
	"testing"

	"github.com/smartystreets/assertions/should"

	"github.com/smartystreets/gunit"
)

type SequenceHandlerFixture struct {
	*gunit.Fixture

	input   chan *Envelope
	output  chan *Envelope
	handler *SequenceHandler
}

func TestSequenceHandlerFixture(t *testing.T) {
	gunit.Run(new(SequenceHandlerFixture), t)
}

func (shf *SequenceHandlerFixture) Setup() {
	shf.input = make(chan *Envelope, 10)
	shf.output = make(chan *Envelope, 10)
	shf.handler = NewSequenceHandler(shf.input, shf.output)
}

func (shf *SequenceHandlerFixture) TestExpectedEnvelopeSentToOutput() {
	shf.sendEnvelopesInSequence(0, 1, 2, 3)

	shf.handler.Handle()

	shf.So(shf.sequenceOrder(), should.Resemble, []int{0, 1, 2})
	shf.So(shf.handler.buffer, should.BeEmpty)
}

func (shf *SequenceHandlerFixture) TestEnvelopesReceivedOutOfOrder_BufferedUntilContiguousBlock() {
	shf.sendEnvelopesInSequence(4, 3, 2, 1, 0)

	shf.handler.Handle()

	shf.So(shf.sequenceOrder(), should.Resemble, []int{0, 1, 2, 3})
	shf.So(shf.handler.buffer, should.BeEmpty)
}

func (shf *SequenceHandlerFixture) sendEnvelopesInSequence(sequences ...int) {
	max := maxInt(sequences)
	for sequence := range sequences {
		shf.input <- &Envelope{Sequence: sequence, EOF: max == sequence}
	}
}

func maxInt(ints []int) (max int) {
	max = 0
	for _, value := range ints {
		if value > max {
			max = value
		}
	}
	return max
}

func (shf *SequenceHandlerFixture) sequenceOrder() (order []int) {
	for envelope := range shf.output {
		order = append(order, envelope.Sequence)
	}
	return order
}
