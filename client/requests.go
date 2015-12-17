import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	errEmptyUrl = errors.New("Url for request is empty")
	errUndefinedMethod = errors.New("Method is undefined")
	errEmptyItem = errors.New("Item is empty")
)

// Helping function for getting requests to the server
func request(url, method string, item interface{}) (*http.Response, error) {
	if url == "" {
		return nil, errEmptyItem
	}

	if method != "GET" && method != "POST" {
		return nil, errUndefinedMethod
	}

	if item == nil {
		return nil, errEmptyItem
	}

	body, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return req, err
	}
	req.Header.Set("content-type", "application/json")
	return req, err
}

func unmarshal(resp http*Response, respitem interface{}) error {
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err := json.Unmarshal(respBody)
	if err != nil {
		return err
	}

	return nil
}