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
    name = "codegen_lib",
    srcs = ["main.go"],
    importpath = "github.com/kahlys/battools/blaze/codegen",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "codegen",
    embed = [":codegen_lib"],
    visibility = ["//visibility:public"],
)
```

It can be used as follows:

```sh
bazel run //blaze/codegen/generator:generator -- -name bruce
```

## genrule

Using genrules to generate the code. Add the following to the BUILD.bazel file.

```python
genrule(
    name = "genworld",
    outs = ["hello_world.go"],
    cmd = "$(location //blaze/codegen/generator:generator) -name hello > $@",
    tools = ["//blaze/codegen/generator"],
)
```

Use the following command to generate the code:

```sh
bazel build //blaze/codegen:genworld
```

## macros

To simplify the use of the genrule, a macro can be defined.

In the generator.bzl file:

```python
def hello_generator(name, arg, visibility=None):
  native.genrule(
    name = name,
    outs = ["hello_" + arg + ".go"],
    cmd = "$(location //blaze/codegen/generator:generator) -name %s > $@" % arg,
    tools = ["//blaze/codegen/generator"],
    visibility = visibility,
  )
```

In the BUILD.bazel file:

```python
# import hello_generator macro from generator.bzl
load("//blaze/codegen/generator:generator.bzl", "hello_generator")

hello_generator(
    name = "genmacrolion",
    arg = "lion",
)
```

Use the following command to generate the code:

```sh
bazel build //blaze/codegen:genmacrolion
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
            cmd = "$(location //blaze/codegen/generator:generator) -name %s > $@" % arg,
            tools = ["//blaze/codegen/generator"],
            visibility = visibility,
        )
```

In the BUILD.bazel file:

```python
# import chained_generator macro from generator.bzl
load("//blaze/codegen:generator.bzl", "chained_generator")

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
bazel build //blaze/codegen:genchainred
bazel build //blaze/codegen:genchainblue
bazel build //blaze/codegen:genchaingreen
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
            default = "//blaze/codegen/generator:generator",
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
    importpath = "github.com/kahlys/battools/blaze/codegen/cmd",
    visibility = ["//visibility:private"],
    deps = [
        "//blaze/codegen:go_hello",  # keep
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
bazel run //blaze/codegen/cmd
```

## Aspect bazel build

I will use [write_source_files](https://github.com/aspect-build/bazel-lib/blob/main/docs/write_source_files.md#write_source_files) rule to generate files in my source tree. The benefit of this rule is that it will add a test to ensure that the generated files are up to date.

In the `MODULE.bazel` file, add the following rule:

```starlark
bazel_dep(name = "aspect_bazel_lib", version = "2.7.3")
```

Then in the `BUILD.bazel` file, located in the directory where the generated files are, add the following rule:

```starlark
load("@aspect_bazel_lib//lib:write_source_files.bzl", "write_source_files")

write_source_files(
    name = "write_gen_hello",
    files = {
        "hello.go": "//blaze/codegen:genworld",
    },
)
```

It will generate the `hello.go` file in the same directory where the `BUILD.bazel` file is located, by running the `genworld` rule from the `blaze/codegen` directory.