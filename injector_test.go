package goose

import (
//	"fmt";
	"log";
	"reflect";
	"testing";
)

const (
	THING1 = iota
	THING2
)

func TestTypeInstanceBinding(t *testing.T) {
	log.Println("Testing instance binding")
	injector := CreateInjector()
	injector.BindToInstance(CreateKeyForType(reflect.TypeOf("")), "foo")
	value := injector.GetInstance(reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'");
	}
}

func TestTypeInstanceBindingThroughMethod(t *testing.T) {
	log.Println("Testing instance binding")
	injector := CreateInjector()
	injector.BindToInstance(CreateKeyForType(reflect.TypeOf("")), "foo")
	value := getStringInstance(injector)
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'");
	}
}

func getStringInstance(injector Injector) interface{} {
	return injector.GetInstance(reflect.TypeOf(""))
}

func TestTypeInstanceBindingFailure(t *testing.T) {
	log.Println("Testing instance binding (failure)")
	var injector = CreateInjector()
	defer func() {
		recover()
	}()

	injector.GetInstance(reflect.TypeOf(""))
	t.Error("Expected to panic on missing string key")
}

func TestTypeInstanceBindingWithTag(t *testing.T) {
	log.Println("Testing instance binding with tag")
	injector := CreateInjector()
	injector.BindToInstance(CreateKeyForTaggedType(
		reflect.TypeOf(""), THING1), "foo")
	value := injector.GetTaggedInstance(reflect.TypeOf(""), THING1)
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'");
	}
}

func TestTypeInstanceBindingWithTagFailure(t *testing.T) {
	log.Println("Testing instance binding with tag (failure)")
	injector := CreateInjector()
	defer func() {
		recover()
	}()

	injector.GetTaggedInstance(reflect.TypeOf(""), THING1)
	t.Error("Expected to panic on missing string(THING1) key")
}

func TestTypeProviderBinding(t *testing.T) {
	log.Println("Testing provider binding")
	injector := CreateInjector()
	injector.Bind(CreateKeyForType(reflect.TypeOf("")), func () interface{} { return "foo" })
	value := injector.GetInstance(reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'");
	}
}

func TestTypeProviderBindingWithTag(t *testing.T) {
	log.Println("Testing provider binding with tag")
	injector := CreateInjector()
	injector.Bind(CreateKeyForTaggedType(reflect.TypeOf(""), THING2),
			func () interface{} { return "foo" })
	value := injector.GetTaggedInstance(reflect.TypeOf(""), THING2)
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'");
	}
}

func TestDeferToParentInjector(t *testing.T) {
	log.Println("Testing provider binding with parent injector")
	parent := CreateInjector()
	child := parent.CreateChildInjector()

	parent.BindToInstance(CreateKeyForType(reflect.TypeOf("")), "foo")
	value := child.GetInstance(reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get parent binding for string value 'foo'")
	}
}

func TestAlreadyBoundInChildInjector(t *testing.T) {
	defer func () {
		recover()
	}()

	parent := CreateInjector()
	child := parent.CreateChildInjector()

	child.BindToInstance(CreateKeyForType(reflect.TypeOf("")), "foo")
	parent.BindToInstance(CreateKeyForType(reflect.TypeOf("")), "foo")
	t.Error("Expected to fail because already bound in child injector")
}

func TestAlreadyBoundInParentInjector(t *testing.T) {
	defer func () {
		recover()
	}()

	parent := CreateInjector()
	child := parent.CreateChildInjector()

	parent.BindToInstance(CreateKeyForType(reflect.TypeOf("")), "foo")
	child.BindToInstance(CreateKeyForType(reflect.TypeOf("")), "foo")
	t.Error("Expected to fail because already bound in parent injector")
}
