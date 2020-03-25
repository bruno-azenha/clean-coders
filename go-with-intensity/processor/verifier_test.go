package processor

import (
	"bytes"
	"io/ioutil"
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

	vf.verifier.Verify(input)

	vf.So(vf.client.request.Method, should.Equal, "GET")
	vf.So(vf.client.request.URL.Path, should.Equal, "/street-address")
	vf.AssertQueryStringValue("street", "Street1")
	vf.AssertQueryStringValue("city", "City")
	vf.AssertQueryStringValue("state", "State")
	vf.AssertQueryStringValue("zipcode", "ZIPCode")
}

func (vf *VerifierFixture) TestResponseParsed() {
	vf.client.response = &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString(`[{""}]`)),
		StatusCode: http.StatusOK,
	}

	result := vf.verifier.Verify(AddressInput{})
	vf.So(result.DeliveryLine1, should.Equal, "Hello World!")
}

func (vf VerifierFixture) rawQuery() string {
	return vf.client.request.URL.RawQuery
}

func (vf VerifierFixture) AssertQueryStringValue(key, expected string) {
	query := vf.client.request.URL.Query()
	vf.So(query.Get(key), should.Equal, expected)
}

////////////////////////////

type FakeHTTPClient struct {
	request  *http.Request
	response *http.Response
	err      error
}

func (fc *FakeHTTPClient) Configure(responseText string, statusCode int, err error) {
	fc.response = &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString(responseText)),
		StatusCode: statusCode,
	}
	fc.err = err
}

func (fc *FakeHTTPClient) Do(request *http.Request) (*http.Response, error) {
	fc.request = request
	return fc.response, fc.err
}
