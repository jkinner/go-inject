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

package goose

import (
	"fmt"
	"reflect"
	"testing"
)

// Tags that are primitive types. It is recommended to use const values if you go this way.
const (
	THING1 = iota
	THING2
)

// Types that are used as tags.
type Thing1 struct{}
type Thing2 struct{}

type MyContext struct{
	string
}

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
	value := container.GetInstance(nil /* context */, reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'")
	}
}

func TestContainerRepeatedLookupsDisallowed(t *testing.T) {
	injector := CreateInjector()
	injector.BindInstance(reflect.TypeOf(""), "foo")
	container := injector.CreateContainer()
	value := container.GetInstance(nil /* context */, reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'")
	}

	defer func() {
		recover()
	}()
	container.GetInstance(nil /* context */, reflect.TypeOf(""))
	t.Error("Expected to fail when retrieving the same key from the container")
}

func TestTypeInstanceBindingThroughMethod(t *testing.T) {
	injector := CreateInjector()
	injector.BindInstance(reflect.TypeOf(""), "foo")
	value := getStringInstance(injector)
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'")
	}
}

func getStringInstance(injector Injector) interface{} {
	return injector.CreateContainer().GetInstance(nil /* context */, reflect.TypeOf(""))
}

func TestTypeInstanceBindingFailure(t *testing.T) {
	var injector = CreateInjector()
	container := injector.CreateContainer()

	defer func() {
		recover()
	}()
	container.GetInstance(nil, reflect.TypeOf(""))
	t.Error("Expected to panic on missing string key")
}

func TestTypeInstanceBindingWithTag(t *testing.T) {
	injector := CreateInjector()
	injector.BindTaggedInstance(reflect.TypeOf(""), THING1, "foo")
	container := injector.CreateContainer()
	value := container.GetTaggedInstance(nil /* context */, reflect.TypeOf(""), THING1)
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'")
	}
}

func TestTypeInstanceBindingWithTagFailure(t *testing.T) {
	injector := CreateInjector()
	container := injector.CreateContainer()

	defer func() {
		recover()
	}()
	container.GetTaggedInstance(nil /* context */, reflect.TypeOf(""), THING1)
	t.Error("Expected to panic on missing string(THING1) key")
}

func TestTypeProviderBinding(t *testing.T) {
	injector := CreateInjector()
	injector.Bind(reflect.TypeOf(""), func(_ Context, container Container) interface{} { return "foo" })
	container := injector.CreateContainer()
	value := container.GetInstance(nil /* context */, reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'")
	}
}

func TestTypeProviderBindingWithTag(t *testing.T) {
	injector := CreateInjector()
	injector.BindTagged(reflect.TypeOf(""), THING2,
		func(_ Context, _ Container) interface{} { return "foo" })
	container := injector.CreateContainer()
	value := container.GetTaggedInstance(nil /* context */, reflect.TypeOf(""), THING2)
	if value == nil || value != "foo" {
		t.Error("Expected to get the string binding value 'foo'")
	}
}

func TestDelegateToParentInjector(t *testing.T) {
	parent := CreateInjector()
	child := parent.CreateChildInjector()

	parent.BindInstance(reflect.TypeOf(""), "foo")
	value := child.CreateContainer().GetInstance(nil /* context */, reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get parent binding for string value 'foo'")
	}
}

func TestChildInjectorUsedForChildBindings(t *testing.T) {
	parent := CreateInjector()
	child := parent.CreateChildInjector()

	child.BindTaggedInstance(reflect.TypeOf(""), Thing1{}, "foo")
	child.Bind(reflect.TypeOf(""), func(context Context, container Container) interface{} {
		fmt.Printf("Calling child provider with container %+v\n", container)
		return container.GetTaggedInstance(context, reflect.TypeOf(""), Thing1{})
	})

	child.Expose(reflect.TypeOf(""))
	value := parent.CreateContainer().GetInstance(nil /* context */, reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get child binding for string value 'foo' via child tagged binding")
	}
}

func TestSameBindingMultipleChildContainers(t *testing.T) {
	parent := CreateInjector()
	child1 := parent.CreateChildInjector()
	child2 := parent.CreateChildInjector()

	child1.BindTaggedInstance(reflect.TypeOf(""), Thing1{}, "foo")
	child2.BindTaggedInstance(reflect.TypeOf(""), Thing1{}, "bar")

	child1.ExposeTagged(reflect.TypeOf(""), Thing1{})
	child2.ExposeTaggedAndRenameTagged(reflect.TypeOf(""), Thing1{}, reflect.TypeOf(""), Thing2{})

	container := parent.CreateContainer()
	value := container.GetTaggedInstance(nil /* context */, reflect.TypeOf(""), Thing1{})
	if value == nil || value != "foo" {
		t.Error("Expected to get child binding for string value 'foo' via child container")
	}
	value = container.GetTaggedInstance(nil /* context */, reflect.TypeOf(""), Thing2{})
	if value == nil || value != "bar" {
		t.Error("Expected to get child binding for string value 'bar' via child container")
	}
}

