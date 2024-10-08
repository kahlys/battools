# codex

personnal tools

## Build

It uses bazel to build the project. To find bazel modules, visit [the registry](https://registry.bazel.build/)

### Golang

_[Go with Bzlmod](https://github.com/bazelbuild/rules_go/blob/master/docs/go/core/bzlmod.md)_

### Rustlang

At the root of the project, create a `Cargo.toml` file with the following content:

```toml
[workspace]

members = [
    "path/to/project",
    "path/to/other_project",
]
```

In your `MODULE.bazel` file, add the following lines:

```starlark
bazel_dep(name = "rules_rust", version = "0.45.1")


rust = use_extension("@rules_rust//rust:extensions.bzl", "rust")
rust.toolchain(
    edition = "2021",
    versions = ["1.78.0"],
)
use_repo(rust, "rust_toolchains")

register_toolchains("@rust_toolchains//:all")

crate = use_extension("@rules_rust//crate_universe:extension.bzl", "crate")
crate.from_cargo(
    name = "crates",
    cargo_lockfile = "//:Cargo.lock",
    manifests = [
        "//:Cargo.toml",
        "//path/to/project:Cargo.toml",
        "//path/to/other_project:Cargo.toml",
    ],
)
use_repo(crate, "crates")
```

And in a project's `BUILD.bazel` file (at the same level of the project Cargo.toml file), add the following lines:

```starlark
load("@crates//:defs.bzl", "aliases", "all_crate_deps")
load("@rules_rust//rust:defs.bzl", "rust_binary")

rust_binary(
    name = "project_name",
    srcs = ["src/main.rs"],
    crate_name = "project_name",
    deps = all_crate_deps(),
)
```

## Bazel tips

### Missing generated go import path and dependencies

When running bazel commands like `bazel run @rules_go//go -- mod tidy`, you may encounter an error if generated code is not yet generated. A workaround if to write a script that generates empty files in the expected directories.

```bash
#!/bin/bash
set -e


mkdir -p blaze/hello
echo 'package hello' > blaze/hello/hello.go

bazel run @rules_go//go -- mod tidy
bazel mod tidy

rm -rf blaze/hello
```

## Gazelle tips

### GO import path resolve

To tell gazelle which rule to use for a given import path, you can add a comment in the BUILD.bazel file.

```starlark
# gazelle:resolve go github.com/kahlys/codex/blaze/hello //blaze/codegen:go_hello
```

This can be useful when you have multiple rules that match the same import path. When you run `bazel run //:gazelle`, if you encounter an error about multiple rules, you can tell bazel to resolve the dependencies by adding a comment in the BUILD.bazel file.

Take the following example:

```bash
gazelle: rule //blaze/codegen/cmd:cmd_lib imports "github.com/kahlys/codex/blaze/hello" which matches multiple rules: //blaze/codegen:go_hello and //blaze/codegen:go_hello_sea. # gazelle:resolve may be used to disambiguate
```

In the root BUILD.bazel file, add the following line to resolve the go dependencies.

```starlark
# gazelle:resolve go github.com/kahlys/codex/blaze/hello //blaze/codegen:go_hello
```
