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
	"code.google.com/p/go-inject"
	"code.google.com/p/go-inject/http"
	"flag"
	"fmt"
	"log"
	"net/http"
)

type Counter struct{}

// Configure a counter provider that increments on each invocation.
func ConfigureCounter(injector inject.Injector) {
	i := 0
	injector.BindInScope(Counter{}, func(_ inject.Context, _ inject.Container) interface{} {
		i += 1
		return i
	}, inject_http.RequestScoped{})
}

// Configure HTTP server path mappings and provider functions.
func ConfigureInjector(injector inject.Injector) {

	inject_http.BindHandlerFunc(
		injector,
		"/",
		func(w http.ResponseWriter, request *http.Request) {
			w.Header().Add(
				"Content-Type",
				"text/plain",
			)
			w.Write([]byte(fmt.Sprintf("Hello, %d!\n", injector.CreateContainer().GetInstance(request, Counter{}).(int))))
			w.Write([]byte(fmt.Sprintf("Hello, %d!", injector.CreateContainer().GetInstance(request, Counter{}).(int))))
		})
}

type OneServer struct{}
type TwoServer struct{}

// Configure a single HTTP server on a port and expose it using the tag.
func ConfigureServer(injector inject.Injector, port int, tag inject.Tag) {
	// Each HTTP server needs to have a port bound. The port is statically configured here, but
	// it could be configured using a Provider that allocates an unused port.
	injector.BindInstance(inject_http.Port{}, port)
	inject_http.ConfigureInjector(injector)
	ConfigureInjector(injector)
	injector.ExposeAndTag(inject_http.Server{}, tag)
}

func main() {
	// Initialize the flag(s)
	flag.Parse()

	// Create the injector to begin configuring the bindings.
	injector := inject.CreateInjector()

	// Don't configure flags, since the HTTP module assumes a single HTTP port in use.

	// Bind any scope tags used in the inject_http module. Scope bindings are in a separate function
	// so the module can be used more than once. The scope binding can only ever happen once.
	inject_http.ConfigureScopes(injector)

	// Use a global counter. Both HTTP servers will increment at the same time.
	ConfigureCounter(
		injector,
	)

	// Create a child injector for each server. Each child injector is isolated, so the HTTP module
	// can be bound multiple times. Only the exposed bindings have the potential to collide. Note
	// that if either server were bound in the parent injector, all the bindings would collide.
	oneServerInjector := injector.CreateChildInjector()
	twoServerInjector := injector.CreateChildInjector()

	// Configure each server on its own port, exposing the server in the parent injector with its
	// own tag.
	ConfigureServer(oneServerInjector, 8080, OneServer{})
	ConfigureServer(twoServerInjector, 8081, TwoServer{})

	oneHttpServer := injector.CreateContainer().GetTaggedInstance(
		nil,
		inject_http.Server{},
		OneServer{}).(http.Server)

	// What would I do with go-routines? Two HTTP servers at the same time.
	go func() {
		log.Fatal(oneHttpServer.ListenAndServe())
	}()

	twoHttpServer := injector.CreateContainer().GetTaggedInstance(
		nil,
		inject_http.Server{},
		TwoServer{}).(http.Server)
	log.Fatal(twoHttpServer.ListenAndServe())
}
