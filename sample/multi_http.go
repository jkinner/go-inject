package main

import (
	"flag"
	"fmt"
	"github.com/jkinner/goose"
	"github.com/jkinner/goose/http"
	"log"
	"net/http"
)

type Counter struct{}

// Configure a counter provider that increments on each invocation.
func ConfigureCounter(injector goose.Injector) {
	i := 0
	injector.BindInScope(Counter{}, func(_ goose.Context, _ goose.Container) interface{} {
		i += 1
		return i
	}, goose_http.RequestScoped{})
}

// Configure HTTP server path mappings and provider functions.
func ConfigureInjector(injector goose.Injector) {

	goose_http.BindHandlerFunc(
		injector,
		"/",
		func(w http.ResponseWriter, request *http.Request) {
			w.Header().Add(
				"Content-Type",
				"text/plain",
			)
			w.Write([]byte(fmt.Sprintf("Hello, %d!\n", injector.CreateContainer().GetInstance(request, Counter{}).(int))))
			w.Write([]byte(fmt.Sprintf("Hello, %d!", injector.CreateContainer().GetInstance(request, Counter{}).(int))))
		})
}

type OneServer struct{}
type TwoServer struct{}

// Configure a single HTTP server on a port and expose it using the tag.
func ConfigureServer(injector goose.Injector, port int, tag goose.Tag) {
	// Each HTTP server needs to have a port bound. The port is statically configured here, but
	// it could be configured using a Provider that allocates an unused port.
	injector.BindInstance(goose_http.Port{}, port)
	goose_http.ConfigureInjector(injector)
	ConfigureInjector(injector)
	injector.ExposeAndTag(goose_http.Server{}, tag)
}

func main() {
	// Initialize the flag(s)
	flag.Parse()

	// Create the injector to begin configuring the bindings.
	injector := goose.CreateInjector()

	// Don't configure flags, since the HTTP module assumes a single HTTP port in use.

	// Bind any scope tags used in the goose_http module. Scope bindings are in a separate function
	// so the module can be used more than once. The scope binding can only ever happen once.
	goose_http.ConfigureScopes(injector)

	// Use a global counter. Both HTTP servers will increment at the same time.
	ConfigureCounter(
		injector,
	)

	// Create a child injector for each server. Each child injector is isolated, so the HTTP module
	// can be bound multiple times. Only the exposed bindings have the potential to collide. Note
	// that if either server were bound in the parent injector, all the bindings would collide.
	oneServerInjector := injector.CreateChildInjector()
	twoServerInjector := injector.CreateChildInjector()

	// Configure each server on its own port, exposing the server in the parent injector with its
	// own tag.
	ConfigureServer(oneServerInjector, 8080, OneServer{})
	ConfigureServer(twoServerInjector, 8081, TwoServer{})

	oneHttpServer := injector.CreateContainer().GetTaggedInstance(
		nil,
		goose_http.Server{},
		OneServer{}).(http.Server)

	// What would I do with go-routines? Two HTTP servers at the same time.
	go func() {
		log.Fatal(oneHttpServer.ListenAndServe())
	}()

	twoHttpServer := injector.CreateContainer().GetTaggedInstance(
		nil,
		goose_http.Server{},
		TwoServer{}).(http.Server)
	log.Fatal(twoHttpServer.ListenAndServe())
}
