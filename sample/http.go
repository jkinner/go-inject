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

func sayFoo(w http.ResponseWriter, request *http.Request) {
	w.Header().Add(
		"Content-Type",
		"text/plain",
	)
	w.Write([]byte(fmt.Sprintf("Foo\n")))
}

func ConfigureInjector(injector goose.Injector) {
	i := 0
	injector.BindInScope(Counter{}, func(_ goose.Context, _ goose.Container) interface{} {
		i += 1
		return i
	}, goose_http.RequestScoped{})

	goose_http.BindHandlerFunc(injector, "/",
		func (w http.ResponseWriter, request *http.Request) {
			w.Header().Add(
				"Content-Type",
				"text/plain",
			)
			w.Write([]byte(fmt.Sprintf("Hello! (%d)\n",
				injector.CreateContainer().GetInstance(request, Counter{}).(int))))
			w.Write([]byte(fmt.Sprintf("Hello! (%d)\n",
				injector.CreateContainer().GetInstance(request, Counter{}).(int))))
		})
}

func ConfigureFooInjector(injector goose.Injector) {
	goose_http.BindHandlerFunc(injector, "/foo/", sayFoo)
}

func main() {
	flag.Parse()
	injector := goose.CreateInjector()
	goose_http.ConfigureFlags(injector)
	goose_http.ConfigureScopes(injector)
	goose_http.ConfigureInjector(injector)
	ConfigureFooInjector(injector)
	ConfigureInjector(injector)
	container := injector.CreateContainer()
	httpServer := container.GetInstance(nil, goose_http.Server{}).(http.Server)
	log.Fatal(httpServer.ListenAndServe())
}
