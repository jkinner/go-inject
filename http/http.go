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

package inject_http

import (
	"code.google.com/p/go-inject"
	"code.google.com/p/go-inject/multi"
	"flag"
	"fmt"
	"log"
	"net/http"
)

// Value for the http_port flag (default: 80)
var httpPort *int = flag.Int("http_port", 80, "port on which to start the http listener")

// Container key for the HTTP server itself
type Server struct{}

// Container key for looking up handlers.
type Handlers struct{}

// Tag on the Handlers key that specifies that the handlers are functions.
type Func struct{}

// Container key for the HTTP port for the server.
type Port struct{}

// Scope tag for the caching in the context of the *http.Request
type RequestScoped struct{}

type HandlerMap map[interface{}]interface{}
type HandlerFuncMap map[string]func(http.ResponseWriter, *http.Request)

var requestScope = inject.CreateSimpleScopeWithName("HTTP Request")

type scopingHandler struct {
	handler http.Handler
}

func (this scopingHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer requestScope.Exit(request)
	requestScope.Enter(request)
	this.handler.ServeHTTP(writer, request)
}

func providesHttpServer(_ inject.Context, container inject.Container) interface{} {
	port := container.GetInstance(nil, Port{}).(int)

	serveMux := http.NewServeMux()

	handlers := container.GetTaggedInstance(
		nil,
		Handlers{},
		inject_multi.Values{},
	).(map[interface{}]interface{})
	for path, handler := range handlers {
		serveMux.Handle(path.(string), scopingHandler{handler.(http.Handler)})
	}

	handlerFuncs := container.GetTaggedInstance(
		nil,
		inject.TaggedKey{Handlers{}, Func{}},
		inject_multi.Values{},
	).(map[interface{}]interface{})
	for path, handlerFunc := range handlerFuncs {
		serveMux.HandleFunc(path.(string),
			func(w http.ResponseWriter, request *http.Request) {
				defer requestScope.Exit(request)
				requestScope.Enter(request)
				handlerFunc.(func(http.ResponseWriter, *http.Request))(w, request)
			})
	}

	log.Printf("Creating HTTP server listening on port %d", port)
	return http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: serveMux,
	}
}

// Binds the following:
//   inject_http.Port - the value of the http_port flag
func ConfigureFlags(injector inject.Injector) {
	injector.Bind(Port{}, func(_ inject.Context, _ inject.Container) interface{} { return *httpPort })
}

// Binds the following:
//   RequestScoped scope
func ConfigureScopes(injector inject.Injector) {
	injector.BindScope(requestScope, RequestScoped{})
}

// Binds the following:
//   inject_http.Server{} - the HTTP server itself
// Requires these bindings:
//   inject_http.Handlers - a HandlerMap assigning a path to a http.Handler
//   inject_http.Handlers<inject_http.Func> - a HanderFuncMap assigning a path to a handler func
func ConfigureInjector(injector inject.Injector) {
	injector.Bind(Server{}, providesHttpServer)
	inject_multi.EnsureMapBound(injector, Handlers{})
	inject_multi.EnsureMapBound(injector, inject.TaggedKey{Handlers{}, Func{}})
}

func BindHandler(injector inject.Injector, pattern string, handler http.Handler) {
	inject_multi.BindMapInstance(
		injector,
		Handlers{},
		pattern,
		handler,
	)
}

func BindHandlerFunc(injector inject.Injector, pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
	inject_multi.BindMapTaggedInstance(
		injector,
		Handlers{},
		Func{},
		pattern,
		handlerFunc,
	)
}
