"""macro for genrule using generator"""

def hello_generator(name, arg, visibility = None):
    native.genrule(
        name = name,
        outs = ["hello_" + arg + ".go"],
        cmd = "$(location //blaze/code_gen/generator:generator) -name %s > $@" % arg,
        tools = ["//blaze/code_gen/generator"],
        visibility = visibility,
    )

def chained_generator(name, args, visibility = None):
    for arg in args:
        native.genrule(
            name = name + arg,
            outs = ["hello_" + arg + ".go"],
            cmd = "$(location //blaze/code_gen/generator:generator) -name %s > $@" % arg,
            tools = ["//blaze/code_gen/generator"],
            visibility = visibility,
        )
