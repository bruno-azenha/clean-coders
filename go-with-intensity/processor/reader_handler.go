package processor

import (
	"encoding/csv"
	"errors"
	"io"
)

type ReaderHandler struct {
	reader   *csv.Reader
	closer   io.Closer
	output   chan *Envelope
	sequence int
}

func NewReaderHandler(reader io.ReadCloser, output chan *Envelope) *ReaderHandler {
	return &ReaderHandler{
		reader:   csv.NewReader(reader),
		closer:   reader,
		output:   output,
		sequence: initialSequenceValue,
	}
}

func (rh *ReaderHandler) Handle() error {
	defer rh.close()
	rh.skipHeader()

	for {
		record, err := rh.reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.New("Malformed input")
		}
		rh.sendEnvelope(record)
	}

	return nil
}

func (rh *ReaderHandler) skipHeader() {
	rh.reader.Read()
}

func (rh *ReaderHandler) sendEnvelope(record []string) {
	rh.output <- &Envelope{
		Sequence: rh.sequence,
		Input:    createInput(record),
	}
	rh.sequence++
}

func createInput(record []string) AddressInput {
	return AddressInput{
		Street1: record[0],
		City:    record[1],
		State:   record[2],
		ZIPCode: record[3],
	}
}

func (rh *ReaderHandler) close() {
	rh.output <- &Envelope{Sequence: rh.sequence, EOF: true}
	close(rh.output)
	rh.closer.Close()
}
