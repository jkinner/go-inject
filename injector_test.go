package goose

import (
	"reflect";
	"testing";
)

const (
	THING1 = iota
	THING2
)

type Thing1 struct {}
type Thing2 struct {}

func TestTaggedInstanceBinding(t *testing.T) {
	injector := CreateInjector()
	injector.BindTaggedInstance(
		reflect.TypeOf(""),
		THING1,
		"foo",
	)
	injector.BindTaggedInstance(
		reflect.TypeOf(""),
		THING2,
		"bar",
	)
}

func TestTypeTaggedInstanceBinding(t *testing.T) {
	injector := CreateInjector()
	injector.BindTaggedInstance(
		reflect.TypeOf(""),
		Thing1{},
		"foo",
	)
	injector.BindTaggedInstance(
		reflect.TypeOf(""),
		Thing2{},
		"bar",
	)
}

func TestNoDuplicateTypeTaggedInstanceBinding(t *testing.T) {
	injector := CreateInjector()
	injector.BindTaggedInstance(
		reflect.TypeOf(""),
		Thing1{},
		"foo",
	)
	injector.BindTaggedInstance(
		reflect.TypeOf(""),
		Thing2{},
		"bar",
	)
	defer func() {
		recover()
	}()

	injector.BindTaggedInstance(
		reflect.TypeOf(""),
		Thing2{},
		"bar",
	)

	t.Error("Expected a panic when binding Thing2{} again")
}

func TestTypeInstanceBinding(t *testing.T) {
	injector := CreateInjector()
	injector.BindInstance(reflect.TypeOf(""), "foo")
	container := injector.CreateContainer()
	value := container.GetInstance(reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'");
	}
}

func TestNoLookupCycles(t *testing.T) {
	injector := CreateInjector()
	injector.BindInstance(reflect.TypeOf(""), "foo")
	container := injector.CreateContainer()
	value := container.GetInstance(reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'");
	}
	defer func () {
		recover()
	}()
	container.GetInstance(reflect.TypeOf(""))
	t.Error("Expected to fail when retrieving the same key from the container")
}

func TestTypeInstanceBindingThroughMethod(t *testing.T) {
	injector := CreateInjector()
	injector.BindInstance(reflect.TypeOf(""), "foo")
	value := getStringInstance(injector)
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'");
	}
}

func getStringInstance(injector Injector) interface{} {
	return injector.CreateContainer().GetInstance(reflect.TypeOf(""))
}

func TestTypeInstanceBindingFailure(t *testing.T) {
	var injector = CreateInjector()
	defer func() {
		recover()
	}()

	container := injector.CreateContainer()
	container.GetInstance(reflect.TypeOf(""))
	t.Error("Expected to panic on missing string key")
}

func TestTypeInstanceBindingWithTag(t *testing.T) {
	injector := CreateInjector()
	injector.BindTaggedInstance(reflect.TypeOf(""), THING1, "foo")
	container := injector.CreateContainer()
	value := container.GetTaggedInstance(reflect.TypeOf(""), THING1)
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'");
	}
}

func TestTypeInstanceBindingWithTagFailure(t *testing.T) {
	injector := CreateInjector()
	defer func() {
		recover()
	}()

	container := injector.CreateContainer()
	container.GetTaggedInstance(reflect.TypeOf(""), THING1)
	t.Error("Expected to panic on missing string(THING1) key")
}

func TestTypeProviderBinding(t *testing.T) {
	injector := CreateInjector()
	injector.Bind(reflect.TypeOf(""), func (container Container) interface{} { return "foo" })
	container := injector.CreateContainer()
	value := container.GetInstance(reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'");
	}
}

func TestTypeProviderBindingWithTag(t *testing.T) {
	injector := CreateInjector()
	injector.BindTagged(reflect.TypeOf(""), THING2,
			func (container Container) interface{} { return "foo" })
	container := injector.CreateContainer()
	value := container.GetTaggedInstance(reflect.TypeOf(""), THING2)
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'");
	}
}

func TestDeferToParentInjector(t *testing.T) {
	parent := CreateInjector()
	child := parent.CreateChildInjector()

	parent.BindInstance(reflect.TypeOf(""), "foo")
	value := child.CreateContainer().GetInstance(reflect.TypeOf(""))
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

	child.BindInstance(reflect.TypeOf(""), "foo")
	parent.BindInstance(reflect.TypeOf(""), "foo")
	t.Error("Expected to fail because already bound in child injector")
}

func TestAlreadyBoundInParentInjector(t *testing.T) {
	defer func () {
		recover()
	}()

	parent := CreateInjector()
	child := parent.CreateChildInjector()

	parent.BindInstance(reflect.TypeOf(""), "foo")
	child.BindInstance(reflect.TypeOf(""), "foo")
	t.Error("Expected to fail because already bound in parent injector")
}
