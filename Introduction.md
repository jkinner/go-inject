# Overview #

Go-inject is a way to do dependency injection in the go programming language (golang). It is not officially part of the go language nor is it endorsed by the team working on the go language. Use at your own risk.

The style of injection is inspired by Guice, but takes into account the limited reflection capabilities of go. So, in style, it is more functional. If you ever peek under to hood of Guice, you'll find it's this way, too, but it does a lot of the functional work for Java developers.

## Hello, world! ##
Here is a (very contrived) "Hello, world!" example, in the style of Guice (using the type of the object to be retrieved as the key):

```
import (
  "fmt"
  "code.google.com/p/go-inject"
  "reflect"
)

func main() {
  injector := inject.CreateInjector()
  injector.BindInstance(reflect.TypeOf(""), "world")
  container := injector.CreateContainer()
  fmt.Println(fmt.Sprintf("Hello, %s!", container.GetInstance(nil, reflect.TypeOf(""))))
}
```

Let's walk through the code of the `main()` function, line-by-line:

```
  injector := inject.CreateInjector()
```
This code begins every go-inject program. The `injector` is the interface to the injection system. You could have multiple injectors in a single program, although it is typical to have only one. Each injector is isolated from all the others. In its initial state, the injector is empty.

```
  injector.BindInstance(reflect.TypeOf(""), "world")
```
Without any bindings, there is no way to find any objects in the injector. Binding the instances or factory functions should all be done in one stage of the program. Many configuration functions may be used (similar to using multiple Modules in Guice) and they can all configure the same injector. Within each configuration function, child injectors might also be created and configured.

```
  container := injector.CreateContainer()
```
A `container` is used to retrieve values from an `injector`. Typically, a Container is created to retrieve an instance of each object in the system. The Container enforces certain rules for injection, but most especially that a program cannot have a circular dependency. So, a program should not reuse a container to create multiple top-level objects.

```
  fmt.Println(fmt.Sprintf("Hello, %s!", container.GetInstance(nil, reflect.TypeOf(""))))
```
Let's break this down:
```
  container.GetInstance(nil, reflect.TypeOf(""))
```
This is the only go-inject code in this statement. It retrieves the string object from the container. Because the string object is bound to an instance, there is no danger of circularity. However, if the string object were created by a function, that function may retrieve other objects from the container. It would be an error to try to retrieve the string object while executing the function to retrieve the string object.

# go-inject for Guice users #

Because go has limited reflection capabilities, there are some restrictions to the dependency injection framework in go-inject. Here are some common constructs in Guice and their corresponding forms in go-inject:

## Modules ##
### Guice ###
```
class MyModule extends Module {
  void configure() {
    // Bind everything
  }
}
```

### go-inject ###
```
func ConfigureInjector(injector inject.Injector) {
  // Bind everything
}
```

## Injectors ##
### Guice ###
```
Guice.createInjector()
```

### go-inject ###
```
inject.CreateInjector()
```

## Creating instances ##
### Guice ###
```
object = injector.getInstance(key)
```

### go-inject ###
```
container := injector.CreateContainer()
object := container.getInstance(key)
```

## Private modules ##
### Guice ###
```
class MyModule extends PrivateModule {
  void configure() {
    // Bind everything
    expose(...);
  }
}
```

### go-inject ###
```
func ConfigureInjector(injector inject.Injector) {
  childInjector := injector.CreateChildInjector()
  // Bind everything
  childInjector.Expose(...)
}
```