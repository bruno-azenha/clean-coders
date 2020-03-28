package processor

type SequenceHandler struct {
	input   chan *Envelope
	output  chan *Envelope
	counter int
	buffer  map[int]*Envelope
}

func NewSequenceHandler(input, output chan *Envelope) *SequenceHandler {
	return &SequenceHandler{
		input:   input,
		output:  output,
		counter: 0,
		buffer:  map[int]*Envelope{},
	}
}

func (sh *SequenceHandler) Handle() {
	for envelope := range sh.input {
		sh.processEnvelope(envelope)
	}

}

func (sh *SequenceHandler) processEnvelope(envelope *Envelope) {
	sh.buffer[envelope.Sequence] = envelope
	sh.sendBufferedEnvelopesInOrder()
}

func (sh *SequenceHandler) sendBufferedEnvelopesInOrder() {
	for {
		next, found := sh.buffer[sh.counter]
		if !found {
			break
		}
		sh.output <- next
		delete(sh.buffer, sh.counter)
		sh.counter++
	}
}
