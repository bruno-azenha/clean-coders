package processor

import "io"

type Pipeline struct {
	reader  io.ReadCloser
	writer  io.WriteCloser
	workers int

	verifier      Verifier
	verifyInput   chan *Envelope
	sequenceInput chan *Envelope
	writerInput   chan *Envelope
}

func NewPipeline(reader io.ReadCloser, writer io.WriteCloser, client HTTPClient, workers int) *Pipeline {
	return &Pipeline{
		reader:  reader,
		writer:  writer,
		workers: workers,

		verifier:      NewSmartyVerifier(client),
		verifyInput:   make(chan *Envelope, 1024),
		sequenceInput: make(chan *Envelope, 1024),
		writerInput:   make(chan *Envelope, 1024),
	}
}

func (p *Pipeline) Process() (err error) {
	go func() {
		err = NewReaderHandler(p.reader, p.verifyInput).Handle()
	}()

	p.startVerifyHandlers()
	p.startSequenceHandler()
	p.awaitWriterHandler()

	return err
}

func (p *Pipeline) startSequenceHandler() {
	go NewSequenceHandler(p.sequenceInput, p.writerInput).Handle()
}

func (p *Pipeline) awaitWriterHandler() {
	NewWriterHandler(p.writerInput, p.writer).Handle()
}

func (p *Pipeline) startVerifyHandlers() {
	for i := 0; i < p.workers; i++ {
		go NewVerifyHandler(p.verifyInput, p.sequenceInput, p.verifier).Handle()
	}
}
