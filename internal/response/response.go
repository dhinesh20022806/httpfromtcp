package response

import (
	"fmt"
	"http-server/internal/headers"
	"io"
	"log/slog"
	"strings"
)


type Response struct {

}

type Writer struct {
	Writer io.Writer

}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{Writer: writer}
}



type StatusCode int


const (
	StatusOK StatusCode = 200
	
	StatusBadRequest = 400
	StatusInternalServerError = 500

)

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	rn := "\r\n";

	header := []string{};

	
header = append(header, rn)
	headers.ForEach(func(n string, k string) {
            header = append(header,  n + ":" + k + rn)
	});

	header = append(header, rn);

	slog.Info("WriteHeaders#func", "header", strings.Join(header, ""))

	byteHeader := []byte(strings.Join(header, ""))

	_, err := w.Writer.Write(byteHeader);

	return err;
}

func GetDefaultHeaders(contentLength int) headers.Headers {

	headers := headers.NewHeaders();

	headers.Set("content-length", fmt.Sprintf("%d", contentLength))
	headers.Set("content-type", "text/plain")
	headers.Set("connection", "closed")

	return *headers
}

func (w *Writer)WriteStatusLine(statusCode StatusCode) error {

	statusLine := []byte{}

	switch statusCode {
	case StatusOK: statusLine = []byte("HTTP/1.1 200 OK")
	case StatusBadRequest: statusLine = []byte("HTTP/1.1 400 Bad Request")
	case StatusInternalServerError: statusLine = []byte("HTTP/1.1 500 Internal Server Error")	

	default:
		return fmt.Errorf("unrecognized status code")

	}

	_, err := w.Writer.Write(statusLine);

	return err

}

func (w *Writer) WriteBody(p []byte) (int, error){

	n, err := w.Writer.Write(p);

	if err != nil {
		return 0, err
	}
	return n , nil
}