cc_library(
    name = "many_attrs",
    deps = [":many_attrs_file_1_cpp", ":many_attrs_file_2_cpp"],
    hdrs = ["header.h"],
    copts = ["-Wall"],
    defines = ["MY_DEFINE=1"],
    includes = ["include"],
    linkopts = ["-lm"],
    visibility = ["//visibility:public"],
)

cc_library(
    name = "many_attrs_file_1_cpp",
    srcs = ["file_1.cpp"],
    hdrs = ["header.h"],
    copts = ["-Wall"],
    defines = ["MY_DEFINE=1"],
    includes = ["include"],
    linkopts = ["-lm"],
    visibility = ["//visibility:public"],
)

cc_library(
    name = "many_attrs_file_2_cpp",
    srcs = ["file_2.cpp"],
    hdrs = ["header.h"],
    copts = ["-Wall"],
    defines = ["MY_DEFINE=1"],
    includes = ["include"],
    linkopts = ["-lm"],
    visibility = ["//visibility:public"],
)
