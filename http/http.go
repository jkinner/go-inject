package goose_http

import (
	"flag"
	"fmt"
	"github.com/jkinner/goose"
	"log"
	"net/http"
	"reflect"
)

// Value for the http_port flag (default: 80)
var httpPort *int = flag.Int("http_port", 80, "port on which to start the http listener")

type Port struct{}
type Handlers struct{}
// Tag on the Handlers key (specifying that the handlers are func handlers)
type Func struct{}

// Scope tag for the caching in the context of the *http.Request
type RequestScoped struct{}

type HandlerMap map[string]http.Handler
type HandlerFuncMap map[string]func(http.ResponseWriter, *http.Request)

var requestScope = goose.CreateSimpleScopeWithName("HTTP Request")

type scopingHandler struct {
	handler http.Handler
}

func (this scopingHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer requestScope.Exit(request)
	requestScope.Enter(request)
	this.handler.ServeHTTP(writer, request)
}

func providesHttpServer(_ goose.Context, container goose.Container) interface{} {
	port := container.GetInstance(nil, Port{}).(int)

	serveMux := http.NewServeMux()
	handlers := container.GetInstance(nil, Handlers{}).(HandlerMap)
	for path, handler := range(handlers) {
		serveMux.Handle(path, scopingHandler { handler })
	}

	handlerFuncs := container.GetTaggedInstance(nil, Handlers{}, Func{}).(HandlerFuncMap)
	for path, handlerFunc := range(handlerFuncs) {
		serveMux.HandleFunc(path, func(w http.ResponseWriter, request *http.Request) {
			defer requestScope.Exit(request)
			requestScope.Enter(request)
			handlerFunc(w, request)
		})
	}

	log.Printf("Creating HTTP server listening on port %d", port)
	return http.Server {
		Addr:			fmt.Sprintf(":%d", port),
		Handler:	serveMux,
	}
}

// Binds the following:
//   goose_http.Port - the value of the http_port flag
func ConfigureFlags(injector goose.Injector) {
	injector.Bind(Port{}, func(_ goose.Context, _ goose.Container) interface{} { return *httpPort })
}

// Binds the following:
//   reflect.TypeOf(http.Server{}) - the HTTP server itself
// Requires these bindings:
//   goose_http.Handlers - a HandlerMap assigning a path to a http.Handler
//   goose_http.Handlers<goose_http.Func> - a HanderFuncMap assigning a path to a handler func
func ConfigureInjector(injector goose.Injector) {
	injector.BindScope(requestScope, RequestScoped{})
	injector.Bind(reflect.TypeOf(http.Server{}), providesHttpServer)
}
