package main

import (
	"flag";
	"fmt";
	"github.com/jkinner/goose";
	"reflect";
)

var name *string = flag.String("name", "world", "whom to say hello to")

func main() {
	flag.Parse()
	injector := goose.CreateInjector()
	injector.Bind(reflect.TypeOf(""), func () (interface{}) { return *name })
	sayHello(injector.CreateContainer())
}

func sayHello(container goose.Container) {
	fmt.Println("Hello", container.GetInstance(reflect.TypeOf("")), "!")
}
