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
	"flag"
	"fmt"
	"code.google.com/p/go-inject"
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
	injector := inject.CreateInjector()
	// Bind flag
	injector.Bind(Name{}, func(_ inject.Context, _ inject.Container) interface{} { return *name })
	sayHello(injector.CreateContainer())
}

func sayHello(container inject.Container) {
	fmt.Println(fmt.Sprintf("Hello, %s!", container.GetInstance(nil, Name{})))
}
