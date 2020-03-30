package processor

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/smartystreets/assertions/should"

	"github.com/smartystreets/gunit"
)

type PipelineFixture struct {
	*gunit.Fixture

	reader   *ReadWriteSpyBuffer
	writer   *ReadWriteSpyBuffer
	client   *IntegrationHTTPClient
	pipeline *Pipeline
}

func TestPipelineFixture(t *testing.T) {
	gunit.Run(new(PipelineFixture), t)
}

func (pf *PipelineFixture) Setup() {
	pf.reader = NewReadWriteSpyBuffer("")
	pf.writer = NewReadWriteSpyBuffer("")
	pf.client = &IntegrationHTTPClient{}
	pf.pipeline = NewPipeline(pf.reader, pf.writer, pf.client, 2)
}

func (pf *PipelineFixture) LongTestPipeline() {
	fmt.Fprintln(pf.reader, "Street1,City,State,ZIPCode")
	fmt.Fprintln(pf.reader, "A,B,C,D")
	fmt.Fprintln(pf.reader, "A,B,C,D")

	err := pf.pipeline.Process()

	expectedOutput := "Status,DeliveryLine1,LastLine,City,State,ZIPCode\n" +
		"Deliverable,AA,BB,CC,DD,EE\n" +
		"Deliverable,AA,BB,CC,DD,EE\n"

	pf.So(err, should.BeNil)
	pf.So(pf.writer.String(), should.Equal, expectedOutput)
}

type IntegrationHTTPClient struct{}

func (this *IntegrationHTTPClient) Do(request *http.Request) (*http.Response, error) {
	return &http.Response{
		Body:       NewReadWriteSpyBuffer(integrationJSONOutput),
		StatusCode: http.StatusOK,
	}, nil
}

const integrationJSONOutput = `
[
	{
		"delivery_line_1": "AA",
		"last_line": "BB",
		"components": {
			"city_name": "CC",
			"state_abbreviation": "DD",
			"zipcode": "EE"
		},
		"analysis": {
			"dpv_match_code": "Y",
			"dpv_vacant": "N",
			"active": "Y"
		}
	}
]`
