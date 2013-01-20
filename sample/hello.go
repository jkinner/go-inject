package main

import (
	"fmt";
	"github.com/jkinner/goose";
	"os";
	"reflect";
)

func main() {
	injector := goose.CreateInjector()
	var name = "world"
	if len(os.Args) > 1 {
  	name = os.Args[1]
	}

	injector.BindToInstance(goose.CreateKeyForType(reflect.TypeOf("")), name)
	sayHello(injector)
}

func sayHello(injector goose.Injector) {
	fmt.Println("Hello", injector.GetInstance(reflect.TypeOf("")), "!")
}
