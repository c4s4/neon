package task

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"neon/build"
	"neon/util"
	"net/http"
	"strconv"
)

const (
	DEFAULT_METHOD = "GET"
	DEFAULT_STATUS = "200"
)

func init() {
	build.TaskMap["request"] = build.TaskDescriptor{
		Constructor: Request,
		Help: `Perform an HTTP request.

Arguments:

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
	}
}

func Request(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"request", "method", "headers", "body", "file", "status", "username", "password"}
	if err := CheckFields(args, fields, fields[:1]); err != nil {
		return nil, err
	}
	url, err := args.GetString("request")
	if err != nil {
		return nil, fmt.Errorf("argument request must be a string")
	}
	var method = DEFAULT_METHOD
	if args.HasField("method") {
		method, err = args.GetString("method")
		if err != nil {
			return nil, fmt.Errorf("argument method of task request must be a string")
		}
	}
	var headers string
	if args.HasField("headers") {
		headers, err = args.GetString("headers")
		if err != nil {
			return nil, fmt.Errorf("argument headers of task request must be a string")
		}
	}
	var body string
	if args.HasField("body") {
		body, err = args.GetString("body")
		if err != nil {
			return nil, fmt.Errorf("argument body of task request must be a string")
		}
	}
	var file string
	if args.HasField("file") {
		file, err = args.GetString("file")
		if err != nil {
			return nil, fmt.Errorf("argument file of task request must be a string")
		}
	}
	var status string
	if args.HasField("status") {
		status, err = args.GetString("status")
		if err != nil {
			return nil, fmt.Errorf("argument status of task request must be a string")
		}
	}
	var username string
	if args.HasField("username") {
		username, err = args.GetString("username")
		if err != nil {
			return nil, fmt.Errorf("argument username of task request must be a string")
		}
	}
	var password string
	if args.HasField("password") {
		password, err = args.GetString("password")
		if err != nil {
			return nil, fmt.Errorf("argument password of task request must be a string")
		}
	}
	if file != "" && body != "" {
		return nil, fmt.Errorf("body and file can't be set at the same time")
	}
	return func() error {
		// evaluate arguments
		_url, _err := target.Build.Context.EvaluateString(url)
		if _err != nil {
			return fmt.Errorf("evaluating url: %v", _err)
		}
		_method, _err := target.Build.Context.EvaluateString(method)
		if _err != nil {
			return fmt.Errorf("evaluating method: %v", _err)
		}
		_result, _err := target.Build.Context.EvaluateExpression(headers)
		if _err != nil {
			return fmt.Errorf("evaluating headers: %v", _err)
		}
		_headers, _err := util.ToMapStringString(_result)
		if _err != nil {
			return fmt.Errorf("evaluating headers: %v", _err)
		}
		_str, _err := target.Build.Context.EvaluateString(body)
		if _err != nil {
			return fmt.Errorf("evaluating body: %v", _err)
		}
		_body := []byte(_str)
		_file, _err := target.Build.Context.EvaluateString(file)
		if _err != nil {
			return fmt.Errorf("evaluating file: %v", _err)
		}
		if _file != "" {
			_file = util.ExpandAndJoinToRoot(target.Build.Dir, _file)
		}
		_status, _err := target.Build.Context.EvaluateString(status)
		if _err != nil {
			return fmt.Errorf("evaluating status: %v", _err)
		}
		if _status == "" {
			_status = DEFAULT_STATUS
		}
		_username, _err := target.Build.Context.EvaluateString(username)
		if _err != nil {
			return fmt.Errorf("evaluating username: %v", _err)
		}
		_password, _err := target.Build.Context.EvaluateString(password)
		if _err != nil {
			return fmt.Errorf("evaluating password: %v", _err)
		}
		// perform request
		if _file != "" {
			_body, _err = ioutil.ReadFile(_file)
			if _err != nil {
				return _err
			}
		}
		_request, _err := http.NewRequest(_method, _url, bytes.NewBuffer([]byte(_body)))
		if _err != nil {
			return fmt.Errorf("building request: %v", _err)
		}
		for _name, _value := range _headers {
			_request.Header.Set(_name, _value)
		}
		if _username != "" {
			_request.SetBasicAuth(_username, _password)
		}
		client := &http.Client{}
		_response, _err := client.Do(_request)
		if _err != nil {
			return fmt.Errorf("requesting '%s': %v", _url, _err)
		}
		defer _response.Body.Close()
		_response_status := strconv.Itoa(_response.StatusCode)
		target.Build.Context.SetProperty("_status", _response_status)
		target.Build.Context.SetProperty("_headers", _response.Header)
		_response_body, _err := ioutil.ReadAll(_response.Body)
		if _err != nil {
			return fmt.Errorf("reading response body: %v", _err)
		}
		target.Build.Context.SetProperty("_body", string(_response_body))
		if _response_status != _status {
			return fmt.Errorf("bad response status: %s", _response_status)
		}
		return nil
	}, nil
}
