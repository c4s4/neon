package task

import (
	"bytes"
	"fmt"
	"github.com/c4s4/neon/build"
	"io/ioutil"
	"net/http"
	"reflect"
)

const (
	// DefaultMethod is the default request method
	DefaultMethod = "GET"
	// DefaultStatus is the default expected response status
	DefaultStatus = 200
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "request",
		Func: request,
		Args: reflect.TypeOf(requestArgs{}),
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

type requestArgs struct {
	Request  string
	Method   string            `neon:"optional"`
	Headers  map[string]string `neon:"optional"`
	Body     string            `neon:"optional"`
	File     string            `neon:"optional,file"`
	Username string            `neon:"optional"`
	Password string            `neon:"optional"`
	Status   int               `neon:"optional"`
}

func request(context *build.Context, args interface{}) error {
	params := args.(requestArgs)
	var err error
	method := params.Method
	if method == "" {
		method = DefaultMethod
	}
	status := params.Status
	if status == 0 {
		status = DefaultStatus
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
