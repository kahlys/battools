"""Bazel custom rule to generate go code using the generator."""

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

# create a rule called mygenrule that uses the go_rule function to define it.
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
