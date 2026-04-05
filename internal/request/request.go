package request

import (
	"bytes"
	"fmt"
	"io"
	"slices"
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

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
	StateError parserState = "error"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

var ERROR_BAD_START_LINE = fmt.Errorf("bad start line")
var ERROR_BAD_HTTP_VERSION = fmt.Errorf("bad http version")
var ERROR_BAD_HTTP_METHOD = fmt.Errorf("bad http method")
var ERROR_REQUEST_IN_ERROR_STATE = fmt.Errorf("request in error state")
var SEPARATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {

	idx := bytes.Index(b, SEPARATOR)

	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_BAD_START_LINE
	}

	httpParts := bytes.Split(parts[2], []byte("/"))

	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {

		return nil, 0, ERROR_BAD_START_LINE
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	if !rl.ValidHTTP() {
		return nil, 0, ERROR_BAD_HTTP_VERSION
	}

	if !rl.ValidMethod() {
		return nil, 0, ERROR_BAD_HTTP_METHOD
	}

	return rl, read, nil

}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
outer:
	for {

		switch r.state {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n

			r.state = StateDone
		case StateDone:
			break outer
		}

	}
	return read, nil
}

func (r *Request) isDone() bool {
	return r.state == StateDone
}

func (r *Request) isError() bool {
	return r.state == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	request := newRequest()

	buf := make([]byte, 1024)
	bufLen := 0

	for !request.isDone() && !request.isError() {

		n, err := reader.Read(buf[bufLen:])

		if err != nil {
			return nil, err
		}

		bufLen += n

		readN, err := request.parse(buf[:bufLen+n])

		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN

	}
	// data, err := io.ReadAll(reader)

	// if err != nil {
	// 	return nil, errors.Join(
	// 		fmt.Errorf("Unable to io.ReadAll:"), err,
	// 	)
	// }

	// str := string(data)

	// rl, _, err := parseRequestLine(str)

	// if err != nil {
	// 	return nil, err
	// }

	// return &Request{
	// 	RequestLine: *rl,
	// }, nil

	return request, nil
}
