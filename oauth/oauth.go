package oauth

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/azzamjiul/bookstore_oauth-go/oauth/error_utils"
	"github.com/azzamjiul/bookstore_oauth-go/oauth/http_utils"
)

const (
	headerXPublic   = "X-Public"
	headerXClientId = "X-Client-Id"
	headerXCallerId = "X-Caller-Id"

	paramsAccessToken = "access_token"
)

type accessToken struct {
	Id       string `json:"id"`
	UserId   int64  `json:"user_id"`
	ClientId int64  `json:"client_id"`
	Expires  int64  `json:"expires"`
}

func IsPublic(request *http.Request) bool {
	if request == nil {
		return true
	}
	return request.Header.Get(headerXPublic) == "true"
}

func GetCallerId(request *http.Request) int64 {
	if request == nil {
		return 0
	}

	callerId, err := strconv.ParseInt(request.Header.Get(headerXCallerId), 10, 64)
	if err != nil {
		return 0
	}

	return callerId
}

func GetClientId(request *http.Request) int64 {
	if request == nil {
		return 0
	}

	clientId, err := strconv.ParseInt(request.Header.Get(headerXClientId), 10, 64)
	if err != nil {
		return 0
	}

	return clientId
}

func AuthenticateRequest(request *http.Request) *error_utils.RestErr {
	if request == nil {
		return nil
	}

	cleanRequest(request)

	// htt://api.bookstore.com/resource?access_token=abc123
	accessToken := strings.TrimSpace(request.URL.Query().Get(paramsAccessToken))
	if accessToken == "" {
		return nil
	}

	at, err := getAccessToken(accessToken)
	if err != nil {
		if err.Status == http.StatusNotFound {
			return nil
		}
		return err
	}

	request.Header.Add(headerXCallerId, fmt.Sprintf("%v", at.UserId))
	request.Header.Add(headerXClientId, fmt.Sprintf("%v", at.ClientId))

	return nil
}

func cleanRequest(request *http.Request) {
	if request == nil {
		return
	}

	request.Header.Del(headerXClientId)
	request.Header.Del(headerXCallerId)
}

func getAccessToken(access_token string) (*accessToken, *error_utils.RestErr) {
	var at accessToken
	headers := map[string]string{"Content-Type": "application/json"}

	_, err := http_utils.New().Get(fmt.Sprintf("http://localhost:8081/oauth/access_token/%s", access_token), headers, &at)
	if err != nil {
		return nil, error_utils.NewInternalServerError(err.Error())
	}

	return &at, nil
}
