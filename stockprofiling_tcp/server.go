package main

import (
	"golang-book/assign1/api"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	// create new rpc server and register service
	server := rpc.NewServer()
	server.Register(new(api.Api))
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error: ", e)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
