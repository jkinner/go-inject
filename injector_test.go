package goose

import (
	"reflect";
	"testing";
)

// Tags that are primitive types. It is recommended to use const values if you go this way.
const (
	THING1 = iota
	THING2
)

// Types that are used as tags.
type Thing1 struct {}
type Thing2 struct {}

// Type used as a tag that is parameterized.
type ParameterizedThing struct {
	string
}

type ThingTag struct {
	int
}

// Not preferred, but allowed. Similar to using @Named(constValue). Remember that const of
// primitive types are just the value of the primitive, making collisions extremely likely.
// Prefer to use a Type tag or a ParameterizedType tag that contains your constant type.
// These are demonstrated in other tests.
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

// Use a parameterized tag type along with a const to assist users in creating the right tag.
func TestParameterizedTagInstanceBinding(t *testing.T) {
	injector := CreateInjector()
	injector.BindTaggedInstance(
		reflect.TypeOf(""),
		ThingTag{THING1},
		"foo",
	)
	injector.BindTaggedInstance(
		reflect.TypeOf(""),
		ThingTag{THING2},
		"bar",
	)
}

// Use separate types as a tag on the type.
func TestTypeTaggedInstanceBinding(t *testing.T) {
	injector := CreateInjector()
	// To use a type as a tag, create an empty instance of that type.
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

// Use a single type with a generic primitive as a type parameter for the tag.
func TestParameterizedTypeTaggedInstanceBinding(t *testing.T) {
	injector := CreateInjector()
	// To use a type as a tag, create an empty instance of that type.
	injector.BindTaggedInstance(
		reflect.TypeOf(""),
		ParameterizedThing{"foo"},
		"foo",
	)
	injector.BindTaggedInstance(
		reflect.TypeOf(""),
		ParameterizedThing{"bar"},
		"bar",
	)
}

// Attempt (and fail) to bind a parameterized type tag twice in the same injector.
func TestDuplicateParameterizedTypeTaggedInstanceBindingPanics(t *testing.T) {
	injector := CreateInjector()
	// To use a type as a tag, create an empty instance of that type.
	injector.BindTaggedInstance(
		reflect.TypeOf(""),
		ParameterizedThing{"foo"},
		"foo",
	)
	defer func() {
		// Expected
		recover()
	}()
	injector.BindTaggedInstance(
		reflect.TypeOf(""),
		ParameterizedThing{"foo"},
		"bar",
	)
	t.Error("Expected a panic because ParameterizedThing<foo> is already been bound")
}

// Attempt (and fail) to bind a type tag twice in the same injector.
func TestDuplicateTaggedTypeInstanceBindingPanics(t *testing.T) {
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

	t.Error("Expected a panic when binding Thing2 again")
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

func TestContainerRepeatedLookupsDisallowed(t *testing.T) {
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
	container := injector.CreateContainer()

	defer func() {
		recover()
	}()
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
	container := injector.CreateContainer()

	defer func() {
		recover()
	}()
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

func TestDelegateToParentInjector(t *testing.T) {
	parent := CreateInjector()
	child := parent.CreateChildInjector()

	parent.BindInstance(reflect.TypeOf(""), "foo")
	value := child.CreateContainer().GetInstance(reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get parent binding for string value 'foo'")
	}
}

func TestExposeToParent(t *testing.T) {
	parent := CreateInjector()
	child := parent.CreateChildInjector()

	child.BindInstance(reflect.TypeOf(""), "foo")
	child.Expose(reflect.TypeOf(""))
	value := parent.CreateContainer().GetInstance(reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get parent binding for string value 'foo'")
	}
}

func TestExposeToParentDoesNotExposeFurther(t *testing.T) {
	parent := CreateInjector()
	child := parent.CreateChildInjector()
	grandchild := child.CreateChildInjector()

	grandchild.BindInstance(reflect.TypeOf(""), "foo")
	grandchild.Expose(reflect.TypeOf(""))
	value := child.CreateContainer().GetInstance(reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get parent binding for string value 'foo'")
	}
	defer func() {
		recover() // Expected
	}()
	parent.CreateContainer().GetInstance(reflect.TypeOf(""))
	t.Error("Expected a panic when trying to get type exposed by grandchild from grandparent")
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
