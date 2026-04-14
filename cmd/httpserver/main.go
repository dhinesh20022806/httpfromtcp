package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"http-server/internal/headers"
	"http-server/internal/request"
	"http-server/internal/response"
	"http-server/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func toStr(bytes []byte) string {

	out := ""

	for _, b := range bytes {

		out += fmt.Sprintf("%02x", b)

	}

	return out;
}

func response400() []byte{

	return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}

func response500() []byte{

	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}

func response200() []byte{

	return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
}

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {

		h := response.GetDefaultHeaders(0)
		body := response200()
		status := response.StatusOK
        
		if req.RequestLine.RequestTarget == "/yourproblem" {

			body = response400()
			status = response.StatusBadRequest
       
		}

		if req.RequestLine.RequestTarget == "/myproblem" {

			body = response500()
			status = response.StatusInternalServerError			
		}

		if req.RequestLine.RequestTarget == "/video" {
         
			f, _ := os.ReadFile("assets/vim.mp4")
            	h.Replace("content-type", "video/mp4")
				h.Replace("content-length", fmt.Sprintf("%d", len(f)))
				w.WriteStatusLine(response.StatusOK)
				w.WriteHeaders(h)
				w.WriteBody(f)

          return;

		}

		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream/") {

			target := req.RequestLine.RequestTarget
			res, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])

			if err != nil {

				body = response500()
				status = response.StatusInternalServerError
			}else {

				w.WriteStatusLine(response.StatusOK)
				h.Delete("content-length")
				h.Set("transfer-encoding", "chunked")
				h.Replace("content-type", "text/plain")
				h.Set("trailer", "X-Content-SHA256")
				h.Set("trailer", "X-Content-Length")
				w.WriteHeaders(h)


				fullBody := []byte{}
				for{
					data := make([]byte, 32)

					n, err := res.Body.Read(data)

					if err != nil {
						if errors.Is(io.EOF, err){
							break
						}
						break
					}

					fullBody = append(fullBody, data[:n]...)

					w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))
					w.WriteBody(data[:n])
					w.WriteBody([]byte("\r\n"))
				}

				w.WriteBody([]byte("0\r\n\r\n"))
				trailers := headers.NewHeaders()


				out := sha256.Sum256(fullBody)

				trailers.Set("X-Content-SHA256",toStr(out[:]))
				trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(string(fullBody))) )

				return;


			}

			

		} 


		    w.WriteStatusLine(status)
			h.Replace("content-type", "text/html")
			h.Replace("content-length", fmt.Sprintf("%d", len(body)))
			w.WriteHeaders(h)
			w.WriteBody(body)

		w.Writer.Write([]byte("All good, frfr\n"))

		
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}