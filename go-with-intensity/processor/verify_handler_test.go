package processor

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestHandlerFixture(t *testing.T) {
	gunit.Run(new(HandlerFixture), t)
}

type HandlerFixture struct {
	*gunit.Fixture

	input       chan *Envelope
	output      chan *Envelope
	application *FakeVerifier
	handler     *VerifyHandler
}

func (hf *HandlerFixture) Setup() {
	hf.input = make(chan *Envelope, 10)
	hf.output = make(chan *Envelope, 10)
	hf.application = NewFakeVerifier()
	hf.handler = NewVerifyHandler(hf.input, hf.output, hf.application)
}

func (hf *HandlerFixture) TestVerifyHandlerReceivesInput() {
	envelope := &Envelope{
		Input: AddressInput{
			Street1: "42",
		},
	}
	hf.input <- envelope
	close(hf.input)

	hf.handler.Handle()

	hf.AssertEqual(envelope, <-hf.output)
	hf.AssertEqual(envelope.Input, hf.application.input)
}

///////////////////////////////////

type FakeVerifier struct {
	input AddressInput
}

func NewFakeVerifier() *FakeVerifier {
	return &FakeVerifier{}
}

func (fv *FakeVerifier) Verify(value AddressInput) {
	fv.input = value
}
