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
	build.AddTask(build.TaskDesc {
		Name: "request",
		Func: Request,
		Args: reflect.TypeOf(RequestArgs{}),
		Help: `Perform an HTTP request.

Arguments:

- request: the URL to request.
- method: the request method (GET, POST, etc), defaults to "GET".
- headers: request headers as an anko map.
- body: the request body as a string.
- file: the request body as a file.
- status: expected status code, on error if different (defaults to 200).
- username: user name for authentication.
- password: user password for authentication.

Response status code is stored in variable _status, response body is stored in
variable _body and response headers in _headers.

Examples:

    # get google.com
    - request: "google.com"`,
	})
}

type RequestArgs struct {
	Request  string
	Method   string            `optional`
	Headers  map[string]string `optional`
	Body     string            `optional`
	File     string            `optional file`
	Status   int               `optional`
	Username string            `optional`
	Password string            `optional`
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
	request, err := http.NewRequest(params.Method, params.Request, bytes.NewBuffer([]byte(body)))
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
	if response.StatusCode != params.Status {
		return fmt.Errorf("bad response status: %s", response.StatusCode)
	}
	return nil
}
