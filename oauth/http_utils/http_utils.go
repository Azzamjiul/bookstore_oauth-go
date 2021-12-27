package http_utils

import (
	"encoding/json"
	"net/http"

	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

type client struct{}

type Client interface {
	Get(string, map[string]string, interface{}) (int, error)
	Post(string, interface{}, map[string]string, interface{}) (int, error)
}

func New() Client {
	return &client{}
}

func (g client) Get(url string, headers map[string]string, dest interface{}) (int, error) {
	getRequest := gorequest.New()

	agent := getRequest.Get(url)
	for k, v := range headers {
		agent = agent.Set(k, v)
	}

	resp, bytes, errs := agent.EndBytes()
	if len(errs) > 0 {
		return http.StatusInternalServerError, errs[0]
	}

	if err := json.Unmarshal(bytes, &dest); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "error while reading response content")
	}

	return resp.StatusCode, nil
}

func (g client) Post(url string, data interface{}, headers map[string]string, dest interface{}) (int, error) {
	postRequest := gorequest.New()

	agent := postRequest.Post(url)
	for k, v := range headers {
		agent = agent.Set(k, v)
	}

	resp, bytes, errs := agent.Send(data).EndBytes()
	if len(errs) > 0 {
		return http.StatusInternalServerError, errs[0]
	}

	if err := json.Unmarshal(bytes, &dest); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "error while reading response content")
	}

	return resp.StatusCode, nil
}
