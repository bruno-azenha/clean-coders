package processor

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"

	"github.com/smartystreets/gunit"
)

func TestAuthenticationClient(t *testing.T) {
	gunit.Run(new(AuthenticationClientFixture), t)
}

type AuthenticationClientFixture struct {
	*gunit.Fixture

	inner  *FakeHTTPClient
	client *AuthenticationClient
}

func (acf *AuthenticationClientFixture) Setup() {
	acf.inner = &FakeHTTPClient{}
	acf.client = NewAuthenticationClient(acf.inner, "https", "HOSTNAME", "authid", "authtoken")
}

func (acf *AuthenticationClientFixture) TestProvidedInformationAddedBeforeRequestSent() {
	request := httptest.NewRequest("GET", "/path", nil)

	acf.client.Do(request)

	acf.So(acf.inner.request.Host, should.Equal, "HOSTNAME")
	acf.So(acf.inner.request.URL.Host, should.Equal, "HOSTNAME")
	acf.So(acf.inner.request.URL.Scheme, should.Equal, "https")
	acf.So(acf.inner.request.URL.Query().Get("auth-id"), should.Equal, "authid")
	acf.So(acf.inner.request.URL.Query().Get("auth-token"), should.Equal, "authtoken")
}

func (acf *AuthenticationClientFixture) TestResponseAndErrorFromInnerClientReturned() {
	acf.inner.response = &http.Response{
		StatusCode: http.StatusTeapot,
	}
	acf.inner.err = errors.New("HTTP Error")

	request := httptest.NewRequest("GET", "/path", nil)
	response, err := acf.client.Do(request)

	acf.So(response.StatusCode, should.Equal, http.StatusTeapot)
	acf.So(err.Error(), should.Equal, "HTTP Error")
}
