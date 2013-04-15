/*
 * Copyright 2013 Google Inc. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *     http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
