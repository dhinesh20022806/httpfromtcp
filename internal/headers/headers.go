package headers

import (
	"bytes"
	"fmt"
)
type Headers struct {
	headers map[string]string
}

var rn = []byte("\r\n")


func NewHeaders() *Headers {
	return &Headers { map[string]string{}}
}

func parseHeader(fieldLine []byte) (string, string, error){
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	fmt.Print(string(value), "value in string")

	if bytes.HasSuffix(name, []byte(" ")){
		return "", "", fmt.Errorf("malformed field name")
	}

	return string(name) , string(value), nil

}

func (h Headers) Parse(data []byte) (int,  bool,  error){

	read := 0
    isDone := false
	for {
		idx := bytes.Index(data[read:], rn)

		if idx == -1 {
			break
		}

		// /r/n/r/n last /r/n is 0 index
		if idx == 0 {
			isDone = true
			read += len(rn)
			break
		}

		name, value, err := parseHeader(data[read:read + idx])

		if err != nil {

			return 0, false, err
			
		}
		fmt.Print(idx, "  idx", read, "  read")

		read += idx + len(rn)
		h[name] = value


	}

	return read, isDone, nil
} 


