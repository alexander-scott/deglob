cc_library(
    name = "files_starting_with_1",
    srcs = glob(["src/file_1*.cpp"]),
    hdrs = ["src/file_1.h"],
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
    hdrs = glob(["src/file_2*.cpp"]),
)
