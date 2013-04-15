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
	"code.google.com/p/go-inject"
	"os"
	"reflect"
)

/*
	This is sort of how Guice works. But it's not really go-ish.
*/
func main() {
	injector := inject.CreateInjector()
	var name = "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	injector.BindInstance(reflect.TypeOf(""), name)
	sayHello(injector.CreateContainer())
}

func sayHello(container inject.Container) {
	fmt.Println(fmt.Sprintf("Hello, %s!", container.GetInstance(nil, reflect.TypeOf(""))))
}
