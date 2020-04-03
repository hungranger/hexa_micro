package restclient

import (
	"hexa_micro/pkg/shortenservice/config"
	"hexa_micro/pkg/shortenservice/container/logger"
	"hexa_micro/pkg/shortenservice/interface/restClient/serializer"
	"hexa_micro/pkg/shortenservice/interface/restClient/serializer/json"
	"hexa_micro/pkg/shortenservice/interface/restClient/serializer/msgpack"
	"hexa_micro/pkg/shortenservice/usecase"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type RedirectHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	shortenURLUseCase usecase.IShortenUseCase
}

func NewHandler(shortenURLUseCase usecase.IShortenUseCase) RedirectHandler {
	return &handler{shortenURLUseCase}
}

func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		logger.Log.Fatalf("%+v", err)
		return
	}
}

func (h *handler) serializer(contentType string) serializer.IRedirectSerializer {
	if contentType == "application/x-msgpack" {
		return &msgpack.Redirect{}
	}
	return &json.Redirect{}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	// code := chi.URLParam(r, "code")
	code := strings.TrimPrefix(r.URL.Path, "/")
	redirect, err := h.shortenURLUseCase.Find(code)
	if err != nil {
		if errors.Cause(err) == config.ErrRedirectNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect.URL, http.StatusMovedPermanently)
}

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	redirect, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = h.shortenURLUseCase.Store(redirect)
	if err != nil {
		if errors.Cause(err) == config.ErrRedirectInvalid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	responseBody, err := h.serializer(contentType).Encode(redirect)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setupResponse(w, contentType, responseBody, http.StatusCreated)
}
