package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	echov1 "github.com/organic-programming/grace-op/internal/cli/testsupport/echoholon/protos/echo/v1"
	"google.golang.org/grpc"
)

const defaultListenURI = "tcp://127.0.0.1:0"

type server struct {
	echov1.UnimplementedEchoServiceServer
}

func (server) Ping(_ context.Context, request *echov1.PingRequest) (*echov1.PingResponse, error) {
	message := request.GetMessage()
	switch request.GetMode() {
	case echov1.EchoMode_ECHO_MODE_UPPER:
		message = strings.ToUpper(message)
	case echov1.EchoMode_ECHO_MODE_LOWER:
		message = strings.ToLower(message)
	}

	return &echov1.PingResponse{
		Message: message,
		Count:   int32(len(request.GetTags())),
		Mode:    request.GetMode(),
	}, nil
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
	}

	listen := flag.String("listen", defaultListenURI, "tcp URI to listen on")
	flag.Parse()

	listener, publicURI, err := listenTCP(*listen)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen failed: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	echov1.RegisterEchoServiceServer(grpcServer, server{})

	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- grpcServer.Serve(listener)
	}()

	fmt.Println(publicURI)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(sigCh)

	select {
	case <-sigCh:
		shutdown(grpcServer)
	case err := <-serveErrCh:
		if err != nil && !strings.Contains(strings.ToLower(err.Error()), "use of closed network connection") {
			fmt.Fprintf(os.Stderr, "serve failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := <-serveErrCh; err != nil && !strings.Contains(strings.ToLower(err.Error()), "use of closed network connection") {
		fmt.Fprintf(os.Stderr, "serve failed: %v\n", err)
		os.Exit(1)
	}
}

func listenTCP(uri string) (net.Listener, string, error) {
	if !strings.HasPrefix(uri, "tcp://") {
		return nil, "", fmt.Errorf("unsupported listen URI %q", uri)
	}

	address := strings.TrimPrefix(uri, "tcp://")
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, "", err
	}

	host, _, err := net.SplitHostPort(address)
	if err != nil || host == "" || host == "0.0.0.0" || host == "::" || host == "[::]" {
		host = "127.0.0.1"
	}
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return nil, "", err
	}

	return listener, fmt.Sprintf("tcp://%s:%s", host, port), nil
}

func shutdown(server *grpc.Server) {
	done := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		server.Stop()
	}
}
