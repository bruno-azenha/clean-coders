package processor

type VerifyHandler struct {
	input    chan *Envelope
	output   chan *Envelope
	verifier Verifier
}

type Verifier interface {
	Verify(AddressInput) AddressOutput
}

func NewVerifyHandler(input, output chan *Envelope, verifier Verifier) *VerifyHandler {
	return &VerifyHandler{
		input:    input,
		output:   output,
		verifier: verifier,
	}
}

func (vh *VerifyHandler) Handle() {
	for envelope := range vh.input {
		envelope.Output = vh.verifier.Verify(envelope.Input)
		vh.output <- envelope
	}
}
