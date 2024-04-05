# CODE_GEN

Playground to test code generation with Bazel.

The generator is a simple code generator that generates a GO package with a single hello function.
A basic BUILD.bazel file is provided to build the code generator.

```bzl
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

```bzl
genrule(
    name = "genworld",
    outs = ["hello_world.go"],
    cmd = "$(location //blaze/code_gen/generator:generator) -name hello > $@",
    tools = ["//blaze/code_gen/generator"],
)

genrule(
    name = "genbruce",
    outs = ["hello_bruce.go"],
    cmd = "$(location //blaze/code_gen/generator:generator) -name bruce > $@",
    tools = ["//blaze/code_gen/generator"],
)
```