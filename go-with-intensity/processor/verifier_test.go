package processor

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type VerifierFixture struct {
	*gunit.Fixture

	client   *FakeHTTPClient
	verifier *SmartyVerifier
}

func TestVerifierFixture(t *testing.T) {
	gunit.Run(new(VerifierFixture), t)
}

func (vf *VerifierFixture) Setup() {
	vf.client = &FakeHTTPClient{}
	vf.verifier = NewSmartyVerifier(vf.client)
}

func NewSmartyVerifier(client HTTPClient) *SmartyVerifier {
	return &SmartyVerifier{
		client: client,
	}
}

func (vf *VerifierFixture) TestRequestComposedProperly() {
	input := AddressInput{
		Street1: "Street1",
		City:    "City",
		State:   "State",
		ZIPCode: "ZIPCode",
	}

	vf.client.Configure("[{}]", http.StatusOK, nil)

	vf.verifier.Verify(input)

	vf.So(vf.client.request.Method, should.Equal, "GET")
	vf.So(vf.client.request.URL.Path, should.Equal, "/street-address")
	vf.AssertQueryStringValue("street", input.Street1)
	vf.AssertQueryStringValue("city", input.City)
	vf.AssertQueryStringValue("state", input.State)
	vf.AssertQueryStringValue("zipcode", input.ZIPCode)
}

func (vf *VerifierFixture) TestResponseParsed() {
	vf.client.Configure(rawJSONOutput, http.StatusOK, nil)
	result := vf.verifier.Verify(AddressInput{})

	vf.So(result.DeliveryLine1, should.Equal, "1 Santa Claus ln")
	vf.So(result.LastLine, should.Equal, "North Pole AK 99705-9901")
	vf.So(result.City, should.Equal, "North Pole")
	vf.So(result.State, should.Equal, "AK")
	vf.So(result.ZIPCode, should.Equal, "99705")

}

func (vf *VerifierFixture) TestMalformedJSONHandled() {
	const malformedRawJSONOutput = `I am not JSON!`
	vf.client.Configure(malformedRawJSONOutput, http.StatusOK, nil)
	result := vf.verifier.Verify(AddressInput{})
	vf.So(result.Status, should.Equal, "Invalid API response")
}

func (vf *VerifierFixture) TestHTTPErrorHandled() {
	vf.client.Configure("", 0, errors.New("gophers"))
	result := vf.verifier.Verify(AddressInput{})
	vf.So(result.Status, should.Equal, "Invalid API response")
}

func (vf *VerifierFixture) TestHTTPResponseBodyClosed() {
	vf.client.Configure(rawJSONOutput, http.StatusOK, nil)
	vf.verifier.Verify(AddressInput{})
	vf.So(vf.client.responseBody.timesClosed, should.Equal, 1)
}

func (vf *VerifierFixture) TestAddressStatus() {
	var (
		deliverableJSON      = buildAnalysisJSON("Y", "N", "Y")
		missingSecondaryJSON = buildAnalysisJSON("D", "N", "Y")
		droppedSecondaryJSON = buildAnalysisJSON("S", "N", "Y")
		vacantJSON           = buildAnalysisJSON("Y", "Y", "Y")
		inactiveJSON         = buildAnalysisJSON("Y", "N", "?")
		invalidJSON          = buildAnalysisJSON("N", "?", "?")
	)

	vf.verifyAndAssertStatus(deliverableJSON, "Deliverable")
	vf.verifyAndAssertStatus(missingSecondaryJSON, "Deliverable")
	vf.verifyAndAssertStatus(droppedSecondaryJSON, "Deliverable")
	vf.verifyAndAssertStatus(inactiveJSON, "Inactive")
	vf.verifyAndAssertStatus(vacantJSON, "Vacant")
	vf.verifyAndAssertStatus(invalidJSON, "Invalid")
}

func (vf *VerifierFixture) verifyAndAssertStatus(jsonResponse, expectedStatus string) {
	vf.client.Configure(jsonResponse, http.StatusOK, nil)
	output := vf.verifier.Verify(AddressInput{})
	vf.So(output.Status, should.Equal, expectedStatus)
}

func buildAnalysisJSON(match, vacant, active string) string {
	template := `[
		{
			"analysis": {
				"dpv_match_code": "%s",
				"dpv_vacant": "%s",
				"active": "%s"
			}
		}
	]`
	return fmt.Sprintf(template, match, vacant, active)
}

const rawJSONOutput = `
[
	{
		"delivery_line_1": "1 Santa Claus ln",
		"last_line": "North Pole AK 99705-9901",
		"components": {
			"city_name": "North Pole",
			"state_abbreviation": "AK",
			"zipcode": "99705"
		}
	}
]`

func (vf VerifierFixture) rawQuery() string {
	return vf.client.request.URL.RawQuery
}

func (vf VerifierFixture) AssertQueryStringValue(key, expected string) {
	query := vf.client.request.URL.Query()
	vf.So(query.Get(key), should.Equal, expected)
}

////////////////////////////

type FakeHTTPClient struct {
	request      *http.Request
	response     *http.Response
	responseBody *VerifierSpyBuffer
	err          error
}

func (fc *FakeHTTPClient) Configure(responseText string, statusCode int, err error) {
	if err == nil {
		fc.responseBody = NewVerifierSpyBuffer(responseText)
		fc.response = &http.Response{
			Body:       fc.responseBody,
			StatusCode: statusCode,
		}
	}
	fc.err = err
}

func (fc *FakeHTTPClient) Do(request *http.Request) (*http.Response, error) {
	fc.request = request
	return fc.response, fc.err
}

///////////////////////////

type VerifierSpyBuffer struct {
	*bytes.Buffer
	timesClosed int
}

func NewVerifierSpyBuffer(value string) *VerifierSpyBuffer {
	return &VerifierSpyBuffer{
		Buffer: bytes.NewBufferString(value),
	}
}

func (sb *VerifierSpyBuffer) Close() error {
	sb.timesClosed++
	sb.Buffer.Reset()
	return nil
}
