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

package inject_multi

import (
	"code.google.com/p/go-inject"
)

type mapsKey struct {
	injector inject.Injector
	key inject.Key
}

// Map that holds the maps for each bound key.
var maps map[mapsKey]map[interface{}]inject.Provider = make(map[mapsKey]map[interface{}]inject.Provider)

type Values struct{}

func createOrGetMap(injector inject.Injector, key inject.Key) map[interface{}]inject.Provider {
	mapKey := mapsKey { injector, key }
	if theMap, ok := maps[mapKey]; ok {
		return theMap
	}

	theMap := make(map[interface{}]inject.Provider)
	maps[mapKey] = theMap
	// Inject the map itself to get a map to providers.
	injector.BindInstance(key, theMap)
	// Inject the map tagged with Values{} to get a map of values.
	valueKey := inject.TaggedKey { key, Values{} }
	injector.Bind(valueKey,
		func (context inject.Context, container inject.Container) interface{} {
			valueMap := make(map[interface{}]interface{})
			providerMap := container.GetInstance(context, key).(map[interface{}]inject.Provider)
			for mapKey, provider := range(providerMap) {
				valueMap[mapKey] = provider(context, container)
			}
			return valueMap
		})
	return theMap
}

func EnsureMapBound(injector inject.Injector, key inject.Key) {
	createOrGetMap(injector, key)
}

// Binds a type to a inject.Provider function.
func BindMap(injector inject.Injector, key inject.Key, mapKey interface{}, provider inject.Provider) {
	createOrGetMap(injector, key)[mapKey] = provider
}

// Binds a type to a single instance.
func BindMapInstance(injector inject.Injector, key inject.Key, mapKey interface{}, instance interface{}) {
	BindMap(injector, key, mapKey, func (_ inject.Context, _ inject.Container) interface{} { return instance })
}

// Binds a key to a inject.Provider function, caching it within the specified scope.
func BindMapInScope(injector inject.Injector, key inject.Key, mapKey interface{}, provider inject.Provider, scopeTag inject.Tag) {
	BindMap(injector, key, mapKey, injector.Scope(key, provider, scopeTag))
}

// Binds a key to a single instance.
func BindMapInstanceInScope(injector inject.Injector, key inject.Key, mapKey interface{}, instance interface{}, scopeTag inject.Tag) {
	BindMapInScope(
		injector,
		key,
		mapKey,
		func (_ inject.Context, _ inject.Container) interface{} { return instance },
		scopeTag,
	)
}

// Binds a tagged type to a inject.Provider function.
func BindMapTagged(injector inject.Injector, key inject.Key, tag inject.Tag, mapKey interface{}, provider inject.Provider) {
	BindMap(injector, inject.TaggedKey { key, tag }, mapKey, provider)
}

// Binds a tagged type to a inject.Provider function.
func BindMapTaggedInScope(injector inject.Injector, key inject.Key, tag inject.Tag, mapKey interface{}, provider inject.Provider, scopeTag inject.Tag) {
	BindMapInScope(injector, inject.TaggedKey { key, tag }, mapKey, provider, scopeTag)
}

// Binds a tagged type to a single instance.
func BindMapTaggedInstance(injector inject.Injector, key inject.Key, tag inject.Tag, mapKey interface{}, instance interface{}) {
	BindMap(injector, inject.TaggedKey { key, tag }, mapKey, func (_ inject.Context, _ inject.Container) interface{} { return instance })
}

// Binds a tagged type to a single instance.
func BindMapTaggedInstanceInScope(injector inject.Injector, key inject.Key, tag inject.Tag, mapKey interface{}, instance interface{}, scopeTag inject.Tag) {
	BindMapInScope(
		injector,
		inject.TaggedKey { key, tag },
		mapKey,
		func (_ inject.Context, _ inject.Container) interface{} { return instance },
		scopeTag,
	)
}
