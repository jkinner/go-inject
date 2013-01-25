package main

import (
	"flag"
	"fmt"
	"github.com/jkinner/goose"
)

var name *string = flag.String("name", "world", "whom to say hello to")

// Injector key
type Name struct{}

/*
	Bind using an empty struct instead of the type. This is more idomatic go because of the
	different in reflection between go and Java.
*/
func main() {
	flag.Parse()
	injector := goose.CreateInjector()
	// Bind flag
	injector.Bind(Name{}, func(container goose.Container) interface{} { return *name })
	sayHello(injector.CreateContainer())
}

func sayHello(container goose.Container) {
	fmt.Println(fmt.Sprintf("Hello, %s!", container.GetInstance(Name{})))
}
