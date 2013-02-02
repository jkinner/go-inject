package goose_multi

import (
	"github.com/jkinner/goose"
)

type mapsKey struct {
	injector goose.Injector
	key goose.Key
}

// Map that holds the maps for each bound key.
var maps map[mapsKey]map[interface{}]goose.Provider = make(map[mapsKey]map[interface{}]goose.Provider)

type Values struct{}

func createOrGetMap(injector goose.Injector, key goose.Key) map[interface{}]goose.Provider {
	mapKey := mapsKey { injector, key }
	if theMap, ok := maps[mapKey]; ok {
		return theMap
	}

	theMap := make(map[interface{}]goose.Provider)
	maps[mapKey] = theMap
	// Inject the map itself to get a map to providers.
	injector.BindInstance(key, theMap)
	// Inject the map tagged with Values{} to get a map of values.
	valueKey := goose.TaggedKey { key, Values{} }
	injector.Bind(valueKey,
		func (context goose.Context, container goose.Container) interface{} {
			valueMap := make(map[interface{}]interface{})
			providerMap := container.GetInstance(context, key).(map[interface{}]goose.Provider)
			for mapKey, provider := range(providerMap) {
				valueMap[mapKey] = provider(context, container)
			}
			return valueMap
		})
	return theMap
}

func EnsureMapBound(injector goose.Injector, key goose.Key) {
	createOrGetMap(injector, key)
}

// Binds a type to a goose.Provider function.
func BindMap(injector goose.Injector, key goose.Key, mapKey interface{}, provider goose.Provider) {
	createOrGetMap(injector, key)[mapKey] = provider
}

// Binds a type to a single instance.
func BindMapInstance(injector goose.Injector, key goose.Key, mapKey interface{}, instance interface{}) {
	BindMap(injector, key, mapKey, func (_ goose.Context, _ goose.Container) interface{} { return instance })
}

// Binds a key to a goose.Provider function, caching it within the specified scope.
func BindMapInScope(injector goose.Injector, key goose.Key, mapKey interface{}, provider goose.Provider, scopeTag goose.Tag) {
	BindMap(injector, key, mapKey, injector.Scope(key, provider, scopeTag))
}

// Binds a key to a single instance.
func BindMapInstanceInScope(injector goose.Injector, key goose.Key, mapKey interface{}, instance interface{}, scopeTag goose.Tag) {
	BindMapInScope(
		injector,
		key,
		mapKey,
		func (_ goose.Context, _ goose.Container) interface{} { return instance },
		scopeTag,
	)
}

// Binds a tagged type to a goose.Provider function.
func BindMapTagged(injector goose.Injector, key goose.Key, tag goose.Tag, mapKey interface{}, provider goose.Provider) {
	BindMap(injector, goose.TaggedKey { key, tag }, mapKey, provider)
}

// Binds a tagged type to a goose.Provider function.
func BindMapTaggedInScope(injector goose.Injector, key goose.Key, tag goose.Tag, mapKey interface{}, provider goose.Provider, scopeTag goose.Tag) {
	BindMapInScope(injector, goose.TaggedKey { key, tag }, mapKey, provider, scopeTag)
}

// Binds a tagged type to a single instance.
func BindMapTaggedInstance(injector goose.Injector, key goose.Key, tag goose.Tag, mapKey interface{}, instance interface{}) {
	BindMap(injector, goose.TaggedKey { key, tag }, mapKey, func (_ goose.Context, _ goose.Container) interface{} { return instance })
}

// Binds a tagged type to a single instance.
func BindMapTaggedInstanceInScope(injector goose.Injector, key goose.Key, tag goose.Tag, mapKey interface{}, instance interface{}, scopeTag goose.Tag) {
	BindMapInScope(
		injector,
		goose.TaggedKey { key, tag },
		mapKey,
		func (_ goose.Context, _ goose.Container) interface{} { return instance },
		scopeTag,
	)
}
