package task

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"neon/build"
	"net/http"
	"reflect"
)

const (
	DEFAULT_METHOD = "GET"
	DEFAULT_STATUS = 200
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "request",
		Func: Request,
		Args: reflect.TypeOf(RequestArgs{}),
		Help: `Perform an HTTP request.

Arguments:

- request: URL to request (string).
- method: request method ('GET', 'POST', etc), defaults to 'GET' (string,
  optional).
- headers: request headers (map with string keys and values, optional).
- body: request body (string, optional).
- file: request body as a file (string, optional, file).
- username: user name for authentication (string, optional).
- password: user password for authentication (string, optional).
- status: expected status code, raise an error if different, defaults to 200
  (int, optional).

Examples:

    # get google.com
    - request: 'google.com'

Notes:

- Response status code is stored in variable _status.
- Response body is stored in variable _body.
- Response headers are stored in variable _headers.`,
	})
}

type RequestArgs struct {
	Request  string
	Method   string            `optional`
	Headers  map[string]string `optional`
	Body     string            `optional`
	File     string            `optional file`
	Username string            `optional`
	Password string            `optional`
	Status   int               `optional`
}

func Request(context *build.Context, args interface{}) error {
	params := args.(RequestArgs)
	var err error
	method := params.Method
	if method == "" {
		method = DEFAULT_METHOD
	}
	status := params.Status
	if status == 0 {
		status = DEFAULT_STATUS
	}
	body := []byte(params.Body)
	if params.File != "" {
		body, err = ioutil.ReadFile(params.File)
		if err != nil {
			return err
		}
	}
	request, err := http.NewRequest(method, params.Request, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return fmt.Errorf("building request: %v", err)
	}
	for name, value := range params.Headers {
		request.Header.Set(name, value)
	}
	if params.Username != "" {
		request.SetBasicAuth(params.Username, params.Password)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("requesting '%s': %v", params.Request, err)
	}
	defer response.Body.Close()
	context.SetProperty("_status", response.StatusCode)
	context.SetProperty("_headers", response.Header)
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	context.SetProperty("_body", string(responseBody))
	if response.StatusCode != status {
		return fmt.Errorf("bad response status: %d", response.StatusCode)
	}
	return nil
}
