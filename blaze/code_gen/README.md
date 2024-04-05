# CODE_GEN

Playground to test code generation with Bazel.

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