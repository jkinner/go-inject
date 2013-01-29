package main

import (
	"flag"
	"fmt"
	"github.com/jkinner/goose"
	"github.com/jkinner/goose/http"
	"log"
	"net/http"
)

type UserId struct{}

func ConfigureInjector(injector goose.Injector) {
	handlers := make(goose_http.HandlerMap)
	handlerFuncs := make(goose_http.HandlerFuncMap)

	i := 0
	injector.BindInScope(UserId{}, func(_ goose.Context, _ goose.Container) interface{} {
		i += 1
		return i
	}, goose_http.RequestScoped{})

	handlerFuncs["/"] = func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add(
			"Content-Type",
			"text/plain",
		)
		w.Write([]byte(fmt.Sprintf("Hello, %d!\n", injector.CreateContainer().GetInstance(request, UserId{}).(int))))
		w.Write([]byte(fmt.Sprintf("Hello, %d!", injector.CreateContainer().GetInstance(request, UserId{}).(int))))
	}

	injector.BindInstance(goose_http.Handlers{}, handlers)
	injector.BindTaggedInstance(goose_http.Handlers{}, goose_http.Func{}, handlerFuncs)
}

type OneServer struct{}
type TwoServer struct{}

func main() {
	flag.Parse()
	injector := goose.CreateInjector()

	goose_http.ConfigureScopes(injector)

	// Create a child injector for each server
	oneServerInjector := injector.CreateChildInjector()
	twoServerInjector := injector.CreateChildInjector()
	// Don't use flags; that would confuse both servers.
	oneServerInjector.BindInstance(goose_http.Port{}, 8080)
	goose_http.ConfigureInjector(oneServerInjector)
	ConfigureInjector(oneServerInjector)
	oneServerInjector.ExposeAndTag(goose_http.Server{}, OneServer{})

	twoServerInjector.BindInstance(goose_http.Port{}, 8081)
	ConfigureInjector(twoServerInjector)
	goose_http.ConfigureInjector(twoServerInjector)
	twoServerInjector.ExposeAndTag(goose_http.Server{}, TwoServer{})

	container := injector.CreateContainer()

	oneHttpServer := container.GetTaggedInstance(nil, goose_http.Server{}, OneServer{}).(http.Server)
	go func() {
		log.Fatal(oneHttpServer.ListenAndServe())
	}()

	twoHttpServer := container.GetTaggedInstance(nil, goose_http.Server{}, TwoServer{}).(http.Server)
	log.Fatal(twoHttpServer.ListenAndServe())
}
