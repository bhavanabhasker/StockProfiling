package main

import (
	"golang-book/assign1http/api"
	"log"
	"net/http"
	"os"

	"github.com/bakins/net-http-recover"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/justinas/alice"
)

func main() {

	r := mux.NewRouter()
	// Create the server and register the service
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")

	s.RegisterService(new(api.Api), "")
	//use alice for http handling
	chain := alice.New(
		func(h http.Handler) http.Handler {
			return handlers.CombinedLoggingHandler(os.Stdout, h)
		},
		// to support compressed responses
		handlers.CompressHandler,
		func(h http.Handler) http.Handler {
			return recovery.Handler(os.Stderr, h, true)
		})

	r.Handle("/rpc", chain.Then(s))
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
