package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *RequestLine) ValidHTTP() bool {
	return r.HttpVersion == "1.1"
}

func (r *RequestLine) ValidMethod() bool {

	methods := []string{"GET", "POST", "PUT", "PATCH", "HEAD", "OPTIONS", "DELETE"}

	fmt.Print(slices.Contains(methods, r.Method))

	return slices.Contains(methods, r.Method)

}

type Request struct {
	RequestLine RequestLine
}

var ERROR_BAD_START_LINE = fmt.Errorf("bad start line")
var ERROR_BAD_HTTP_VERSION = fmt.Errorf("bad http version")
var ERROR_BAD_HTTP_METHOD = fmt.Errorf("bad http method")
var SEPARATOR = "\r\n"

func parseRequestLine(b string) (*RequestLine, string, error) {

	idx := strings.Index(b, SEPARATOR)

	if idx == -1 {
		return nil, b, nil
	}

	startLine := b[:idx]
	restOfMsg := b[idx+len(SEPARATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, b, ERROR_BAD_START_LINE
	}

	httpParts := strings.Split(parts[2], "/")



	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	if !rl.ValidHTTP()  {
		return nil, b, ERROR_BAD_HTTP_VERSION
	}

	if !rl.ValidMethod() {
		return nil, b, ERROR_BAD_HTTP_METHOD
	}

	return rl, restOfMsg, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {

	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("Unable to io.ReadAll:"), err,
		)
	}

	str := string(data)

	rl, str, err := parseRequestLine(str)

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *rl,
	}, nil
}
