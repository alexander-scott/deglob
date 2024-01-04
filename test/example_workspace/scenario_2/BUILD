cc_library(
    name = "files_starting_with_1",
    srcs = glob(["file_1*.cpp"]),
    hdrs = ["file_1.h"],
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
