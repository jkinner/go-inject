package goose

import (
	"fmt";
	"reflect";
)

type Provider func () interface{}
type Scope interface{}

type Injector interface {
	CreateChildInjector() Injector
	GetInstanceForKey(key Key) interface{}
	GetInstance(instanceType reflect.Type) interface{}
	GetTaggedInstance(instanceType reflect.Type, tag Tag) interface{}
	Bind(key Key, provider Provider)
	BindToInstance(key Key, instance interface{})
	getBinding(key Key) (Provider, bool)
	findParentBinding(key Key) (Provider, bool)
	findChildBinding(key Key) (Provider, bool)
}

type injector struct {
	bindings map[Key] Provider
	parent Injector
	children map[*injector] *injector
}

func CreateInjector() Injector {
	return &injector {
		bindings: make(map[Key] Provider),
		parent: nil,
		children: make(map[*injector] *injector),
	}
}

func (this injector) GetInstanceForKey(key Key) interface{} {
	var provider, ok = this.bindings[key]
	if ok {
		return provider()
	}

	provider, ok = this.findParentBinding(key)
	if ok {
		return provider()
	}

	panic(fmt.Sprintf("Unable to find %v in injector", key))
}

func (this injector) GetInstance(instanceType reflect.Type) interface{} {
	return this.GetInstanceForKey(CreateKeyForType(instanceType))
}

func (this injector) GetTaggedInstance(instanceType reflect.Type, tag Tag) interface{} {
	return this.GetInstanceForKey(CreateKeyForTaggedType(instanceType, tag))
}

func (this injector) Bind(key Key, provider Provider) {
	_, ok := this.bindings[key]
	if ! ok {
		_, ok = this.findParentBinding(key)
	}
	if ! ok {
		_, ok = this.findChildBinding(key)
	}
	if ok {
		panic(fmt.Sprintf("%v is already bound", key));
	}
	this.bindings[key] = provider
}

func (this injector) BindToInstance(key Key, instance interface{}) {
	this.Bind(key, func () interface{} { return instance })
}

func (this *injector) CreateChildInjector() Injector {
	child := injector {
		bindings: make(map[Key] Provider),
		parent: this,
		children: make(map[*injector] *injector),
	}
	this.children[&child] = &child
	return &child
}

func (this injector) getBinding(key Key) (provider Provider, ok bool) {
	provider, ok = this.bindings[key]
	return
}

func (this injector) findParentBinding(key Key) (Provider, bool) {
	if this.parent != nil {
		if provider, ok := this.parent.getBinding(key); ok {
			return provider, ok
		}
		return this.parent.findParentBinding(key)
	}

	return nil, false
}

func (this injector) findChildBinding(key Key) (Provider, bool) {
	// Bindings MUST NOT be present in any child injectors
	for child := range(this.children) {
		if provider, ok := child.getBinding(key); ok {
			return provider, ok
		}
		return child.findChildBinding(key)
	}

	return nil, false
}
