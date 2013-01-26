package main

import (
	"flag"
	"fmt"
	"github.com/jkinner/goose"
	"github.com/jkinner/goose/http"
	"log"
	"net/http"
	"reflect"
)

type UserId struct{}

func ConfigureInjector(injector goose.Injector) {
	handlers := make(goose_http.HandlerMap)
	handlerFuncs := make(goose_http.HandlerFuncMap)

	i := 0
	injector.BindInScope(UserId{}, func(_ goose.Context, _ goose.Container) interface{} {
		fmt.Println("Creating a new UserId")
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

func main() {
	flag.Parse()
	injector := goose.CreateInjector()
	goose_http.ConfigureFlags(injector)
	goose_http.ConfigureInjector(injector)
	ConfigureInjector(injector)
	container := injector.CreateContainer()
	httpServer := container.GetInstance(nil, reflect.TypeOf(http.Server{})).(http.Server)
	log.Fatal(httpServer.ListenAndServe())
}
