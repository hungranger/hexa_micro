package restclient_test

import (
	"bytes"
	"encoding/json"
	restclient "hexa_micro/pkg/shortenservice/interface/restClient"
	"hexa_micro/pkg/shortenservice/mocks"
	"hexa_micro/pkg/shortenservice/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_Get(t *testing.T) {
	mockRedirect := &model.Redirect{
		Code:     "XH2ey9WR",
		URL:      "http://vnexpress.net",
		CreateAt: 123456,
	}

	mockShortenUCase := new(mocks.ShortenUseCaseFake)
	h := restclient.NewHandler(mockShortenUCase)

	type args struct {
		reqURL string
	}
	type wants struct {
		status int
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		// TODO: Add test cases.
		{
			"Invalid code, return 500 http status code",
			args{"/"},
			wants{
				status: http.StatusInternalServerError,
			},
		},
		{
			"Not Found url, return 404 http status code",
			args{"/abc"},
			wants{
				status: http.StatusNotFound,
			},
		},
		{
			"Found url, return 301 http status code",
			args{"/" + mockRedirect.Code},
			wants{
				status: http.StatusMovedPermanently,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := buildRequest(http.MethodGet, tt.args.reqURL, nil)
			res := httptest.NewRecorder()
			h.Get(res, req)
			assertStatus(t, res.Code, tt.wants.status)
		})
	}
}

func Test_handler_Post(t *testing.T) {
	mockRedirect := &model.Redirect{
		Code:     "XH2ey9WR",
		URL:      "http://vnexpress.net",
		CreateAt: 123456,
	}

	mockShortenUCase := new(mocks.ShortenUseCaseFake)
	h := restclient.NewHandler(mockShortenUCase)

	redirect := model.Redirect{}

	type args struct {
		redirectURL string
	}
	type wants struct {
		status int
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		// TODO: Add test cases.
		{
			"Shorten URL successfully, return status 201",
			args{"http://vnexpress.net"},
			wants{http.StatusCreated},
		},
		{
			"Invalid URL, return status 400",
			args{"http/vnexpress.net"},
			wants{http.StatusBadRequest},
		},
		{
			"Empty URL, return status 500",
			args{""},
			wants{http.StatusInternalServerError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redirect.URL = tt.args.redirectURL
			body, _ := json.Marshal(&redirect)
			req, _ := buildRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			// req.Header.Add("content-type", "application/x-msgpack")

			res := httptest.NewRecorder()
			h.Post(res, req)
			assertStatus(t, res.Code, tt.wants.status)

			wantBody, _ := json.Marshal(mockRedirect)
			assertResponseBody(t, string(res.Body.String()), string(wantBody))
		})
	}
}

func buildRequest(method string, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if json.Valid([]byte(got)) && got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
