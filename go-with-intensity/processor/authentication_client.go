package processor

import "net/http"

type AuthenticationClient struct {
	scheme   string
	hostname string
	inner    HTTPClient
}

func NewAuthenticationClient(inner HTTPClient, scheme string, hostname string, authId string, authToken string) *AuthenticationClient {
	return &AuthenticationClient{
		scheme:   scheme,
		hostname: hostname,
		inner:    inner,
	}
}

func (ac *AuthenticationClient) Do(request *http.Request) (*http.Response, error) {
	request.Host = ac.hostname
	request.URL.Scheme = ac.scheme
	request.URL.Host = ac.hostname
	return ac.inner.Do(request)

}
