package main

import (
	"fmt"
	"github.com/jkinner/goose"
	"os"
)

// Injector tags
type HelloString struct{}
type Name struct{}

func ConfigureInjector(injector goose.Injector) {
	injector.Bind(HelloString{},
		func(context goose.Context, container goose.Container) interface{} {
			return fmt.Sprintf(
				"Hello, %s!", container.GetInstance(context, Name{}))
		})
	var name = "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	injector.BindInstance(Name{}, name)
}

func main() {
	injector := goose.CreateInjector()
	ConfigureInjector(injector)
	sayHello(injector.CreateContainer())
}

func sayHello(container goose.Container) {
	fmt.Println(container.GetInstance(nil, HelloString{}))
}
