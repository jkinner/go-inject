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
	"github.com/jkinner/goose"
	"github.com/jkinner/goose/http"
	"log"
	"net/http"
)

type Counter struct{}

func sayFoo(w http.ResponseWriter, request *http.Request) {
	w.Header().Add(
		"Content-Type",
		"text/plain",
	)
	w.Write([]byte(fmt.Sprintf("Foo\n")))
}

func ConfigureInjector(injector goose.Injector) {
	i := 0
	injector.BindInScope(Counter{}, func(_ goose.Context, _ goose.Container) interface{} {
		i += 1
		return i
	}, goose_http.RequestScoped{})

	goose_http.BindHandlerFunc(injector, "/",
		func (w http.ResponseWriter, request *http.Request) {
			w.Header().Add(
				"Content-Type",
				"text/plain",
			)
			w.Write([]byte(fmt.Sprintf("Hello! (%d)\n",
				injector.CreateContainer().GetInstance(request, Counter{}).(int))))
			w.Write([]byte(fmt.Sprintf("Hello! (%d)\n",
				injector.CreateContainer().GetInstance(request, Counter{}).(int))))
		})
}

func ConfigureFooInjector(injector goose.Injector) {
	goose_http.BindHandlerFunc(injector, "/foo/", sayFoo)
}

func main() {
	// Initialize the flag(s)
	flag.Parse()

	// Create the injector to begin configuring the bindings.
	injector := goose.CreateInjector()

	// Bind any flags used in the goose_http module. Flag bindings are in a separate function
	// so that the module can be used without flags, with explicit bindings (see multi_http.go).
	goose_http.ConfigureFlags(injector)

	// Bind any scope tags used in the goose_http module. Scope bindings are in a separate function
	// so the module can be used more than once. The scope binding can only ever happen once.
	goose_http.ConfigureScopes(injector)

	// Now configure the bindings in the goose_http module. This step can be performed more than
	// once in multiple child injectors.
	goose_http.ConfigureInjector(injector)

	// There are multiple ConfigureInjector functions to demonstrate that multiple modules
	// can configure HTTP mappings. If two modules try to bind the same path, the system will panic.
	ConfigureFooInjector(injector)
	ConfigureInjector(injector)

	// Look up an object by creating a Container. The Container ensures there are no dependency
	// loops configured in the Injector. Top-level code should create a new container for each
	// object lookup. Provider functions will have a container passed to them, and they use that
	// container for their lookups.
	container := injector.CreateContainer()
	httpServer := container.GetInstance(nil, goose_http.Server{}).(http.Server)

	// Start the HTTP server.
	log.Fatal(httpServer.ListenAndServe())
}
