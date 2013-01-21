package main

import (
	"fmt"
	"github.com/jkinner/goose"
	"os"
	"reflect"
)

// Injector tags
type HelloString struct{}
type Name struct{}

func ConfigureInjector(injector goose.Injector) {
	injector.BindTagged(reflect.TypeOf(""), HelloString{},
		func(container goose.Container) interface{} {
			return fmt.Sprintf(
				"Hello, %s!", container.GetTaggedInstance(reflect.TypeOf(""), Name{}))
		})
	var name = "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	injector.BindTaggedInstance(reflect.TypeOf(""), Name{}, name)
}

func main() {
	injector := goose.CreateInjector()
	ConfigureInjector(injector)
	sayHello(injector.CreateContainer())
}

func sayHello(container goose.Container) {
	fmt.Println(container.GetTaggedInstance(reflect.TypeOf(""), HelloString{}))
}
