cc_library(
    name = "files_starting_with_1",
    hdrs = glob(["file_1*.cpp"]),
)

cc_binary(
    name = "main",
    srcs = ["main.cpp"],
    deps = [
        ":files_starting_with_1",
        ":files_starting_with_2",
    ],
)

cc_library(
    name = "files_starting_with_2",
    hdrs = glob(["file_2*.cpp"]),
)
