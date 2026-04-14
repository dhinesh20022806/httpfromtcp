package server

import (
	"fmt"
	"http-server/internal/request"
	"http-server/internal/response"
	"io"
	"log/slog"
	"net"
)

type Server struct {
	closed bool
	handler Handler

}
type HandlerError struct {
	StatusCode response.StatusCode
	Message string
}



type Handler func(w *response.Writer, req *request.Request);

func runConnection(s *Server, conn io.ReadWriteCloser){
    defer conn.Close()

	responsewriter := response.NewWriter(conn)

	headers := response.GetDefaultHeaders(0)

	
    r, err := request.RequestFromReader(conn)
	
	if err != nil {
		responsewriter.WriteStatusLine(response.StatusBadRequest);
		responsewriter.WriteHeaders(headers);
		return
	}

	s.handler(responsewriter, r)



	 
}

func runServer(s *Server, listener net.Listener){

	for {
		conn, err := listener.Accept()
		fmt.Println("Server Conn accpeted")
		if s.closed {
			return;
		}
		if err != nil {
			return


		}


      go runConnection(s, conn)
	}



}

func Serve(port uint16, handler Handler) (*Server, error){

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{ closed: false, 
	 handler: handler,
	}
	go runServer(server, listener)


	return server, nil
}

func (s Server) Close() error {

	slog.Info("Server#Closes", "closes", s.closed)
	s.closed = true;
	return nil
}

