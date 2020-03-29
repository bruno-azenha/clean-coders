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
		counter: initialSequenceValue,
		buffer:  map[int]*Envelope{},
	}
}

func (sh *SequenceHandler) Handle() {
	for envelope := range sh.input {
		sh.processEnvelope(envelope)
	}

	close(sh.output)
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
		sh.sendNextEnvelope(next)
	}
}

func (sh *SequenceHandler) sendNextEnvelope(envelope *Envelope) {
	if envelope.EOF {
		close(sh.input)
	} else {
		sh.output <- envelope

	}
	delete(sh.buffer, sh.counter)
	sh.counter++
}
