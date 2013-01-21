package goose

import (
	"fmt";
	"reflect";
)

// Key used to uniquely identify a binding.
type Key interface {}

// Type used to identify a tagged type binding.
type Tag interface {}

/*
	Signature for provider functions. Provider functions are used to dynamically allocate an instance
	at run-time.
*/
type Provider func (container Container) interface{}

/*
	Injector aggregates binding configuration and creates Containers based on that configuration.
	Binding configuration is defined by a Key that consists of a type and an optional tag.

	A child injector may be used to create bindings that are intended to be used only by part
	of the system. When a Container is created from a child injector, the bindings of the
	parent Injector are also available.

	Keys may only be bound once across an Injector and all of its descendant Injectors (children
	and their children). The Injector will panic if an attempt is made to rebind an already-bound
	Key.

	In order to look up a bound type, use the CreateContainer() method and call the appropriate
	methods on the returned Container.
*/
type Injector interface {
	// Binds a type to a Provider function.
	Bind(bindingType reflect.Type, provider Provider)

	// Binds a type to a single instance.
	BindInstance(bindingType reflect.Type, instance interface{})

	// Binds a key to a Provider function.
	BindKey(key Key, provider Provider)

	// Binds a key to a single instance.
	BindKeyInstance(key Key, instance interface{})

	// Binds a tagged type to a Provider function.
	BindTagged(bindingType reflect.Type, tag Tag, provider Provider)

	// Binds a tagged type to a single instance.
	BindTaggedInstance(bindingType reflect.Type, tag Tag, instance interface{})

	// Creates a child injector that can bind additional types not available from this Injector.
	CreateChildInjector() Injector

	// Creates a Container that can be used to retrieve instance objects from the Injector.
	CreateContainer() Container

	// Exposes a type to its parent injector.
	Expose(bindingType reflect.Type)

	// Exposes a key binding to its parent injector.
	ExposeKey(key Key)

	// Exposes a tagged type to its parent injector.
	ExposeTagged(bindingType reflect.Type, tag Tag)

	// Gets the binding for a key, searching the current injector and all ancestor injectors.
	getBinding(key Key) (Provider, bool)

	/*
		Searches the parent injector for the key, continuing to search upward until the
		root injector is found.
	*/
	findAncestorBinding(key Key) (Provider, bool)
}

// The context holds all the keys used by a given object.
type context map[Key] Key

type injector struct {
	// The bindings present in this injector.
	bindings map[Key] Provider

	// The parent injector. See getBinding(), findAncestorBinding().
	parent *injector

	// The child injectors. Each child will have this injector as the parent.
	children map[*injector] *injector

	// A pointer to the context for this injector and all ancestor and descendant injectors.
	context *context
}

func CreateInjector() Injector {
	context := make(context)
	return &injector {
		bindings: make(map[Key] Provider),
		parent: nil,
		children: make(map[*injector] *injector),
		context: &context,
	}
}

func (this injector) Bind(instanceType reflect.Type, provider Provider) {
	this.BindKey(CreateKeyForType(instanceType), provider)
}

func (this injector) BindInstance(instanceType reflect.Type, instance interface{}) {
	this.BindKeyInstance(CreateKeyForType(instanceType), instance)
}

func (this injector) BindKey(key Key, provider Provider) {
	context := *this.context
	if _, ok := context[key]; ok {
		panic(fmt.Sprintf("%v is already bound", key));
	}
	context[key] = key
	this.bindings[key] = provider
}

func (this injector) BindKeyInstance(key Key, instance interface{}) {
	this.BindKey(key, func (container Container) interface{} { return instance })
}

func (this injector) BindTaggedInstance(instanceType reflect.Type, tag Tag,
		instance interface{}) {
	this.BindKeyInstance(CreateKeyForTaggedType(instanceType, tag), instance)
}

func (this injector) BindTagged(instanceType reflect.Type, tag Tag, provider Provider) {
	this.BindKey(CreateKeyForTaggedType(instanceType, tag), provider)
}

// Creates a child injector that can contain bindings not available to the parent injector.
func (this *injector) CreateChildInjector() Injector {
	child := injector {
		bindings: make(map[Key] Provider),
		parent: this,
		children: make(map[*injector] *injector),
		context: this.context,
	}
	this.children[&child] = &child
	return &child
}

// Creates a Container that is used to request values during object creation.
func (this injector) CreateContainer() Container {
	return container {
		this,
		make(context),
	}
}

func (this injector) Expose(bindingType reflect.Type) {
	this.ExposeKey(CreateKeyForType(bindingType))
}

func (this injector) ExposeTagged(bindingType reflect.Type, tag Tag) {
	this.ExposeKey(CreateKeyForTaggedType(bindingType, tag))
}