func TestExposeToParent(t *testing.T) {
	parent := CreateInjector()
	child := parent.CreateChildInjector()

	child.BindInstance(reflect.TypeOf(""), "foo")
	child.Expose(reflect.TypeOf(""))
	value := parent.CreateContainer().GetInstance(nil /* context */, reflect.TypeOf(""))
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
	value := child.CreateContainer().GetInstance(nil /* context */, reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected to get parent binding for string value 'foo'")
	}
	defer func() {
		recover() // Expected
	}()
	parent.CreateContainer().GetInstance(nil, reflect.TypeOf(""))
	t.Error("Expected a panic when trying to get type exposed by grandchild from grandparent")
}

func TestAlreadyBoundInChildInjector(t *testing.T) {
	defer func() {
		recover()
	}()

	parent := CreateInjector()
	child := parent.CreateChildInjector()

	child.BindInstance(reflect.TypeOf(""), "foo")
	parent.BindInstance(reflect.TypeOf(""), "foo")
	child.Expose(reflect.TypeOf(""))
	t.Error("Expected to fail because already bound in child injector")
}

func TestExposeAndTag(t *testing.T) {
	parent := CreateInjector()
	child := parent.CreateChildInjector()

	child.BindInstance(reflect.TypeOf(""), "bar")
	parent.BindInstance(reflect.TypeOf(""), "foo")
	child.ExposeAndTag(reflect.TypeOf(""), Thing1{})
	container := parent.CreateContainer()
	value := container.GetTaggedInstance(nil, reflect.TypeOf(""), Thing1{})
	if value == nil || value != "bar" {
		t.Error("Expected the parent to have a binding to the tagged value")
	}
	value = container.GetInstance(nil, reflect.TypeOf(""))
	if value == nil || value != "foo" {
		t.Error("Expected the parent to have an untagged binding")
	}
}

func TestAlreadyBoundInParentInjector(t *testing.T) {
	defer func() {
		recover()
	}()

	parent := CreateInjector()
	child := parent.CreateChildInjector()

	parent.BindInstance(reflect.TypeOf(""), "foo")
	child.BindInstance(reflect.TypeOf(""), "foo")
	t.Error("Expected to fail because already bound in parent injector")
}

func TestCanBindInMultipleChildren(t *testing.T) {
	parent := CreateInjector()
	alice := parent.CreateChildInjector()
	bob := parent.CreateChildInjector()
	alice.BindInstance(reflect.TypeOf(""), "foo")
	bob.BindInstance(reflect.TypeOf(""), "bar")
}

func TestSingletonScopedBinding(t *testing.T) {
	context := MyContext{"TestScopedBinding"}
	i := 100
	injector := CreateInjector()
	scope := CreateSimpleScopeWithName("TestScopedBinding Scope")
	injector.BindInScope(
		reflect.TypeOf(0),
		func(_ Context, _ Container) interface{} {
			i += 1
			return i
		},
		Singleton{})
	scope.Enter(context)
	container := injector.CreateContainer()
	first := container.GetInstance(context, reflect.TypeOf(0))
	container = injector.CreateContainer()
	second := container.GetInstance(context, reflect.TypeOf(0))
	if first != second {
		t.Error(fmt.Sprintf("Scoped binding did not return the same value when still in scope (%d != %d)", first, second))
	}
}

type TestScope struct {}

func TestScopedBindingInvokedWhenScopeResets(t *testing.T) {
	context := MyContext{"First context"}
	i := 100
	injector := CreateInjector()
	scope := CreateSimpleScope()
	injector.BindScope(scope, TestScope{})
	injector.BindInScope(
		reflect.TypeOf(0),
		func(_ Context, _ Container) interface{} {
			fmt.Println("Executing underlying provider having value", i)
			i += 1
			return i
		},
		TestScope{})
	scope.Enter(context)
	container := injector.CreateContainer()
	first := container.GetInstance(context, reflect.TypeOf(0))
	scope.Exit(context)
	context = MyContext{"Second context"}
	scope.Enter(context)
	container = injector.CreateContainer()
	second := container.GetInstance(context, reflect.TypeOf(0))
	if first == second {
		t.Error(fmt.Sprintf("Scoped binding returned the same value when scope reset (%d == %d)", first, second))
	}
}
