package main

import (
	"fmt"
	"github.com/jkinner/goose"
	"os"
	"reflect"
)

/*
	This is sort of how Guice works. But it's not really go-ish.
*/
func main() {
	injector := goose.CreateInjector()
	var name = "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	injector.BindInstance(reflect.TypeOf(""), name)
	sayHello(injector.CreateContainer())
}

func sayHello(container goose.Container) {
	fmt.Println(fmt.Sprintf("Hello, %s!", container.GetInstance(nil, reflect.TypeOf(""))))
}
