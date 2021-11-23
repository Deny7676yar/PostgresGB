package requester

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubGetLinks struct {
	urlcrw map[string]string
	title []string
}

func (s *StubGetLinks) Get(links string) string {
	urlcrw := s.urlcrw[links]
	return urlcrw
}

func (s *StubGetLinks) RecordWin(link string){
	s.title = append(s.title, link)
}

func TestRequester_Get(t *testing.T) {
	page := StubGetLinks{
		map[string]string{
			"www.yandex.ru" : "yandexmarket",
			"www.google.com" : "googlplay",
		},
		nil,
	}

	server := &UrlGetServer{&page}

	tests := []struct{
		links string
		urlreq string
		urlresp string
		expectedHTTPStatus int
	}{
		{
			links:      "yandex",
			urlreq:     "www.yandex.ru",
			urlresp:    "yandexmarket",
			expectedHTTPStatus: http.StatusOK,
		},
		{
			links:      "google",
			urlreq:     "www.google.com",
			urlresp:    "googlplay",
			expectedHTTPStatus: http.StatusOK,
		},
		{
			links: "Returns 404",
			urlreq:   "nil",
			urlresp:  "nil",
			expectedHTTPStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.links, func(t *testing.T) {
			request := newGetRequest(tt.links)
			response := httptest.NewRecorder()

			server.(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.urlresp)

		})
	}

}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func newGetRequest(title string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", title), nil)
	return req
}
