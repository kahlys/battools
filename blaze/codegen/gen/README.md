# GEN

This directory will contains generated packages. I mainly use this to avoid linting errors in the `blaze` directory by having the generated code in the source tree. But in my bazel rules, I still use rules that generate code.

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