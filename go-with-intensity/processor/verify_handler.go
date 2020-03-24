package processor

type VerifyHandler struct {
	input       chan *Envelope
	output      chan *Envelope
	application Verifier
}

type Verifier interface {
	Verify(AddressInput)
}

func NewVerifyHandler(input, output chan *Envelope, application Verifier) *VerifyHandler {
	return &VerifyHandler{
		input:       input,
		output:      output,
		application: application,
	}
}

func (vh *VerifyHandler) Handle() {
	received := <-vh.input

	vh.application.Verify(received.Input)

	vh.output <- received
}
