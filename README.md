# battools

personnal tools

## Build

It uses bazel to build the project.

- **go**: [Go with Bzlmod](https://github.com/bazelbuild/rules_go/blob/master/docs/go/core/bzlmod.md)
- **rust**: TODO

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

### Multiple rules

When you run `bazel run //:gazelle`, if you encounter an error about multiple rules, you can tell bazel to resolve the dependencies by adding a comment in the BUILD.bazel file.

Take the following example:

```bash
gazelle: rule //blaze/code_gen/cmd:cmd_lib imports "github.com/kahlys/battools/blaze/hello" which matches multiple rules: //blaze/code_gen:go_hello and //blaze/code_gen:go_hello_sea. # gazelle:resolve may be used to disambiguate
```

In the root BUILD.bazel file, add the following line to resolve the go dependencies.

```starlark
# gazelle:resolve go github.com/kahlys/battools/blaze/hello //blaze/code_gen:go_hello
```