func (this injector) ExposeKey(key Key) {
	if this.parent == nil {
		panic(fmt.Sprintf("No parent injector available when exposing %s", key))
	}
	if _, exists := this.bindings[key]; !exists {
		panic(fmt.Sprintf("No binding for %s is present in this injector", key))
	}

	// TODO(jkinner): Worry about caching in scopes.
	this.parent.bindings[key] = this.bindings[key]
}

func (this injector) getBinding(key Key) (provider Provider, ok bool) {
	provider, ok = this.bindings[key]
	return
}

func (this injector) findAncestorBinding(key Key) (Provider, bool) {
	if this.parent != nil {
		if provider, ok := this.parent.getBinding(key); ok {
			return provider, ok
		}
		return this.parent.findAncestorBinding(key)
	}

	return nil, false
}

/*
	Container provides access to the bindings configured in an Injector. All bindings are available
	as a Provider or as a value. A new Container should be used for each injected type. A Container
	will panic if a key is looked up more than once. This behavior is intended to detect and prevent
	cycles in depedencies.
	
	For example, suppose you have a type A that gets an instance of type B that in turn relies
	on type A again (A -> B -> A). The types would have a structure like this:

	type A struct {
		B
	}

	type B struct {
		A
	}

	func ConfigureInjector(injector goose.Injector) {
		injector.Bind(reflect.TypeOf(A(nil)), func (container Container) interface{} {
			return A { createB(container) }
		}
		injector.Bind(reflect.TypeOf(B(nil)), func (container Container) interface{} {
			return B { createA(container) }
		}
	}

	func createA(container goose.Container) {
		return A { B: container.GetInstanceForKey(reflect.TypeOf(A(nil))) }
	}

	func createB(container goos.Container) {
		return B { A: container.GetInstanceForKey(reflect.TypeOf(B(nil)) }
	}
*/
type Container interface {
	// Returns an instance of the type bound by the key.
	GetInstanceForKey(key Key) interface{}

	// Returns an instance of the type.
	GetInstance(instanceType reflect.Type) interface{}

	// Returns an instance of the type tagged with the tag.
	GetTaggedInstance(instanceType reflect.Type, tag Tag) interface{}

	// Returns a Provider that can return an instance of the type bound by the key.
	GetProviderForKey(key Key) Provider

	// Returns a Provider that can return an instance of the type.
	GetProvider(instanceType reflect.Type) Provider

	// Returns a Provider that can return an instance of the type tagged with the tag.
	GetTaggedProvider(instanceType reflect.Type, tag Tag) Provider
}

type container struct {
	// The injector holding the bindings available to the container.
	injector

	// The invocation context, holding all the previous requests to prevent duplicate requests.
	context
}

// Returns a Provider that can create an instance of the type bound to the key.
func (this container) GetProviderForKey(key Key) Provider {
	if _, exists := this.context[key]; exists {
		panic(fmt.Sprintf("Already looked up %v. Is there a cycle of dependencies?", key))
	}

	this.context[key] = key

	var provider, ok = this.bindings[key]
	if ok {
		return provider
	}

	provider, ok = this.findAncestorBinding(key)
	if ok {
		return provider
	}

	panic(fmt.Sprintf("Unable to find %v in injector", key))
}

// Returns a Provider that can create an instance of the instanceType.
func (this container) GetProvider(instanceType reflect.Type) Provider {
	return this.GetProviderForKey(CreateKeyForType(instanceType))
}

// Returns a Provider that can create an instance of the instanceType tagged with tag.
func (this container) GetTaggedProvider(instanceType reflect.Type, tag Tag) Provider {
	return this.GetProviderForKey(CreateKeyForTaggedType(instanceType, tag))
}

// Returns an instance of the type bound to the key.
func (this container) GetInstanceForKey(key Key) interface{} {
	return this.GetProviderForKey(key)(this)
}

// Returns an instance of the instanceType.
func (this container) GetInstance(instanceType reflect.Type) interface{} {
	return this.GetInstanceForKey(CreateKeyForType(instanceType))
}

// Returns an instance of the instanceType tagged with tag.
func (this container) GetTaggedInstance(instanceType reflect.Type, tag Tag) interface{} {
	return this.GetInstanceForKey(CreateKeyForTaggedType(instanceType, tag))
}

type key struct {
	typeLiteral reflect.Type
	tag Tag
}

func CreateKeyForType(typeLiteral reflect.Type) Key {
	return key {
		typeLiteral,
		nil,
	}
}

func CreateKeyForTaggedType(typeLiteral reflect.Type, tag Tag) Key {
	return key {
		typeLiteral,
		tag,
	}
}

func (this key) String() string {
	if this.tag == nil {
		return fmt.Sprintf("%v", this.typeLiteral)
	}
	return fmt.Sprintf("%v<%s(%v)>", this.typeLiteral, reflect.TypeOf(this.tag).Name(), this.tag)
}
