package request

import (
	"bytes"
	"fmt"
	"http-server/internal/headers"
	"io"
	"slices"
	"strconv"
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
	StateHeaders parserState = "headers"
	StateBody parserState = "body"
	StateError parserState = "error"
)

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	Body        string
	state       parserState

}

func getInt(headers *headers.Headers, name string, defaultVaule int) int {
	valueStr := headers.Get(name)


	if len(valueStr) == 0 {
		return defaultVaule
	}

	value, err := strconv.Atoi(valueStr)
	
	if err != nil {

		return defaultVaule
	}


	return value

}

func newRequest() *Request {
	return &Request{
		state: StateInit,
		Headers: headers.NewHeaders(),
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

func (r *Request) hasBody() bool {

	length := getInt(r.Headers, "content-length", 0)
		
	return length > 0
	
}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
outer:

	for {
       currentData := data[read:]

	//    slog.Info("Request#parse", "currentData", currentData, "data", data, "read", read)

	   if len(currentData) == 0{
		break
	   }
		switch r.state {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			
			r.state = StateHeaders

		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)

			if err != nil {
				r.state = StateError
				return 0, err
			}


			if n == 0 {
				break outer
			}
			read += n

			if done {

				if r.hasBody() {
					r.state = StateBody;

				} else {
					r.state = StateDone;
				}

			} 
		
		case StateBody:
			
			length := getInt(r.Headers, "content-length", 0)
			if length == 0{
				panic("length is zero")
			}
            //  slog.Info("RequestParse#", "currentData", currentData, "length", length, "lenBody", len(r.Body))
			remaining := min(length - len(r.Body), len(currentData))
           r.Body  += string(currentData[:remaining])
		
           read += remaining


			if len(r.Body) == length {
				r.state = StateDone
			}
		case StateDone:
			break outer
		default:
			panic("somehow we have programed prooly")
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
		
		readN, err := request.parse(buf[:bufLen])
		
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
