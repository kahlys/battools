# CODE_GEN

Playground to test code generation with Bazel.

- genrule: _[documentation](https://bazel.build/reference/be/general#genrule)_
- macro creation: _[documentation](https://bazel.build/extending/macros)_
- rule creation: _[documentation](https://bazel.build/extending/rules)_

The generator is a simple code generator that generates a GO package with a single hello function.
A basic BUILD.bazel file is provided to build the code generator.

```python
load("@rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "code_gen_lib",
    srcs = ["main.go"],
    importpath = "github.com/kahlys/battools/blaze/code_gen",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "code_gen",
    embed = [":code_gen_lib"],
    visibility = ["//visibility:public"],
)
```

It can be used as follows:

```sh
bazel run //blaze/code_gen/generator:generator -- -name bruce
```

## genrule

Using genrules to generate the code. Add the following to the BUILD.bazel file.

```python
genrule(
    name = "genworld",
    outs = ["hello_world.go"],
    cmd = "$(location //blaze/code_gen/generator:generator) -name hello > $@",
    tools = ["//blaze/code_gen/generator"],
)
```

Use the following command to generate the code:

```sh
bazel build //blaze/code_gen:genworld
```

## macros

To simplify the use of the genrule, a macro can be defined.

In the generator.bzl file:

```python
def hello_generator(name, arg, visibility=None):
  native.genrule(
    name = name,
    outs = ["hello_" + arg + ".go"],
    cmd = "$(location //blaze/code_gen/generator:generator) -name %s > $@" % arg,
    tools = ["//blaze/code_gen/generator"],
    visibility = visibility,
  )
```

In the BUILD.bazel file:

```python
# import hello_generator macro from generator.bzl
load("//blaze/code_gen/generator:generator.bzl", "hello_generator")

hello_generator(
    name = "genmacrolion",
    arg = "lion",
)
```

Use the following command to generate the code:

```sh
bazel build //blaze/code_gen:genmacrolion
```

## chained macros

To chain macros, the following can be done:

In the generator.bzl file:

```python
def chained_generator(name, args, visibility = None):
    for arg in args:
        native.genrule(
            name = name + arg,
            outs = ["hello_" + arg + ".go"],
            cmd = "$(location //blaze/code_gen/generator:generator) -name %s > $@" % arg,
            tools = ["//blaze/code_gen/generator"],
            visibility = visibility,
        )
```

In the BUILD.bazel file:

```python
# import chained_generator macro from generator.bzl
load("//blaze/code_gen:generator.bzl", "chained_generator")

chained_generator(
    name = "genchain",
    args = [
        "red",
        "blue",
        "green",
    ],
)
```

Use the following command to generate the code:

```sh
bazel build //blaze/code_gen:genchainred
bazel build //blaze/code_gen:genchainblue
bazel build //blaze/code_gen:genchaingreen
```

## custom rule

In a `generated_rule.bzl` file, define a custom rule (here named `mygenrule`) to generate the code.

```python
def _codegen_impl(ctx):
    // rule implementation

mygenrule = rule(
    implementation = _codegen_impl,
    attrs = {
        "param": attr.string(
            doc = "the name passed to the generator",
        ),
        "_generator": attr.label(
            executable = True,
            cfg = "exec",
            default = "//blaze/code_gen/generator:generator",
        ),
    },
)
```

We don't want users to customize the generator attribute, so make it private by prefixing it with an underscore, and assign a default value (our generator). It can be done for others attributes too, but in this example, the `param` attribute is left public, so users can customize it.

Then add the implementation of the rule in the same file.

```python
def _codegen_impl(ctx):
    out = ctx.actions.declare_file("hello_" + ctx.attr.param + ".go")

    args = ctx.actions.args()
    args.add("-name", ctx.attr.param)
    args.add("-out", out.path)

    ctx.actions.run(
        outputs = [out],
        executable = ctx.executable._generator,
        tools = [ctx.executable._generator],
        arguments = [args],
        mnemonic = "HelloGenerator",
    )

    return [
        DefaultInfo(files = depset([out])),
        OutputGroupInfo(
            go_generated_srcs = [out],
        ),
    ]
```

The `ctx.actions.declare_file` method is used to declare the output file. The `ctx.actions.run` method is used to run the generator with the specified arguments.

The rule return some [providers](https://docs.bazel.build/versions/main/skylark/rules.html#providers) that other rules may need.

## usage

To use the generated code, add the following to the BUILD.bazel file:

```python
go_library(
    name = "go_hello",
    srcs = [":gensea"],
    importpath = "github.com/kahlys/battools/blaze/hello",
    visibility = ["//visibility:public"],
)
```

In the `cmd` directory, there is a simple go program that uses the generated code with the configured import path.

```go
package main

import (
	"fmt"

	"github.com/kahlys/battools/blaze/hello"
)

func main() {
	fmt.Println(hello.Hello())
}
```

In the BUILD.bazel file of the `cmd` directory, add _deps_ to the go_library rule with the go_library target previously defined.

```python
go_library(
    name = "cmd_lib",
    srcs = ["main.go"],
    importpath = "github.com/kahlys/battools/blaze/code_gen/cmd",
    visibility = ["//visibility:private"],
    deps = [
        "//blaze/code_gen:go_hello",  # keep
    ],
)

go_binary(
    name = "cmd",
    embed = [":cmd_lib"],
    visibility = ["//visibility:public"],
)
```

Add a `#keep` comment so future run of Gazelle will preserve the manually added target.

Test the code with the following command:

```sh
bazel run //blaze/code_gen/cmd
```