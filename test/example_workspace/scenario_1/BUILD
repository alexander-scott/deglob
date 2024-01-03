cc_library(
    name = "files_with_glob",
    hdrs = glob(["file*.cpp"]),
)

cc_binary(
    name = "main",
    srcs = ["main.cpp"],
    deps = [":files_with_glob"],
)
