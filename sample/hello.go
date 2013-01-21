package main

import (
	"fmt"
	"github.com/jkinner/goose"
	"os"
	"reflect"
)

func main() {
	injector := goose.CreateInjector()
	var name = "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	injector.BindToInstance(reflect.TypeOf(""), name)
	sayHello(injector.CreateContainer())
}

func sayHello(container goose.Container) {
	fmt.Println("Hello", container.GetInstance(reflect.TypeOf("")), "!")
}
